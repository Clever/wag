const discovery = require("@clever/discovery");
const request = require("request");
const opentracing = require("opentracing");

const { Errors } = require("./types");

/**
 * The default retry policy will retry five times with an exponential backoff.
 * @alias module:swagger-test.RetryPolicies.Default
 */
const defaultRetryPolicy = {
  /**
   * backoffs returns an array of five backoffs: 100ms, 200ms, 400ms, 800ms, and
   * 1.6s. It adds a random 5% jitter to each backoff.
   * @function
   * @returns {Array<number>}
   */
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
  /**
   * retry will not retry a request if the HTTP client returns an error, if the
   * is a POST or PATCH, or if the status code is less than 500. It will retry
   * all other requests.
   * @function
   * @returns {boolean}
   */
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
 * @alias module:swagger-test.RetryPolicies.None
 */
const noRetryPolicy = {
  /**
   * returns an empty array
   */
  backoffs() {
    return [];
  },
  /**
   * returns false
   */
  retry() {
    return false;
  },
};

/**
 * swagger-test client library.
 * @module swagger-test
 * @typicalname SwaggerTest
 */

/**
 * The main client object to instantiate.
 * @alias module:swagger-test
 */
class SwaggerTest {

  /**
   * Create a new client object.
   * @param {Object} options - Options for constructing a client object.
   * @param {string} options.address - URL where the server is located. If not
   * specified, the address will be discovered via @clever/discovery.
   * @param {number} options.timeout - The timeout to use for all client requests,
   * in milliseconds. This can be overridden on a per-request basis.
   * @param {Object} [options.retryPolicy=RetryPolicies.Default] - The logic to
   * determine which requests to retry, as well as how many times to retry.
   * @param {function} options.retryPolicy.backoffs
   * @param {function} options.retryPolicy.retry
   */
  constructor(options) {
    options = options || {};

    if (options.discovery) {
      try {
        this.address = discovery("swagger-test", "http").url();
      } catch (e) {
        this.address = discovery("swagger-test", "default").url();
      }
    } else if (options.address) {
      this.address = options.address;
    } else {
      throw new Error("Cannot initialize swagger-test without discovery or address");
    }
    if (options.timeout) {
      this.timeout = options.timeout;
    }
    if (options.retryPolicy) {
      this.retryPolicy = options.retryPolicy;
    }
  }

  /**
   * Returns a list of books
   * @param {Object} params
   * @param {string[]} [params.authors]
   * @param {boolean} [params.available=true]
   * @param {string} [params.state=finished]
   * @param {string} [params.published]
   * @param {string} [params.snakeCase]
   * @param {string} [params.completed]
   * @param {number} [params.maxPages=500.5]
   * @param {number} [params.minPages=5]
   * @param {number} [params.pagesToTime]
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object[]}
   * @reject {module:swagger-test.Errors.BadRequest}
   * @reject {module:swagger-test.Errors.InternalError}
   * @reject {Error}
   */
  getBooks(params, options, cb) {
    if (!cb && typeof options === "function") {
      cb = options;
      options = undefined;
    }

    if (!options) {
      options = {};
    }

    const timeout = options.timeout || this.timeout;
    const span = options.span;

    const headers = {};

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


    if (span) {
      opentracing.inject(span, opentracing.FORMAT_TEXT_MAP, headers);
      span.logEvent("GET /v1/books");
    }

    const requestOptions = {
      method: "GET",
      uri: this.address + "/v1/books",
      json: true,
      timeout,
      headers,
      qs: query,
      useQuerystring: true,
    };

    return new Promise((resolve, reject) => {
      const rejecter = (err) => {
        reject(err);
        if (cb) {
          cb(err);
        }
      };
      const resolver = (data) => {
        resolve(data);
        if (cb) {
          cb(null, data);
        }
      };

      const retryPolicy = options.retryPolicy || this.retryPolicy || defaultRetryPolicy;
      const backoffs = retryPolicy.backoffs();
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
            rejecter(err);
            return;
          }
          switch (response.statusCode) {
            case 200:
              resolver(body);
              break;
            case 400:
              rejecter(new Errors.BadRequest(body || {}));
              break;
            case 500:
              rejecter(new Errors.InternalError(body || {}));
              break;
            default:
              rejecter(new Error("Recieved unexpected statusCode " + response.statusCode));
          }
          return;
        });
      }());
    });
  }

  /**
   * Creates a book
   * @param newBook
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:swagger-test.Errors.BadRequest}
   * @reject {module:swagger-test.Errors.InternalError}
   * @reject {Error}
   */
  createBook(newBook, options, cb) {
    const params = {};
    params["newBook"] = newBook;

    if (!cb && typeof options === "function") {
      cb = options;
      options = undefined;
    }

    if (!options) {
      options = {};
    }

    const timeout = options.timeout || this.timeout;
    const span = options.span;

    const headers = {};

    const query = {};

    if (span) {
      opentracing.inject(span, opentracing.FORMAT_TEXT_MAP, headers);
      span.logEvent("POST /v1/books");
    }

    const requestOptions = {
      method: "POST",
      uri: this.address + "/v1/books",
      json: true,
      timeout,
      headers,
      qs: query,
      useQuerystring: true,
    };

    requestOptions.body = params.newBook;

    return new Promise((resolve, reject) => {
      const rejecter = (err) => {
        reject(err);
        if (cb) {
          cb(err);
        }
      };
      const resolver = (data) => {
        resolve(data);
        if (cb) {
          cb(null, data);
        }
      };

      const retryPolicy = options.retryPolicy || this.retryPolicy || defaultRetryPolicy;
      const backoffs = retryPolicy.backoffs();
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
            rejecter(err);
            return;
          }
          switch (response.statusCode) {
            case 200:
              resolver(body);
              break;
            case 400:
              rejecter(new Errors.BadRequest(body || {}));
              break;
            case 500:
              rejecter(new Errors.InternalError(body || {}));
              break;
            default:
              rejecter(new Error("Recieved unexpected statusCode " + response.statusCode));
          }
          return;
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
   * @param {string} [params.randomBytes]
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
    if (!cb && typeof options === "function") {
      cb = options;
      options = undefined;
    }

    if (!options) {
      options = {};
    }

    const timeout = options.timeout || this.timeout;
    const span = options.span;

    const headers = {};
    headers["authorization"] = params.authorization;

    const query = {};
    if (typeof params.authorID !== "undefined") {
      query["authorID"] = params.authorID;
    }

    if (typeof params.randomBytes !== "undefined") {
      query["randomBytes"] = params.randomBytes;
    }


    if (span) {
      opentracing.inject(span, opentracing.FORMAT_TEXT_MAP, headers);
      span.logEvent("GET /v1/books/{book_id}");
    }

    const requestOptions = {
      method: "GET",
      uri: this.address + "/v1/books/" + params.bookID + "",
      json: true,
      timeout,
      headers,
      qs: query,
      useQuerystring: true,
    };

    return new Promise((resolve, reject) => {
      const rejecter = (err) => {
        reject(err);
        if (cb) {
          cb(err);
        }
      };
      const resolver = (data) => {
        resolve(data);
        if (cb) {
          cb(null, data);
        }
      };

      const retryPolicy = options.retryPolicy || this.retryPolicy || defaultRetryPolicy;
      const backoffs = retryPolicy.backoffs();
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
            rejecter(err);
            return;
          }
          switch (response.statusCode) {
            case 200:
              resolver(body);
              break;
            case 400:
              rejecter(new Errors.BadRequest(body || {}));
              break;
            case 401:
              rejecter(new Errors.Unathorized(body || {}));
              break;
            case 404:
              rejecter(new Errors.Error(body || {}));
              break;
            case 500:
              rejecter(new Errors.InternalError(body || {}));
              break;
            default:
              rejecter(new Error("Recieved unexpected statusCode " + response.statusCode));
          }
          return;
        });
      }());
    });
  }

  /**
   * Retrieve a book
   * @param {string} id
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:swagger-test.Errors.BadRequest}
   * @reject {module:swagger-test.Errors.Error}
   * @reject {module:swagger-test.Errors.InternalError}
   * @reject {Error}
   */
  getBookByID2(id, options, cb) {
    const params = {};
    params["id"] = id;

    if (!cb && typeof options === "function") {
      cb = options;
      options = undefined;
    }

    if (!options) {
      options = {};
    }

    const timeout = options.timeout || this.timeout;
    const span = options.span;

    const headers = {};

    const query = {};

    if (span) {
      opentracing.inject(span, opentracing.FORMAT_TEXT_MAP, headers);
      span.logEvent("GET /v1/books2/{id}");
    }

    const requestOptions = {
      method: "GET",
      uri: this.address + "/v1/books2/" + params.id + "",
      json: true,
      timeout,
      headers,
      qs: query,
      useQuerystring: true,
    };

    return new Promise((resolve, reject) => {
      const rejecter = (err) => {
        reject(err);
        if (cb) {
          cb(err);
        }
      };
      const resolver = (data) => {
        resolve(data);
        if (cb) {
          cb(null, data);
        }
      };

      const retryPolicy = options.retryPolicy || this.retryPolicy || defaultRetryPolicy;
      const backoffs = retryPolicy.backoffs();
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
            rejecter(err);
            return;
          }
          switch (response.statusCode) {
            case 200:
              resolver(body);
              break;
            case 400:
              rejecter(new Errors.BadRequest(body || {}));
              break;
            case 404:
              rejecter(new Errors.Error(body || {}));
              break;
            case 500:
              rejecter(new Errors.InternalError(body || {}));
              break;
            default:
              rejecter(new Error("Recieved unexpected statusCode " + response.statusCode));
          }
          return;
        });
      }());
    });
  }

  /**
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:swagger-test.Errors.BadRequest}
   * @reject {module:swagger-test.Errors.InternalError}
   * @reject {Error}
   */
  healthCheck(options, cb) {
    const params = {};

    if (!cb && typeof options === "function") {
      cb = options;
      options = undefined;
    }

    if (!options) {
      options = {};
    }

    const timeout = options.timeout || this.timeout;
    const span = options.span;

    const headers = {};

    const query = {};

    if (span) {
      opentracing.inject(span, opentracing.FORMAT_TEXT_MAP, headers);
      span.logEvent("GET /v1/health/check");
    }

    const requestOptions = {
      method: "GET",
      uri: this.address + "/v1/health/check",
      json: true,
      timeout,
      headers,
      qs: query,
      useQuerystring: true,
    };

    return new Promise((resolve, reject) => {
      const rejecter = (err) => {
        reject(err);
        if (cb) {
          cb(err);
        }
      };
      const resolver = (data) => {
        resolve(data);
        if (cb) {
          cb(null, data);
        }
      };

      const retryPolicy = options.retryPolicy || this.retryPolicy || defaultRetryPolicy;
      const backoffs = retryPolicy.backoffs();
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
            rejecter(err);
            return;
          }
          switch (response.statusCode) {
            case 200:
              resolver();
              break;
            case 400:
              rejecter(new Errors.BadRequest(body || {}));
              break;
            case 500:
              rejecter(new Errors.InternalError(body || {}));
              break;
            default:
              rejecter(new Error("Recieved unexpected statusCode " + response.statusCode));
          }
          return;
        });
      }());
    });
  }
};

module.exports = SwaggerTest;

/**
 * Retry policies available to use.
 * @alias module:swagger-test.RetryPolicies
 */
module.exports.RetryPolicies = {
  Default: defaultRetryPolicy,
  None: noRetryPolicy,
};

/**
 * Errors returned by methods.
 * @alias module:swagger-test.Errors
 */
module.exports.Errors = Errors;
