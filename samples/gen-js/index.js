const discovery = require("@clever/discovery");
const request = require("request");
const opentracing = require("opentracing");

const { Errors } = require("./types");

const defaultRetryPolicy = {
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

const noRetryPolicy = {
  backoffs() {
    return [];
  },
  retry() {
    return false;
  },
};

module.exports = class SwaggerTest {

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

module.exports.RetryPolicies = {
  Default: defaultRetryPolicy,
  None: noRetryPolicy,
};
module.exports.Errors = Errors;
