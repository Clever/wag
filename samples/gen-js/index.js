const async = require("async");
const discovery = require("clever-discovery");
const request = require("request");
const opentracing = require("opentracing");

/**
 * @external Span
 * @see {@link https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html}
 */

const { Errors } = require("./types");

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
 * @alias module:swagger-test.RetryPolicies.Single
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
   * in milliseconds. This can be overridden on a per-request basis.
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy=RetryPolicies.Single] - The logic to
   * determine which requests to retry, as well as how many times to retry.
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
   * Gets authors
   * @param {Object} params
   * @param {string} [params.name]
   * @param {string} [params.startingAfter]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:swagger-test.Errors.BadRequest}
   * @reject {module:swagger-test.Errors.InternalError}
   * @reject {Error}
   */
  getAuthors(params, options, cb) {
    if (!cb && typeof options === "function") {
      cb = options;
      options = undefined;
    }

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


      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const span = options.span;

      const headers = {};

      const query = {};
      if (typeof params.name !== "undefined") {
        query["name"] = params.name;
      }
  
      if (typeof params.startingAfter !== "undefined") {
        query["startingAfter"] = params.startingAfter;
      }
  

      if (span) {
        opentracing.inject(span, opentracing.FORMAT_TEXT_MAP, headers);
        span.logEvent("GET /v1/authors");
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/authors",
        json: true,
        timeout,
        headers,
        qs: query,
        useQuerystring: true,
      };
  

      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
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

          let e;
          switch (response.statusCode) {
            case 200:
              resolver(body);
              break;
            
            case 400:
              e = new Errors.BadRequest(body || {});
              rejecter(e);
              return;
            
            case 500:
              e = new Errors.InternalError(body || {});
              rejecter(e);
              return;
            
            default:
              e = new Error("Recieved unexpected statusCode " + response.statusCode);
              rejecter(e);
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
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Object} iter
   * @returns {function} iter.map - takes in a function, applies it to each resource, and returns a promise to the result as an array
   * @returns {function} iter.toArray - returns a promise to the resources as an array
   * @returns {function} iter.forEach - takes in a function, applies it to each resource
   */
  getAuthorsIter(params, options, cb) {
    if (!cb && typeof options === "function") {
      cb = options;
      options = undefined;
    }

    const it = (f, saveResults) => new Promise((resolve, reject) => {
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


      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const span = options.span;

      const headers = {};

      const query = {};
      if (typeof params.name !== "undefined") {
        query["name"] = params.name;
      }
  
      if (typeof params.startingAfter !== "undefined") {
        query["startingAfter"] = params.startingAfter;
      }
  

      if (span) {
        opentracing.inject(span, opentracing.FORMAT_TEXT_MAP, headers);
        span.logEvent("GET /v1/authors");
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/authors",
        json: true,
        timeout,
        headers,
        qs: query,
        useQuerystring: true,
      };
  

      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
      const backoffs = retryPolicy.backoffs();
  
      let results = [];
      async.whilst(
        () => requestOptions.uri !== "",
        cbW => {
      const address = this.address;
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
            cbW(err);
            return;
          }

          let e;
          switch (response.statusCode) {
            case 200:
              if (saveResults) {
                results = results.concat(body.authorSet.results.map(f));
              } else {
                body.authorSet.results.forEach(f);
              }
              break;
            
            case 400:
              e = new Errors.BadRequest(body || {});
              rejecter(e);
              cbW(e);
              return;
            
            case 500:
              e = new Errors.InternalError(body || {});
              rejecter(e);
              cbW(e);
              return;
            
            default:
              e = new Error("Recieved unexpected statusCode " + response.statusCode);
              rejecter(e);
              cbW(e);
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
          if (!err) {
            if (saveResults) {
              resolver(results);
            } else {
              resolver();
            }
          }
        }
      );
    });

    return {
      map: f => it(f, true),
      toArray: () => it(x => x, true),
      forEach: f => it(f, false),
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
   * @param {number} [params.startingAfter]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
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
  
      if (typeof params.startingAfter !== "undefined") {
        query["startingAfter"] = params.startingAfter;
      }
  

      if (span) {
        opentracing.inject(span, opentracing.FORMAT_TEXT_MAP, headers);
        span.logEvent("GET /v1/books");
        span.setTag("span.kind", "client");
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
  

      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
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

          let e;
          switch (response.statusCode) {
            case 200:
              resolver(body);
              break;
            
            case 400:
              e = new Errors.BadRequest(body || {});
              rejecter(e);
              return;
            
            case 500:
              e = new Errors.InternalError(body || {});
              rejecter(e);
              return;
            
            default:
              e = new Error("Recieved unexpected statusCode " + response.statusCode);
              rejecter(e);
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
   * @param {number} [params.startingAfter]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Object} iter
   * @returns {function} iter.map - takes in a function, applies it to each resource, and returns a promise to the result as an array
   * @returns {function} iter.toArray - returns a promise to the resources as an array
   * @returns {function} iter.forEach - takes in a function, applies it to each resource
   */
  getBooksIter(params, options, cb) {
    if (!cb && typeof options === "function") {
      cb = options;
      options = undefined;
    }

    const it = (f, saveResults) => new Promise((resolve, reject) => {
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
  
      if (typeof params.startingAfter !== "undefined") {
        query["startingAfter"] = params.startingAfter;
      }
  

      if (span) {
        opentracing.inject(span, opentracing.FORMAT_TEXT_MAP, headers);
        span.logEvent("GET /v1/books");
        span.setTag("span.kind", "client");
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
  

      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
      const backoffs = retryPolicy.backoffs();
  
      let results = [];
      async.whilst(
        () => requestOptions.uri !== "",
        cbW => {
      const address = this.address;
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
            cbW(err);
            return;
          }

          let e;
          switch (response.statusCode) {
            case 200:
              if (saveResults) {
                results = results.concat(body.map(f));
              } else {
                body.forEach(f);
              }
              break;
            
            case 400:
              e = new Errors.BadRequest(body || {});
              rejecter(e);
              cbW(e);
              return;
            
            case 500:
              e = new Errors.InternalError(body || {});
              rejecter(e);
              cbW(e);
              return;
            
            default:
              e = new Error("Recieved unexpected statusCode " + response.statusCode);
              rejecter(e);
              cbW(e);
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
          if (!err) {
            if (saveResults) {
              resolver(results);
            } else {
              resolver();
            }
          }
        }
      );
    });

    return {
      map: f => it(f, true),
      toArray: () => it(x => x, true),
      forEach: f => it(f, false),
    };
  }

  /**
   * Creates a book
   * @param newBook
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
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
        span.setTag("span.kind", "client");
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
  

      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
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

          let e;
          switch (response.statusCode) {
            case 200:
              resolver(body);
              break;
            
            case 400:
              e = new Errors.BadRequest(body || {});
              rejecter(e);
              return;
            
            case 500:
              e = new Errors.InternalError(body || {});
              rejecter(e);
              return;
            
            default:
              e = new Error("Recieved unexpected statusCode " + response.statusCode);
              rejecter(e);
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
   * @param {string} [params.randomBytes]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
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


      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const span = options.span;

      const headers = {};
      if (!params.bookID) {
        rejecter(new Error("bookID must be non-empty because it's a path parameter"));
        return;
      }
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
        span.setTag("span.kind", "client");
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
  

      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
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

          let e;
          switch (response.statusCode) {
            case 200:
              resolver(body);
              break;
            
            case 400:
              e = new Errors.BadRequest(body || {});
              rejecter(e);
              return;
            
            case 401:
              e = new Errors.Unathorized(body || {});
              rejecter(e);
              return;
            
            case 404:
              e = new Errors.Error(body || {});
              rejecter(e);
              return;
            
            case 500:
              e = new Errors.InternalError(body || {});
              rejecter(e);
              return;
            
            default:
              e = new Error("Recieved unexpected statusCode " + response.statusCode);
              rejecter(e);
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
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
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


      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const span = options.span;

      const headers = {};
      if (!params.id) {
        rejecter(new Error("id must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span) {
        opentracing.inject(span, opentracing.FORMAT_TEXT_MAP, headers);
        span.logEvent("GET /v1/books2/{id}");
        span.setTag("span.kind", "client");
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
  

      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
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

          let e;
          switch (response.statusCode) {
            case 200:
              resolver(body);
              break;
            
            case 400:
              e = new Errors.BadRequest(body || {});
              rejecter(e);
              return;
            
            case 404:
              e = new Errors.Error(body || {});
              rejecter(e);
              return;
            
            case 500:
              e = new Errors.InternalError(body || {});
              rejecter(e);
              return;
            
            default:
              e = new Error("Recieved unexpected statusCode " + response.statusCode);
              rejecter(e);
              return;
          }
        });
      }());
    });
  }

  /**
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:swagger-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
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
        span.setTag("span.kind", "client");
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
  

      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
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

          let e;
          switch (response.statusCode) {
            case 200:
              resolver();
              break;
            
            case 400:
              e = new Errors.BadRequest(body || {});
              rejecter(e);
              return;
            
            case 500:
              e = new Errors.InternalError(body || {});
              rejecter(e);
              return;
            
            default:
              e = new Error("Recieved unexpected statusCode " + response.statusCode);
              rejecter(e);
              return;
          }
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
  Single: singleRetryPolicy,
  Exponential: exponentialRetryPolicy,
  None: noRetryPolicy,
};

/**
 * Errors returned by methods.
 * @alias module:swagger-test.Errors
 */
module.exports.Errors = Errors;
