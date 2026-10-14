const discovery = require("clever-discovery");
const kayvee = require("kayvee");
const {commandFactory, circuitFactory, metricsFactory} = require("hystrixjs");
const RollingNumberEvent = require("hystrixjs/lib/metrics/RollingNumberEvent");

const { Errors } = require("./types");

function parseForBaggage(entries) {
  if (!entries) {
    return "";
  }
  // Regular expression for valid characters in keys and values
  const validChars = /^[a-zA-Z0-9!#$%&'*+`\-.^_`|~]+$/;

  const pairs = [];

  entries.forEach((value, key) => {
    const validKey = key.match(validChars) ? key : encodeURIComponent(key);
    const validValue = value.match(validChars) ? value : encodeURIComponent(value);
    pairs.push(`${validKey}=${validValue}`);
  });

  return pairs.join(",");
}

/**
 * The exponential retry policy will retry five times with an exponential backoff.
 * @alias module:swagger-test.RetryPolicies.Exponential
 */
const exponentialRetryPolicy = {
  backoffs() {
    const ret = [];
    let next = 100.0; // milliseconds
    const e = 0.05; // +/- 5% jitter
    while (ret.length < 5) {
      const jitter = ((Math.random() * 2) - 1) * e * next;
      ret.push(next + jitter);
      next *= 2;
    }
    return ret;
  },
  retry(request, err, res) {
    return !err && request.method !== "POST" && request.method !== "PATCH" && res.status >= 500;
  },
};

/**
 * Use this retry policy to retry a request once.
 * @alias module:swagger-test.RetryPolicies.Single
 */
const singleRetryPolicy = {
  backoffs() {
    return [1000];
  },
  retry(request, err, res) {
    return !err && request.method !== "POST" && request.method !== "PATCH" && res.status >= 500;
  },
};

/**
 * Use this retry policy to turn off retries.
 * @alias module:swagger-test.RetryPolicies.None
 */
const noRetryPolicy = {
  backoffs() {
    return [];
  },
  retry() {
    return false;
  },
};

/**
 * Logs a failed request. 5XX responses are logged as errors, and everything else as warnings.
 * @private
 */
function logRequestError(logger, request, response, err) {
  const status = response ? response.status : 0;
  const logData = {
    "backend": "swagger-test",
    "method": request.method,
    "uri": request.url,
    "message": err,
    "status_code": status,
  };
  if (status >= 500) {
    logger.errorD("client-request-finished", logData);
  } else {
    logger.warnD("client-request-finished", logData);
  }
}

/**
 * Parses a JSON response body, returning undefined if the body is empty.
 * @private
 */
function parseJSON(body) {
  return body ? JSON.parse(body) : undefined;
}

/**
 * Adds a query parameter to a URL, skipping undefined values. Arrays are added
 * as repeated parameters, e.g. ?tag=a&tag=b.
 * @private
 */
function appendQueryParam(url, name, value) {
  if (value === undefined) {
    return;
  }
  for (const v of [].concat(value)) {
    url.searchParams.append(name, v);
  }
}

/**
 * Takes a promise and uses the provided callback (if any) to handle promise
 * resolutions and rejections
 * @private
 */
function applyCallback(promise, cb) {
  if (!cb) {
    return promise;
  }
  return promise.then((result) => {
    cb(null, result);
  }).catch((err) => {
    cb(err);
  });
}

/**
 * Default circuit breaker options.
 * @alias module:swagger-test.DefaultCircuitOptions
 */
const defaultCircuitOptions = {
  forceClosed:            true,
  requestVolumeThreshold: 20,
  maxConcurrentRequests:  100,
  requestVolumeThreshold: 20,
  sleepWindow:            5000,
  errorPercentThreshold:  90,
  logIntervalMs:          30000
};

/**
 * swagger-test client library.
 * @module swagger-test
 * @typicalname SwaggerTest
 */

/**
 * swagger-test client
 * @alias module:swagger-test
 */
class SwaggerTest {

  /**
   * Create a new client object.
   * @param {Object} options - Options for constructing a client object.
   * @param {string} [options.address] - URL where the server is located. Must provide
   * this or the discovery argument
   * @param {bool} [options.discovery] - Use clever-discovery to locate the server. Must provide
   * this or the address argument
   * @param {number} [options.timeout] - The timeout to use for all client requests,
   * in milliseconds. This can be overridden on a per-request basis. Default is 5000ms.
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy=RetryPolicies.Single] - The logic to
   * determine which requests to retry, as well as how many times to retry.
   * @param {module:kayvee.Logger} [options.logger=logger.New("swagger-test-wagclient")] - The Kayvee
   * logger to use in the client.
   * @param {Object} [options.circuit] - Options for constructing the client's circuit breaker.
   * @param {bool} [options.circuit.forceClosed] - When set to true the circuit will always be closed. Default: true.
   * @param {number} [options.circuit.maxConcurrentRequests] - the maximum number of concurrent requests
   * the client can make at the same time. Default: 100.
   * @param {number} [options.circuit.requestVolumeThreshold] - The minimum number of requests needed
   * before a circuit can be tripped due to health. Default: 20.
   * @param {number} [options.circuit.sleepWindow] - how long, in milliseconds, to wait after a circuit opens
   * before testing for recovery. Default: 5000.
   * @param {number} [options.circuit.errorPercentThreshold] - the threshold to place on the rolling error
   * rate. Once the error rate exceeds this percentage, the circuit opens.
   * Default: 90.
   * @param {object} [options.asynclocalstore] a request scoped async store
   */
  constructor(options) {
    options = options || {};

    if (options.discovery) {
      try {
        this.address = discovery(options.serviceName || "swagger-test", "http").url();
      } catch (e) {
        this.address = discovery(options.serviceName || "swagger-test", "default").url();
      }
    } else if (options.address) {
      this.address = options.address;
    } else {
      throw new Error("Cannot initialize swagger-test without discovery or address");
    }
    this.timeout = options.timeout || 5000;
    this.retryPolicy = options.retryPolicy || singleRetryPolicy;
    this.logger = options.logger || new kayvee.logger((options.serviceName || "swagger-test") + "-wagclient");
    this.asynclocalstore = options.asynclocalstore;

    const circuitOptions = Object.assign({}, defaultCircuitOptions, options.circuit);
    // hystrix implements a caching mechanism, we don't want this or we can't trust that clients
    // are initialized with the values passed in.
    commandFactory.resetCache();
    circuitFactory.resetCache();
    metricsFactory.resetCache();
    this._hystrixCommand = commandFactory.getOrCreate(options.serviceName || "swagger-test").
      errorHandler(this._hystrixCommandErrorHandler).
      circuitBreakerForceClosed(circuitOptions.forceClosed).
      requestVolumeRejectionThreshold(circuitOptions.maxConcurrentRequests).
      circuitBreakerRequestVolumeThreshold(circuitOptions.requestVolumeThreshold).
      circuitBreakerSleepWindowInMilliseconds(circuitOptions.sleepWindow).
      circuitBreakerErrorThresholdPercentage(circuitOptions.errorPercentThreshold).
      timeout(0).
      statisticalWindowLength(10000).
      statisticalWindowNumberOfBuckets(10).
      run(this._hystrixCommandRun).
      context(this).
      build();

    this._logCircuitStateInterval = setInterval(() => this._logCircuitState(), circuitOptions.logIntervalMs);
  }

  /**
  * Releases handles used in client
  */
  close() {
    clearInterval(this._logCircuitStateInterval);
  }

  _hystrixCommandErrorHandler(err) {
    // to avoid counting 4XXs as errors, only count an error if it comes from sending the request
    if (err._fromRequest === true) {
      return err;
    }
    return false;
  }

  _hystrixCommandRun(method, args) {
    return method.apply(this, args);
  }

  _logCircuitState(logger) {
    // code below heavily borrows from hystrix's internal HystrixSSEStream.js logic
    const metrics = this._hystrixCommand.metrics;
    const healthCounts = metrics.getHealthCounts()
    const circuitBreaker = this._hystrixCommand.circuitBreaker;
    this.logger.infoD("swagger-test", {
      "requestCount":                    healthCounts.totalCount,
      "errorCount":                      healthCounts.errorCount,
      "errorPercentage":                 healthCounts.errorPercentage,
      "isCircuitBreakerOpen":            circuitBreaker.isOpen(),
      "rollingCountFailure":             metrics.getRollingCount(RollingNumberEvent.FAILURE),
      "rollingCountShortCircuited":      metrics.getRollingCount(RollingNumberEvent.SHORT_CIRCUITED),
      "rollingCountSuccess":             metrics.getRollingCount(RollingNumberEvent.SUCCESS),
      "rollingCountTimeout":             metrics.getRollingCount(RollingNumberEvent.TIMEOUT),
      "currentConcurrentExecutionCount": metrics.getCurrentExecutionCount(),
      "latencyTotalMean":                metrics.getExecutionTime("mean") || 0,
    });
  }

  /**
   * Returns the headers sent with every request: the caller's headers, plus baggage
   * and the operation and client version.
   * @private
   */
  _headers(operation, options) {
    const storeContext = this.asynclocalstore?.get("context") || new Map();
    const baggage = new Map([...storeContext, ...(options.baggage || new Map())]);

    const headers = new Headers(options.headers);
    headers.set("baggage", parseForBaggage(baggage));
    headers.set("Canonical-Resource", operation);
    headers.set(versionHeader, version);
    return headers;
  }

  /**
   * Sends a request with fetch, retrying it according to the retry policy. Resolves with the
   * response and its body as text. The timeout covers the entire request, including reading
   * the body.
   * @private
   */
  async _send(request, options) {
    const timeout = options.timeout || this.timeout;
    const retryPolicy = options.retryPolicy || this.retryPolicy;
    const backoffs = retryPolicy.backoffs();

    for (let retries = 0; ; retries++) {
      let response, body, err;
      try {
        response = await fetch(request.url, {
          method: request.method,
          headers: request.headers,
          body: request.body,
          signal: AbortSignal.timeout(timeout),
        });
        body = await response.text();
      } catch (e) {
        err = e;
      }

      if (retries < backoffs.length && retryPolicy.retry(request, err, response, body)) {
        await new Promise((resolve) => setTimeout(resolve, backoffs[retries]));
        continue;
      }
      if (err) {
        err._fromRequest = true;
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      return { response, body };
    }
  }

  /**
   * Gets authors
   * @param {Object} params
   * @param {string} [params.name]
   * @param {string} [params.startingAfter]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {Object.<string, string>} [options.headers] - Additional headers to send with the request
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:swagger-test.Errors.BadRequest}
   * @reject {module:swagger-test.Errors.InternalError}
   * @reject {Error}
   */
  getAuthors(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAuthors, arguments), callback);
  }

  async _getAuthors(params, options, cb) {
    if (!options || typeof options === "function") {
      options = {};
    }

    const headers = this._headers("getAuthors", options);

    const url = new URL(this.address + "/v1/authors");
    appendQueryParam(url, "name", params.name);
    appendQueryParam(url, "startingAfter", params.startingAfter);

    const request = {
      method: "GET",
      url: url.href,
      headers,
    };

    const { response, body } = await this._send(request, options);
    switch (response.status) {
      case 200: {
        return parseJSON(body);
      }
      case 400: {
        const err = new Errors.BadRequest(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      case 500: {
        const err = new Errors.InternalError(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      default: {
        const err = new Error("Received unexpected statusCode " + response.status);
        logRequestError(this.logger, request, response, err);
        throw err;
      }
    }
  }


  /**
   * Gets authors
   * @param {Object} params
   * @param {string} [params.name]
   * @param {string} [params.startingAfter]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {Object.<string, string>} [options.headers] - Additional headers to send with the request
   * @returns {Object} iter
   * @returns {function} iter.map - takes in a function, applies it to each resource, and returns a promise to the result as an array
   * @returns {function} iter.toArray - returns a promise to the resources as an array
   * @returns {function} iter.forEach - takes in a function, applies it to each resource
   * @returns {function} iter.forEachAsync - takes in an async function, applies it to each resource
   */
  getAuthorsIter(params, options) {
    const it = (f, saveResults, isAsync) => this._getAuthorsIter(params, options, f, saveResults, isAsync);
    return {
      map: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, true, false]), cb),
      toArray: cb => applyCallback(this._hystrixCommand.execute(it, [x => x, true, false]), cb),
      forEach: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, false]), cb),
      forEachAsync: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, true]), cb),
    };
  }

  async _getAuthorsIter(params, options, f, saveResults, isAsync) {
    options = options || {};

    const headers = this._headers("getAuthors", options);

    const url = new URL(this.address + "/v1/authors");
    appendQueryParam(url, "name", params.name);
    appendQueryParam(url, "startingAfter", params.startingAfter);

    const request = {
      method: "GET",
      url: url.href,
      headers,
    };

    let results = [];
    while (request.url) {
      const { response, body } = await this._send(request, options);
      switch (response.status) {
        case 200: {
          const page = parseJSON(body);
          const items = page.authorSet.results;
          if (saveResults) {
            results = results.concat(items.map(f));
          } else if (isAsync) {
            for (let i = 0; i < items.length; i++) {
              await f(items[i], i, page);
            }
          } else {
            items.forEach(f);
          }
          break;
        }
        case 400: {
          const err = new Errors.BadRequest(parseJSON(body) || {});
          logRequestError(this.logger, request, response, err);
          throw err;
        }
        case 500: {
          const err = new Errors.InternalError(parseJSON(body) || {});
          logRequestError(this.logger, request, response, err);
          throw err;
        }
        default: {
          const err = new Error("Received unexpected statusCode " + response.status);
          logRequestError(this.logger, request, response, err);
          throw err;
        }
      }

      const nextPagePath = response.headers.get("x-next-page-path");
      request.url = nextPagePath ? this.address + nextPagePath : "";
    }
    return saveResults ? results : undefined;
  }

  /**
   * Gets authors, but needs to use the body so it's a PUT
   * @param {Object} params
   * @param {string} [params.name]
   * @param {string} [params.startingAfter]
   * @param [params.favoriteBooks]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {Object.<string, string>} [options.headers] - Additional headers to send with the request
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:swagger-test.Errors.BadRequest}
   * @reject {module:swagger-test.Errors.InternalError}
   * @reject {Error}
   */
  getAuthorsWithPut(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAuthorsWithPut, arguments), callback);
  }

  async _getAuthorsWithPut(params, options, cb) {
    if (!options || typeof options === "function") {
      options = {};
    }

    const headers = this._headers("getAuthorsWithPut", options);
    headers.set("Content-Type", "application/json");

    const url = new URL(this.address + "/v1/authors");
    appendQueryParam(url, "name", params.name);
    appendQueryParam(url, "startingAfter", params.startingAfter);

    const request = {
      method: "PUT",
      url: url.href,
      headers,
      body: JSON.stringify(params.favoriteBooks),
    };

    const { response, body } = await this._send(request, options);
    switch (response.status) {
      case 200: {
        return parseJSON(body);
      }
      case 400: {
        const err = new Errors.BadRequest(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      case 500: {
        const err = new Errors.InternalError(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      default: {
        const err = new Error("Received unexpected statusCode " + response.status);
        logRequestError(this.logger, request, response, err);
        throw err;
      }
    }
  }


  /**
   * Gets authors, but needs to use the body so it's a PUT
   * @param {Object} params
   * @param {string} [params.name]
   * @param {string} [params.startingAfter]
   * @param [params.favoriteBooks]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {Object.<string, string>} [options.headers] - Additional headers to send with the request
   * @returns {Object} iter
   * @returns {function} iter.map - takes in a function, applies it to each resource, and returns a promise to the result as an array
   * @returns {function} iter.toArray - returns a promise to the resources as an array
   * @returns {function} iter.forEach - takes in a function, applies it to each resource
   * @returns {function} iter.forEachAsync - takes in an async function, applies it to each resource
   */
  getAuthorsWithPutIter(params, options) {
    const it = (f, saveResults, isAsync) => this._getAuthorsWithPutIter(params, options, f, saveResults, isAsync);
    return {
      map: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, true, false]), cb),
      toArray: cb => applyCallback(this._hystrixCommand.execute(it, [x => x, true, false]), cb),
      forEach: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, false]), cb),
      forEachAsync: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, true]), cb),
    };
  }

  async _getAuthorsWithPutIter(params, options, f, saveResults, isAsync) {
    options = options || {};

    const headers = this._headers("getAuthorsWithPut", options);
    headers.set("Content-Type", "application/json");

    const url = new URL(this.address + "/v1/authors");
    appendQueryParam(url, "name", params.name);
    appendQueryParam(url, "startingAfter", params.startingAfter);

    const request = {
      method: "PUT",
      url: url.href,
      headers,
      body: JSON.stringify(params.favoriteBooks),
    };

    let results = [];
    while (request.url) {
      const { response, body } = await this._send(request, options);
      switch (response.status) {
        case 200: {
          const page = parseJSON(body);
          const items = page.authorSet.results;
          if (saveResults) {
            results = results.concat(items.map(f));
          } else if (isAsync) {
            for (let i = 0; i < items.length; i++) {
              await f(items[i], i, page);
            }
          } else {
            items.forEach(f);
          }
          break;
        }
        case 400: {
          const err = new Errors.BadRequest(parseJSON(body) || {});
          logRequestError(this.logger, request, response, err);
          throw err;
        }
        case 500: {
          const err = new Errors.InternalError(parseJSON(body) || {});
          logRequestError(this.logger, request, response, err);
          throw err;
        }
        default: {
          const err = new Error("Received unexpected statusCode " + response.status);
          logRequestError(this.logger, request, response, err);
          throw err;
        }
      }

      const nextPagePath = response.headers.get("x-next-page-path");
      request.url = nextPagePath ? this.address + nextPagePath : "";
    }
    return saveResults ? results : undefined;
  }

  /**
   * Returns a list of books
   * @param {Object} params
   * @param {string[]} [params.authors] - A list of authors. Must specify at least one and at most two
   * @param {boolean} [params.available=true]
   * @param {string} [params.state=finished]
   * @param {string} [params.published]
   * @param {string} [params.snakeCase]
   * @param {string} [params.completed]
   * @param {number} [params.maxPages=500.5]
   * @param {number} [params.minPages=5]
   * @param {number} [params.pagesToTime]
   * @param {string} [params.authorization]
   * @param {number} [params.startingAfter]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {Object.<string, string>} [options.headers] - Additional headers to send with the request
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object[]}
   * @reject {module:swagger-test.Errors.BadRequest}
   * @reject {module:swagger-test.Errors.InternalError}
   * @reject {Error}
   */
  getBooks(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getBooks, arguments), callback);
  }

  async _getBooks(params, options, cb) {
    if (!options || typeof options === "function") {
      options = {};
    }

    const headers = this._headers("getBooks", options);
    if (params.authorization !== undefined) {
      headers.set("authorization", params.authorization);
    }

    const url = new URL(this.address + "/v1/books");
    appendQueryParam(url, "authors", params.authors);
    appendQueryParam(url, "available", params.available);
    appendQueryParam(url, "state", params.state);
    appendQueryParam(url, "published", params.published);
    appendQueryParam(url, "snake_case", params.snakeCase);
    appendQueryParam(url, "completed", params.completed);
    appendQueryParam(url, "maxPages", params.maxPages);
    appendQueryParam(url, "min_pages", params.minPages);
    appendQueryParam(url, "pagesToTime", params.pagesToTime);
    appendQueryParam(url, "startingAfter", params.startingAfter);

    const request = {
      method: "GET",
      url: url.href,
      headers,
    };

    const { response, body } = await this._send(request, options);
    switch (response.status) {
      case 200: {
        return parseJSON(body);
      }
      case 400: {
        const err = new Errors.BadRequest(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      case 500: {
        const err = new Errors.InternalError(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      default: {
        const err = new Error("Received unexpected statusCode " + response.status);
        logRequestError(this.logger, request, response, err);
        throw err;
      }
    }
  }


  /**
   * Returns a list of books
   * @param {Object} params
   * @param {string[]} [params.authors] - A list of authors. Must specify at least one and at most two
   * @param {boolean} [params.available=true]
   * @param {string} [params.state=finished]
   * @param {string} [params.published]
   * @param {string} [params.snakeCase]
   * @param {string} [params.completed]
   * @param {number} [params.maxPages=500.5]
   * @param {number} [params.minPages=5]
   * @param {number} [params.pagesToTime]
   * @param {string} [params.authorization]
   * @param {number} [params.startingAfter]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {Object.<string, string>} [options.headers] - Additional headers to send with the request
   * @returns {Object} iter
   * @returns {function} iter.map - takes in a function, applies it to each resource, and returns a promise to the result as an array
   * @returns {function} iter.toArray - returns a promise to the resources as an array
   * @returns {function} iter.forEach - takes in a function, applies it to each resource
   * @returns {function} iter.forEachAsync - takes in an async function, applies it to each resource
   */
  getBooksIter(params, options) {
    const it = (f, saveResults, isAsync) => this._getBooksIter(params, options, f, saveResults, isAsync);
    return {
      map: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, true, false]), cb),
      toArray: cb => applyCallback(this._hystrixCommand.execute(it, [x => x, true, false]), cb),
      forEach: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, false]), cb),
      forEachAsync: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, true]), cb),
    };
  }

  async _getBooksIter(params, options, f, saveResults, isAsync) {
    options = options || {};

    const headers = this._headers("getBooks", options);
    if (params.authorization !== undefined) {
      headers.set("authorization", params.authorization);
    }

    const url = new URL(this.address + "/v1/books");
    appendQueryParam(url, "authors", params.authors);
    appendQueryParam(url, "available", params.available);
    appendQueryParam(url, "state", params.state);
    appendQueryParam(url, "published", params.published);
    appendQueryParam(url, "snake_case", params.snakeCase);
    appendQueryParam(url, "completed", params.completed);
    appendQueryParam(url, "maxPages", params.maxPages);
    appendQueryParam(url, "min_pages", params.minPages);
    appendQueryParam(url, "pagesToTime", params.pagesToTime);
    appendQueryParam(url, "startingAfter", params.startingAfter);

    const request = {
      method: "GET",
      url: url.href,
      headers,
    };

    let results = [];
    while (request.url) {
      const { response, body } = await this._send(request, options);
      switch (response.status) {
        case 200: {
          const page = parseJSON(body);
          const items = page;
          if (saveResults) {
            results = results.concat(items.map(f));
          } else if (isAsync) {
            for (let i = 0; i < items.length; i++) {
              await f(items[i], i, page);
            }
          } else {
            items.forEach(f);
          }
          break;
        }
        case 400: {
          const err = new Errors.BadRequest(parseJSON(body) || {});
          logRequestError(this.logger, request, response, err);
          throw err;
        }
        case 500: {
          const err = new Errors.InternalError(parseJSON(body) || {});
          logRequestError(this.logger, request, response, err);
          throw err;
        }
        default: {
          const err = new Error("Received unexpected statusCode " + response.status);
          logRequestError(this.logger, request, response, err);
          throw err;
        }
      }

      const nextPagePath = response.headers.get("x-next-page-path");
      request.url = nextPagePath ? this.address + nextPagePath : "";
    }
    return saveResults ? results : undefined;
  }

  /**
   * Creates a book
   * @param newBook
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {Object.<string, string>} [options.headers] - Additional headers to send with the request
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:swagger-test.Errors.BadRequest}
   * @reject {module:swagger-test.Errors.InternalError}
   * @reject {Error}
   */
  createBook(newBook, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._createBook, arguments), callback);
  }

  async _createBook(newBook, options, cb) {
    const params = {};
    params["newBook"] = newBook;

    if (!options || typeof options === "function") {
      options = {};
    }

    const headers = this._headers("createBook", options);
    headers.set("Content-Type", "application/json");

    const url = new URL(this.address + "/v1/books");

    const request = {
      method: "POST",
      url: url.href,
      headers,
      body: JSON.stringify(params.newBook),
    };

    const { response, body } = await this._send(request, options);
    switch (response.status) {
      case 200: {
        return parseJSON(body);
      }
      case 400: {
        const err = new Errors.BadRequest(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      case 500: {
        const err = new Errors.InternalError(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      default: {
        const err = new Error("Received unexpected statusCode " + response.status);
        logRequestError(this.logger, request, response, err);
        throw err;
      }
    }
  }

  /**
   * Puts a book
   * @param newBook
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {Object.<string, string>} [options.headers] - Additional headers to send with the request
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:swagger-test.Errors.BadRequest}
   * @reject {module:swagger-test.Errors.InternalError}
   * @reject {Error}
   */
  putBook(newBook, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._putBook, arguments), callback);
  }

  async _putBook(newBook, options, cb) {
    const params = {};
    params["newBook"] = newBook;

    if (!options || typeof options === "function") {
      options = {};
    }

    const headers = this._headers("putBook", options);
    headers.set("Content-Type", "application/json");

    const url = new URL(this.address + "/v1/books");

    const request = {
      method: "PUT",
      url: url.href,
      headers,
      body: JSON.stringify(params.newBook),
    };

    const { response, body } = await this._send(request, options);
    switch (response.status) {
      case 200: {
        return parseJSON(body);
      }
      case 400: {
        const err = new Errors.BadRequest(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      case 500: {
        const err = new Errors.InternalError(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      default: {
        const err = new Error("Received unexpected statusCode " + response.status);
        logRequestError(this.logger, request, response, err);
        throw err;
      }
    }
  }

  /**
   * Returns a book
   * @param {Object} params
   * @param {number} params.bookID
   * @param {string} [params.authorID]
   * @param {string} [params.authorization]
   * @param {string} [params.XDontRateLimitMeBro]
   * @param {string} [params.randomBytes]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {Object.<string, string>} [options.headers] - Additional headers to send with the request
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:swagger-test.Errors.BadRequest}
   * @reject {module:swagger-test.Errors.Unathorized}
   * @reject {module:swagger-test.Errors.Error}
   * @reject {module:swagger-test.Errors.InternalError}
   * @reject {Error}
   */
  getBookByID(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getBookByID, arguments), callback);
  }

  async _getBookByID(params, options, cb) {
    if (!options || typeof options === "function") {
      options = {};
    }
    if (!params.bookID) {
      throw new Error("bookID must be non-empty because it's a path parameter");
    }

    const headers = this._headers("getBookByID", options);
    if (params.authorization !== undefined) {
      headers.set("authorization", params.authorization);
    }
    if (params.XDontRateLimitMeBro !== undefined) {
      headers.set("X-Dont-Rate-Limit-Me-Bro", params.XDontRateLimitMeBro);
    }

    const url = new URL(this.address + "/v1/books/" + encodeURIComponent(params.bookID) + "");
    appendQueryParam(url, "authorID", params.authorID);
    appendQueryParam(url, "randomBytes", params.randomBytes);

    const request = {
      method: "GET",
      url: url.href,
      headers,
    };

    const { response, body } = await this._send(request, options);
    switch (response.status) {
      case 200: {
        return parseJSON(body);
      }
      case 400: {
        const err = new Errors.BadRequest(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      case 401: {
        const err = new Errors.Unathorized(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      case 404: {
        const err = new Errors.Error(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      case 500: {
        const err = new Errors.InternalError(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      default: {
        const err = new Error("Received unexpected statusCode " + response.status);
        logRequestError(this.logger, request, response, err);
        throw err;
      }
    }
  }

  /**
   * Retrieve a book
   * @param {string} id
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {Object.<string, string>} [options.headers] - Additional headers to send with the request
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:swagger-test.Errors.BadRequest}
   * @reject {module:swagger-test.Errors.Error}
   * @reject {module:swagger-test.Errors.InternalError}
   * @reject {Error}
   */
  getBookByID2(id, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getBookByID2, arguments), callback);
  }

  async _getBookByID2(id, options, cb) {
    const params = {};
    params["id"] = id;

    if (!options || typeof options === "function") {
      options = {};
    }
    if (!params.id) {
      throw new Error("id must be non-empty because it's a path parameter");
    }

    const headers = this._headers("getBookByID2", options);

    const url = new URL(this.address + "/v1/books2/" + encodeURIComponent(params.id) + "");

    const request = {
      method: "GET",
      url: url.href,
      headers,
    };

    const { response, body } = await this._send(request, options);
    switch (response.status) {
      case 200: {
        return parseJSON(body);
      }
      case 400: {
        const err = new Errors.BadRequest(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      case 404: {
        const err = new Errors.Error(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      case 500: {
        const err = new Errors.InternalError(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      default: {
        const err = new Error("Received unexpected statusCode " + response.status);
        logRequestError(this.logger, request, response, err);
        throw err;
      }
    }
  }

  /**
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {Object.<string, string>} [options.headers] - Additional headers to send with the request
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:swagger-test.Errors.BadRequest}
   * @reject {module:swagger-test.Errors.InternalError}
   * @reject {Error}
   */
  healthCheck(options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._healthCheck, arguments), callback);
  }

  async _healthCheck(options, cb) {
    const params = {};

    if (!options || typeof options === "function") {
      options = {};
    }

    const headers = this._headers("healthCheck", options);

    const url = new URL(this.address + "/v1/health/check");

    const request = {
      method: "GET",
      url: url.href,
      headers,
    };

    const { response, body } = await this._send(request, options);
    switch (response.status) {
      case 200: {
        return;
      }
      case 400: {
        const err = new Errors.BadRequest(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      case 500: {
        const err = new Errors.InternalError(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      default: {
        const err = new Error("Received unexpected statusCode " + response.status);
        logRequestError(this.logger, request, response, err);
        throw err;
      }
    }
  }

  /**
   * testing that we can use a lowercase name for a model
   * @param {Object} params
   * @param params.lowercase
   * @param {string} params.pathParam
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {Object.<string, string>} [options.headers] - Additional headers to send with the request
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:swagger-test.Errors.BadRequest}
   * @reject {module:swagger-test.Errors.InternalError}
   * @reject {Error}
   */
  lowercaseModelsTest(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._lowercaseModelsTest, arguments), callback);
  }

  async _lowercaseModelsTest(params, options, cb) {
    if (!options || typeof options === "function") {
      options = {};
    }
    if (!params.pathParam) {
      throw new Error("pathParam must be non-empty because it's a path parameter");
    }

    const headers = this._headers("lowercaseModelsTest", options);
    headers.set("Content-Type", "application/json");

    const url = new URL(this.address + "/v1/lowercaseModelsTest/" + encodeURIComponent(params.pathParam) + "");

    const request = {
      method: "POST",
      url: url.href,
      headers,
      body: JSON.stringify(params.lowercase),
    };

    const { response, body } = await this._send(request, options);
    switch (response.status) {
      case 200: {
        return;
      }
      case 400: {
        const err = new Errors.BadRequest(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      case 500: {
        const err = new Errors.InternalError(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      default: {
        const err = new Error("Received unexpected statusCode " + response.status);
        logRequestError(this.logger, request, response, err);
        throw err;
      }
    }
  }
};

module.exports = SwaggerTest;

/**
 * Retry policies available to use.
 * @alias module:swagger-test.RetryPolicies
 */
module.exports.RetryPolicies = {
  Single: singleRetryPolicy,
  Exponential: exponentialRetryPolicy,
  None: noRetryPolicy,
};

/**
 * Errors returned by methods.
 * @alias module:swagger-test.Errors
 */
module.exports.Errors = Errors;

module.exports.DefaultCircuitOptions = defaultCircuitOptions;

const version = "9.0.0";
const versionHeader = "X-Client-Version";
module.exports.Version = version;
module.exports.VersionHeader = versionHeader;
