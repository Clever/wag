const async = require("async");
const discovery = require("clever-discovery");
const kayvee = require("kayvee");
const request = require("request");
const {commandFactory} = require("hystrixjs");
const RollingNumberEvent = require("hystrixjs/lib/metrics/RollingNumberEvent");

const { Errors } = require("./types");

/**
 * The exponential retry policy will retry five times with an exponential backoff.
 * @alias module:wag/samples.RetryPolicies.Exponential
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
  retry(requestOptions, err, res) {
    if (err || requestOptions.method === "POST" ||
        requestOptions.method === "PATCH" ||
        res.statusCode < 500) {
      return false;
    }
    return true;
  },
};

/**
 * Use this retry policy to retry a request once.
 * @alias module:wag/samples.RetryPolicies.Single
 */
const singleRetryPolicy = {
  backoffs() {
    return [1000];
  },
  retry(requestOptions, err, res) {
    if (err || requestOptions.method === "POST" ||
        requestOptions.method === "PATCH" ||
        res.statusCode < 500) {
      return false;
    }
    return true;
  },
};

/**
 * Use this retry policy to turn off retries.
 * @alias module:wag/samples.RetryPolicies.None
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
 * Request status log is used to
 * to output the status of a request returned
 * by the client.
 * @private
 */
function responseLog(logger, req, res, err) {
  var res = res || { };
  var req = req || { };
  var logData = {
	"backend": "wag/samples",
	"method": req.method || "",
	"uri": req.uri || "",
    "message": err || (res.statusMessage || ""),
    "status_code": res.statusCode || 0,
  };

  if (err) {
    logger.errorD("client-request-finished", logData);
  } else {
    logger.infoD("client-request-finished", logData);
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
 * @alias module:wag/samples.DefaultCircuitOptions
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
 * wag/samples client library.
 * @module wag/samples
 * @typicalname WagSamples
 */

/**
 * wag/samples client
 * @alias module:wag/samples
 */
class WagSamples {

  /**
   * Create a new client object.
   * @param {Object} options - Options for constructing a client object.
   * @param {string} [options.address] - URL where the server is located. Must provide
   * this or the discovery argument
   * @param {bool} [options.discovery] - Use clever-discovery to locate the server. Must provide
   * this or the address argument
   * @param {number} [options.timeout] - The timeout to use for all client requests,
   * in milliseconds. This can be overridden on a per-request basis. Default is 5000ms.
   * @param {bool} [options.keepalive] - Set keepalive to true for client requests. This sets the
   * forever: true attribute in request. Defaults to true.
   * @param {module:wag/samples.RetryPolicies} [options.retryPolicy=RetryPolicies.Single] - The logic to
   * determine which requests to retry, as well as how many times to retry.
   * @param {module:kayvee.Logger} [options.logger=logger.New("wag/samples-wagclient")] - The Kayvee
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
   */
  constructor(options) {
    options = options || {};

    if (options.discovery) {
      try {
        this.address = discovery(options.serviceName || "wag/samples", "http").url();
      } catch (e) {
        this.address = discovery(options.serviceName || "wag/samples", "default").url();
      }
    } else if (options.address) {
      this.address = options.address;
    } else {
      throw new Error("Cannot initialize wag/samples without discovery or address");
    }
    if (options.keepalive !== undefined) {
      this.keepalive = options.keepalive;
    } else {
      this.keepalive = true;
    }
    if (options.timeout) {
      this.timeout = options.timeout;
    } else {
      this.timeout = 5000;
    }
    if (options.retryPolicy) {
      this.retryPolicy = options.retryPolicy;
    }
    if (options.logger) {
      this.logger = options.logger;
    } else {
      this.logger = new kayvee.logger((options.serviceName || "wag/samples") + "-wagclient");
    }

    const circuitOptions = Object.assign({}, defaultCircuitOptions, options.circuit);
    this._hystrixCommand = commandFactory.getOrCreate(options.serviceName || "wag/samples").
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
    // to avoid counting 4XXs as errors, only count an error if it comes from the request library
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
    this.logger.infoD("wag/samples", {
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
   * Gets authors
   * @param {Object} params
   * @param {string} [params.name]
   * @param {string} [params.startingAfter]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {module:wag/samples.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:wag/samples.Errors.BadRequest}
   * @reject {module:wag/samples.Errors.InternalError}
   * @reject {Error}
   */
  getAuthors(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAuthors, arguments), callback);
  }

  _getAuthors(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;

      const headers = {};
      headers["Canonical-Resource"] = "getAuthors";
      headers[versionHeader] = version;

      const query = {};
      if (typeof params.name !== "undefined") {
        query["name"] = params.name;
      }

      if (typeof params.startingAfter !== "undefined") {
        query["startingAfter"] = params.startingAfter;
      }


      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/authors",
        gzip: true,
        json: true,
        timeout,
        headers,
        qs: query,
        useQuerystring: true,
      };
      if (this.keepalive) {
        requestOptions.forever = true;
      }


      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      const logger = this.logger;

      let retries = 0;
      (function requestOnce() {
        request(requestOptions, (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            setTimeout(requestOnce, backoff);
            return;
          }
          if (err) {
            err._fromRequest = true;
            responseLog(logger, requestOptions, response, err)
            reject(err);
            return;
          }

          switch (response.statusCode) {
            case 200:
              resolve(body);
              break;

            case 400:
              var err = new Errors.BadRequest(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 500:
              var err = new Errors.InternalError(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            default:
              var err = new Error("Received unexpected statusCode " + response.statusCode);
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;
          }
        });
      }());
    });
  }


  /**
   * Gets authors
   * @param {Object} params
   * @param {string} [params.name]
   * @param {string} [params.startingAfter]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {module:wag/samples.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @returns {Object} iter
   * @returns {function} iter.map - takes in a function, applies it to each resource, and returns a promise to the result as an array
   * @returns {function} iter.toArray - returns a promise to the resources as an array
   * @returns {function} iter.forEach - takes in a function, applies it to each resource
   * @returns {function} iter.forEachAsync - takes in an async function, applies it to each resource
   */
  getAuthorsIter(params, options) {
    const it = (f, saveResults, isAsync) => new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;

      const headers = {};
      headers["Canonical-Resource"] = "getAuthors";
      headers[versionHeader] = version;

      const query = {};
      if (typeof params.name !== "undefined") {
        query["name"] = params.name;
      }

      if (typeof params.startingAfter !== "undefined") {
        query["startingAfter"] = params.startingAfter;
      }


      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/authors",
        gzip: true,
        json: true,
        timeout,
        headers,
        qs: query,
        useQuerystring: true,
      };
      if (this.keepalive) {
        requestOptions.forever = true;
      }


      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      const logger = this.logger;

      let results = [];
      async.whilst(
        () => requestOptions.uri !== "",
        cbW => {
      const address = this.address;
      let retries = 0;
      (function requestOnce() {
        request(requestOptions, async (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            setTimeout(requestOnce, backoff);
            return;
          }
          if (err) {
            err._fromRequest = true;
            responseLog(logger, requestOptions, response, err)
            cbW(err);
            return;
          }

          switch (response.statusCode) {
            case 200:
              if (saveResults) {
                results = results.concat(body.authorSet.results.map(f));
              } else {
                if (isAsync) {
                  for (let i = 0; i < body.authorSet.results.length; i++) {
                    try {
                      await f(body.authorSet.results[i], i, body);
                    } catch(err) {
                      reject(err);
                    }
                  }
                } else {
                  body.forEach(f)
                }
              }
              break;

            case 400:
              var err = new Errors.BadRequest(body || {});
              responseLog(logger, requestOptions, response, err);
              cbW(err);
              return;

            case 500:
              var err = new Errors.InternalError(body || {});
              responseLog(logger, requestOptions, response, err);
              cbW(err);
              return;

            default:
              var err = new Error("Received unexpected statusCode " + response.statusCode);
              responseLog(logger, requestOptions, response, err);
              cbW(err);
              return;
          }

          requestOptions.qs = null;
          requestOptions.useQuerystring = false;
          requestOptions.uri = "";
          if (response.headers["x-next-page-path"]) {
            requestOptions.uri = address + response.headers["x-next-page-path"];
          }
          cbW();
        });
      }());
        },
        err => {
          if (err) {
            reject(err);
            return;
          }
          if (saveResults) {
            resolve(results);
          } else {
            resolve();
          }
        }
      );
    });

    return {
      map: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, true, false]), cb),
      toArray: cb => applyCallback(this._hystrixCommand.execute(it, [x => x, true, false]), cb),
      forEach: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, false]), cb),
      forEachAsync: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, true]), cb),
    };
  }

  /**
   * Gets authors, but needs to use the body so it's a PUT
   * @param {Object} params
   * @param {string} [params.name]
   * @param {string} [params.startingAfter]
   * @param [params.favoriteBooks]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {module:wag/samples.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:wag/samples.Errors.BadRequest}
   * @reject {module:wag/samples.Errors.InternalError}
   * @reject {Error}
   */
  getAuthorsWithPut(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAuthorsWithPut, arguments), callback);
  }

  _getAuthorsWithPut(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;

      const headers = {};
      headers["Canonical-Resource"] = "getAuthorsWithPut";
      headers[versionHeader] = version;

      const query = {};
      if (typeof params.name !== "undefined") {
        query["name"] = params.name;
      }

      if (typeof params.startingAfter !== "undefined") {
        query["startingAfter"] = params.startingAfter;
      }


      const requestOptions = {
        method: "PUT",
        uri: this.address + "/v1/authors",
        gzip: true,
        json: true,
        timeout,
        headers,
        qs: query,
        useQuerystring: true,
      };
      if (this.keepalive) {
        requestOptions.forever = true;
      }

      requestOptions.body = params.favoriteBooks;


      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      const logger = this.logger;

      let retries = 0;
      (function requestOnce() {
        request(requestOptions, (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            setTimeout(requestOnce, backoff);
            return;
          }
          if (err) {
            err._fromRequest = true;
            responseLog(logger, requestOptions, response, err)
            reject(err);
            return;
          }

          switch (response.statusCode) {
            case 200:
              resolve(body);
              break;

            case 400:
              var err = new Errors.BadRequest(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 500:
              var err = new Errors.InternalError(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            default:
              var err = new Error("Received unexpected statusCode " + response.statusCode);
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;
          }
        });
      }());
    });
  }


  /**
   * Gets authors, but needs to use the body so it's a PUT
   * @param {Object} params
   * @param {string} [params.name]
   * @param {string} [params.startingAfter]
   * @param [params.favoriteBooks]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {module:wag/samples.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @returns {Object} iter
   * @returns {function} iter.map - takes in a function, applies it to each resource, and returns a promise to the result as an array
   * @returns {function} iter.toArray - returns a promise to the resources as an array
   * @returns {function} iter.forEach - takes in a function, applies it to each resource
   * @returns {function} iter.forEachAsync - takes in an async function, applies it to each resource
   */
  getAuthorsWithPutIter(params, options) {
    const it = (f, saveResults, isAsync) => new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;

      const headers = {};
      headers["Canonical-Resource"] = "getAuthorsWithPut";
      headers[versionHeader] = version;

      const query = {};
      if (typeof params.name !== "undefined") {
        query["name"] = params.name;
      }

      if (typeof params.startingAfter !== "undefined") {
        query["startingAfter"] = params.startingAfter;
      }


      const requestOptions = {
        method: "PUT",
        uri: this.address + "/v1/authors",
        gzip: true,
        json: true,
        timeout,
        headers,
        qs: query,
        useQuerystring: true,
      };
      if (this.keepalive) {
        requestOptions.forever = true;
      }

      requestOptions.body = params.favoriteBooks;


      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      const logger = this.logger;

      let results = [];
      async.whilst(
        () => requestOptions.uri !== "",
        cbW => {
      const address = this.address;
      let retries = 0;
      (function requestOnce() {
        request(requestOptions, async (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            setTimeout(requestOnce, backoff);
            return;
          }
          if (err) {
            err._fromRequest = true;
            responseLog(logger, requestOptions, response, err)
            cbW(err);
            return;
          }

          switch (response.statusCode) {
            case 200:
              if (saveResults) {
                results = results.concat(body.authorSet.results.map(f));
              } else {
                if (isAsync) {
                  for (let i = 0; i < body.authorSet.results.length; i++) {
                    try {
                      await f(body.authorSet.results[i], i, body);
                    } catch(err) {
                      reject(err);
                    }
                  }
                } else {
                  body.forEach(f)
                }
              }
              break;

            case 400:
              var err = new Errors.BadRequest(body || {});
              responseLog(logger, requestOptions, response, err);
              cbW(err);
              return;

            case 500:
              var err = new Errors.InternalError(body || {});
              responseLog(logger, requestOptions, response, err);
              cbW(err);
              return;

            default:
              var err = new Error("Received unexpected statusCode " + response.statusCode);
              responseLog(logger, requestOptions, response, err);
              cbW(err);
              return;
          }

          requestOptions.qs = null;
          requestOptions.useQuerystring = false;
          requestOptions.uri = "";
          if (response.headers["x-next-page-path"]) {
            requestOptions.uri = address + response.headers["x-next-page-path"];
          }
          cbW();
        });
      }());
        },
        err => {
          if (err) {
            reject(err);
            return;
          }
          if (saveResults) {
            resolve(results);
          } else {
            resolve();
          }
        }
      );
    });

    return {
      map: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, true, false]), cb),
      toArray: cb => applyCallback(this._hystrixCommand.execute(it, [x => x, true, false]), cb),
      forEach: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, false]), cb),
      forEachAsync: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, true]), cb),
    };
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
   * @param {module:wag/samples.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object[]}
   * @reject {module:wag/samples.Errors.BadRequest}
   * @reject {module:wag/samples.Errors.InternalError}
   * @reject {Error}
   */
  getBooks(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getBooks, arguments), callback);
  }

  _getBooks(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;

      const headers = {};
      headers["Canonical-Resource"] = "getBooks";
      headers[versionHeader] = version;
      headers["authorization"] = params.authorization;

      const query = {};
      if (typeof params.authors !== "undefined") {
        query["authors"] = params.authors;
      }

      if (typeof params.available !== "undefined") {
        query["available"] = params.available;
      }

      if (typeof params.state !== "undefined") {
        query["state"] = params.state;
      }

      if (typeof params.published !== "undefined") {
        query["published"] = params.published;
      }

      if (typeof params.snakeCase !== "undefined") {
        query["snake_case"] = params.snakeCase;
      }

      if (typeof params.completed !== "undefined") {
        query["completed"] = params.completed;
      }

      if (typeof params.maxPages !== "undefined") {
        query["maxPages"] = params.maxPages;
      }

      if (typeof params.minPages !== "undefined") {
        query["min_pages"] = params.minPages;
      }

      if (typeof params.pagesToTime !== "undefined") {
        query["pagesToTime"] = params.pagesToTime;
      }

      if (typeof params.startingAfter !== "undefined") {
        query["startingAfter"] = params.startingAfter;
      }


      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/books",
        gzip: true,
        json: true,
        timeout,
        headers,
        qs: query,
        useQuerystring: true,
      };
      if (this.keepalive) {
        requestOptions.forever = true;
      }


      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      const logger = this.logger;

      let retries = 0;
      (function requestOnce() {
        request(requestOptions, (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            setTimeout(requestOnce, backoff);
            return;
          }
          if (err) {
            err._fromRequest = true;
            responseLog(logger, requestOptions, response, err)
            reject(err);
            return;
          }

          switch (response.statusCode) {
            case 200:
              resolve(body);
              break;

            case 400:
              var err = new Errors.BadRequest(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 500:
              var err = new Errors.InternalError(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            default:
              var err = new Error("Received unexpected statusCode " + response.statusCode);
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;
          }
        });
      }());
    });
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
   * @param {module:wag/samples.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @returns {Object} iter
   * @returns {function} iter.map - takes in a function, applies it to each resource, and returns a promise to the result as an array
   * @returns {function} iter.toArray - returns a promise to the resources as an array
   * @returns {function} iter.forEach - takes in a function, applies it to each resource
   * @returns {function} iter.forEachAsync - takes in an async function, applies it to each resource
   */
  getBooksIter(params, options) {
    const it = (f, saveResults, isAsync) => new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;

      const headers = {};
      headers["Canonical-Resource"] = "getBooks";
      headers[versionHeader] = version;
      headers["authorization"] = params.authorization;

      const query = {};
      if (typeof params.authors !== "undefined") {
        query["authors"] = params.authors;
      }

      if (typeof params.available !== "undefined") {
        query["available"] = params.available;
      }

      if (typeof params.state !== "undefined") {
        query["state"] = params.state;
      }

      if (typeof params.published !== "undefined") {
        query["published"] = params.published;
      }

      if (typeof params.snakeCase !== "undefined") {
        query["snake_case"] = params.snakeCase;
      }

      if (typeof params.completed !== "undefined") {
        query["completed"] = params.completed;
      }

      if (typeof params.maxPages !== "undefined") {
        query["maxPages"] = params.maxPages;
      }

      if (typeof params.minPages !== "undefined") {
        query["min_pages"] = params.minPages;
      }

      if (typeof params.pagesToTime !== "undefined") {
        query["pagesToTime"] = params.pagesToTime;
      }

      if (typeof params.startingAfter !== "undefined") {
        query["startingAfter"] = params.startingAfter;
      }


      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/books",
        gzip: true,
        json: true,
        timeout,
        headers,
        qs: query,
        useQuerystring: true,
      };
      if (this.keepalive) {
        requestOptions.forever = true;
      }


      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      const logger = this.logger;

      let results = [];
      async.whilst(
        () => requestOptions.uri !== "",
        cbW => {
      const address = this.address;
      let retries = 0;
      (function requestOnce() {
        request(requestOptions, async (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            setTimeout(requestOnce, backoff);
            return;
          }
          if (err) {
            err._fromRequest = true;
            responseLog(logger, requestOptions, response, err)
            cbW(err);
            return;
          }

          switch (response.statusCode) {
            case 200:
              if (saveResults) {
                results = results.concat(body.map(f));
              } else {
                if (isAsync) {
                  for (let i = 0; i < body.length; i++) {
                    try {
                      await f(body[i], i, body);
                    } catch(err) {
                      reject(err);
                    }
                  }
                } else {
                  body.forEach(f)
                }
              }
              break;

            case 400:
              var err = new Errors.BadRequest(body || {});
              responseLog(logger, requestOptions, response, err);
              cbW(err);
              return;

            case 500:
              var err = new Errors.InternalError(body || {});
              responseLog(logger, requestOptions, response, err);
              cbW(err);
              return;

            default:
              var err = new Error("Received unexpected statusCode " + response.statusCode);
              responseLog(logger, requestOptions, response, err);
              cbW(err);
              return;
          }

          requestOptions.qs = null;
          requestOptions.useQuerystring = false;
          requestOptions.uri = "";
          if (response.headers["x-next-page-path"]) {
            requestOptions.uri = address + response.headers["x-next-page-path"];
          }
          cbW();
        });
      }());
        },
        err => {
          if (err) {
            reject(err);
            return;
          }
          if (saveResults) {
            resolve(results);
          } else {
            resolve();
          }
        }
      );
    });

    return {
      map: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, true, false]), cb),
      toArray: cb => applyCallback(this._hystrixCommand.execute(it, [x => x, true, false]), cb),
      forEach: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, false]), cb),
      forEachAsync: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, true]), cb),
    };
  }

  /**
   * Creates a book
   * @param newBook
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {module:wag/samples.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:wag/samples.Errors.BadRequest}
   * @reject {module:wag/samples.Errors.InternalError}
   * @reject {Error}
   */
  createBook(newBook, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._createBook, arguments), callback);
  }

  _createBook(newBook, options, cb) {
    const params = {};
    params["newBook"] = newBook;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;

      const headers = {};
      headers["Canonical-Resource"] = "createBook";
      headers[versionHeader] = version;

      const query = {};

      const requestOptions = {
        method: "POST",
        uri: this.address + "/v1/books",
        gzip: true,
        json: true,
        timeout,
        headers,
        qs: query,
        useQuerystring: true,
      };
      if (this.keepalive) {
        requestOptions.forever = true;
      }

      requestOptions.body = params.newBook;


      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      const logger = this.logger;

      let retries = 0;
      (function requestOnce() {
        request(requestOptions, (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            setTimeout(requestOnce, backoff);
            return;
          }
          if (err) {
            err._fromRequest = true;
            responseLog(logger, requestOptions, response, err)
            reject(err);
            return;
          }

          switch (response.statusCode) {
            case 200:
              resolve(body);
              break;

            case 400:
              var err = new Errors.BadRequest(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 500:
              var err = new Errors.InternalError(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            default:
              var err = new Error("Received unexpected statusCode " + response.statusCode);
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;
          }
        });
      }());
    });
  }

  /**
   * Puts a book
   * @param newBook
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {module:wag/samples.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:wag/samples.Errors.BadRequest}
   * @reject {module:wag/samples.Errors.InternalError}
   * @reject {Error}
   */
  putBook(newBook, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._putBook, arguments), callback);
  }

  _putBook(newBook, options, cb) {
    const params = {};
    params["newBook"] = newBook;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;

      const headers = {};
      headers["Canonical-Resource"] = "putBook";
      headers[versionHeader] = version;

      const query = {};

      const requestOptions = {
        method: "PUT",
        uri: this.address + "/v1/books",
        gzip: true,
        json: true,
        timeout,
        headers,
        qs: query,
        useQuerystring: true,
      };
      if (this.keepalive) {
        requestOptions.forever = true;
      }

      requestOptions.body = params.newBook;


      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      const logger = this.logger;

      let retries = 0;
      (function requestOnce() {
        request(requestOptions, (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            setTimeout(requestOnce, backoff);
            return;
          }
          if (err) {
            err._fromRequest = true;
            responseLog(logger, requestOptions, response, err)
            reject(err);
            return;
          }

          switch (response.statusCode) {
            case 200:
              resolve(body);
              break;

            case 400:
              var err = new Errors.BadRequest(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 500:
              var err = new Errors.InternalError(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            default:
              var err = new Error("Received unexpected statusCode " + response.statusCode);
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;
          }
        });
      }());
    });
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
   * @param {module:wag/samples.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:wag/samples.Errors.BadRequest}
   * @reject {module:wag/samples.Errors.Unathorized}
   * @reject {module:wag/samples.Errors.Error}
   * @reject {module:wag/samples.Errors.InternalError}
   * @reject {Error}
   */
  getBookByID(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getBookByID, arguments), callback);
  }

  _getBookByID(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;

      const headers = {};
      headers["Canonical-Resource"] = "getBookByID";
      headers[versionHeader] = version;
      if (!params.bookID) {
        reject(new Error("bookID must be non-empty because it's a path parameter"));
        return;
      }
      headers["authorization"] = params.authorization;
      headers["X-Dont-Rate-Limit-Me-Bro"] = params.XDontRateLimitMeBro;

      const query = {};
      if (typeof params.authorID !== "undefined") {
        query["authorID"] = params.authorID;
      }

      if (typeof params.randomBytes !== "undefined") {
        query["randomBytes"] = params.randomBytes;
      }


      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/books/" + params.bookID + "",
        gzip: true,
        json: true,
        timeout,
        headers,
        qs: query,
        useQuerystring: true,
      };
      if (this.keepalive) {
        requestOptions.forever = true;
      }


      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      const logger = this.logger;

      let retries = 0;
      (function requestOnce() {
        request(requestOptions, (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            setTimeout(requestOnce, backoff);
            return;
          }
          if (err) {
            err._fromRequest = true;
            responseLog(logger, requestOptions, response, err)
            reject(err);
            return;
          }

          switch (response.statusCode) {
            case 200:
              resolve(body);
              break;

            case 400:
              var err = new Errors.BadRequest(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 401:
              var err = new Errors.Unathorized(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 404:
              var err = new Errors.Error(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 500:
              var err = new Errors.InternalError(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            default:
              var err = new Error("Received unexpected statusCode " + response.statusCode);
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;
          }
        });
      }());
    });
  }

  /**
   * Retrieve a book
   * @param {string} id
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {module:wag/samples.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:wag/samples.Errors.BadRequest}
   * @reject {module:wag/samples.Errors.Error}
   * @reject {module:wag/samples.Errors.InternalError}
   * @reject {Error}
   */
  getBookByID2(id, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getBookByID2, arguments), callback);
  }

  _getBookByID2(id, options, cb) {
    const params = {};
    params["id"] = id;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;

      const headers = {};
      headers["Canonical-Resource"] = "getBookByID2";
      headers[versionHeader] = version;
      if (!params.id) {
        reject(new Error("id must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/books2/" + params.id + "",
        gzip: true,
        json: true,
        timeout,
        headers,
        qs: query,
        useQuerystring: true,
      };
      if (this.keepalive) {
        requestOptions.forever = true;
      }


      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      const logger = this.logger;

      let retries = 0;
      (function requestOnce() {
        request(requestOptions, (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            setTimeout(requestOnce, backoff);
            return;
          }
          if (err) {
            err._fromRequest = true;
            responseLog(logger, requestOptions, response, err)
            reject(err);
            return;
          }

          switch (response.statusCode) {
            case 200:
              resolve(body);
              break;

            case 400:
              var err = new Errors.BadRequest(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 404:
              var err = new Errors.Error(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 500:
              var err = new Errors.InternalError(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            default:
              var err = new Error("Received unexpected statusCode " + response.statusCode);
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;
          }
        });
      }());
    });
  }

  /**
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {module:wag/samples.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:wag/samples.Errors.BadRequest}
   * @reject {module:wag/samples.Errors.InternalError}
   * @reject {Error}
   */
  healthCheck(options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._healthCheck, arguments), callback);
  }

  _healthCheck(options, cb) {
    const params = {};

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;

      const headers = {};
      headers["Canonical-Resource"] = "healthCheck";
      headers[versionHeader] = version;

      const query = {};

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/health/check",
        gzip: true,
        json: true,
        timeout,
        headers,
        qs: query,
        useQuerystring: true,
      };
      if (this.keepalive) {
        requestOptions.forever = true;
      }


      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      const logger = this.logger;

      let retries = 0;
      (function requestOnce() {
        request(requestOptions, (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            setTimeout(requestOnce, backoff);
            return;
          }
          if (err) {
            err._fromRequest = true;
            responseLog(logger, requestOptions, response, err)
            reject(err);
            return;
          }

          switch (response.statusCode) {
            case 200:
              resolve();
              break;

            case 400:
              var err = new Errors.BadRequest(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 500:
              var err = new Errors.InternalError(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            default:
              var err = new Error("Received unexpected statusCode " + response.statusCode);
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;
          }
        });
      }());
    });
  }
};

module.exports = WagSamples;

/**
 * Retry policies available to use.
 * @alias module:wag/samples.RetryPolicies
 */
module.exports.RetryPolicies = {
  Single: singleRetryPolicy,
  Exponential: exponentialRetryPolicy,
  None: noRetryPolicy,
};

/**
 * Errors returned by methods.
 * @alias module:wag/samples.Errors
 */
module.exports.Errors = Errors;

module.exports.DefaultCircuitOptions = defaultCircuitOptions;

const version = "9.0.0";
const versionHeader = "X-Client-Version";
module.exports.Version = version;
module.exports.VersionHeader = versionHeader;
