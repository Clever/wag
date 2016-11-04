const discovery = require("@clever/discovery");
const request = require("request");
const url = require("url");
const opentracing = require("opentracing");

// go-swagger treats handles/expects arrays in the query string to be a string of comma joined values
// so...do that thing. It's worth noting that this has lots of issues ("what if my values have commas in them?")
// but that's an issue with go-swagger
function serializeQueryString(data) {
  if (Array.isArray(data)) {
    return data.join(",");
  }
  return data;
}

const defaultRetryPolicy = {
  backoffs() {
    const ret = [];
    let next = 100.0; // milliseconds
    const e = 0.05; // +/- 5% jitter
    while (ret.length < 5) {
      const jitter = (Math.random()*2-1.0)*e*next;
      ret.push(next + jitter);
      next *= 2;
    }
    return ret;
  },
  retry(requestOptions, err, res, body) {
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
  retry(requestOptions, err, res, body) {
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
      };
    } else if (options.address) {
      this.address = options.address;
    } else {
      throw new Error("Cannot initialize swagger-test without discovery or address");
    }
    if (options.timeout) {
      this.timeout = options.timeout
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
    query["authors"] = serializeQueryString(params.authors);
    query["available"] = serializeQueryString(params.available);
    query["state"] = serializeQueryString(params.state);
    query["published"] = serializeQueryString(params.published);
    query["snake_case"] = serializeQueryString(params.snakeCase);
    query["completed"] = serializeQueryString(params.completed);
    query["maxPages"] = serializeQueryString(params.maxPages);
    query["min_pages"] = serializeQueryString(params.minPages);
    query["pagesToTime"] = serializeQueryString(params.pagesToTime);

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
    };

    return new Promise((resolve, reject) => {
      const rejecter = (err) => {
        reject(err);
        if (cb) {
          cb(err);
        }
      }
      const resolver = (data) => {
        resolve(data);
        if (cb) {
          cb(null, data);
        }
      }

      const retryPolicy = options.retryPolicy || this.retryPolicy || defaultRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      let retries = 0;
      (function requestOnce() {
        request(requestOptions, (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            return setTimeout(requestOnce, backoff);
          }
          if (err) {
            return rejecter(err);
          }
          if (response.statusCode >= 400) {
            return rejecter(new Error(body));
          }
          resolver(body);
        });
      })();
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
    };

    requestOptions.body = params.newBook;

    return new Promise((resolve, reject) => {
      const rejecter = (err) => {
        reject(err);
        if (cb) {
          cb(err);
        }
      }
      const resolver = (data) => {
        resolve(data);
        if (cb) {
          cb(null, data);
        }
      }

      const retryPolicy = options.retryPolicy || this.retryPolicy || defaultRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      let retries = 0;
      (function requestOnce() {
        request(requestOptions, (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            return setTimeout(requestOnce, backoff);
          }
          if (err) {
            return rejecter(err);
          }
          if (response.statusCode >= 400) {
            return rejecter(new Error(body));
          }
          resolver(body);
        });
      })();
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
    query["authorID"] = serializeQueryString(params.authorID);
    query["randomBytes"] = serializeQueryString(params.randomBytes);

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
    };

    return new Promise((resolve, reject) => {
      const rejecter = (err) => {
        reject(err);
        if (cb) {
          cb(err);
        }
      }
      const resolver = (data) => {
        resolve(data);
        if (cb) {
          cb(null, data);
        }
      }

      const retryPolicy = options.retryPolicy || this.retryPolicy || defaultRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      let retries = 0;
      (function requestOnce() {
        request(requestOptions, (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            return setTimeout(requestOnce, backoff);
          }
          if (err) {
            return rejecter(err);
          }
          if (response.statusCode >= 400) {
            return rejecter(new Error(body));
          }
          resolver(body);
        });
      })();
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
    };

    return new Promise((resolve, reject) => {
      const rejecter = (err) => {
        reject(err);
        if (cb) {
          cb(err);
        }
      }
      const resolver = (data) => {
        resolve(data);
        if (cb) {
          cb(null, data);
        }
      }

      const retryPolicy = options.retryPolicy || this.retryPolicy || defaultRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      let retries = 0;
      (function requestOnce() {
        request(requestOptions, (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            return setTimeout(requestOnce, backoff);
          }
          if (err) {
            return rejecter(err);
          }
          if (response.statusCode >= 400) {
            return rejecter(new Error(body));
          }
          resolver(body);
        });
      })();
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
    };

    return new Promise((resolve, reject) => {
      const rejecter = (err) => {
        reject(err);
        if (cb) {
          cb(err);
        }
      }
      const resolver = (data) => {
        resolve(data);
        if (cb) {
          cb(null, data);
        }
      }

      const retryPolicy = options.retryPolicy || this.retryPolicy || defaultRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      let retries = 0;
      (function requestOnce() {
        request(requestOptions, (err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            return setTimeout(requestOnce, backoff);
          }
          if (err) {
            return rejecter(err);
          }
          if (response.statusCode >= 400) {
            return rejecter(new Error(body));
          }
          resolver(body);
        });
      })();
    });
  }
}

module.exports.RetryPolicies = {
  Default: defaultRetryPolicy,
  None: noRetryPolicy,
};
