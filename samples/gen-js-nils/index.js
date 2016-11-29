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
 * @alias module:nil-test.RetryPolicies.Exponential
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
 * @alias module:nil-test.RetryPolicies.Single
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
 * @alias module:nil-test.RetryPolicies.None
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
 * nil-test client library.
 * @module nil-test
 * @typicalname NilTest
 */

/**
 * nil-test client
 * @alias module:nil-test
 */
class NilTest {

  /**
   * Create a new client object.
   * @param {Object} options - Options for constructing a client object.
   * @param {string} [options.address] - URL where the server is located. Must provide
   * this or the discovery argument
   * @param {bool} [options.discovery] - Use clever-discovery to locate the server. Must provide
   * this or the address argument
   * @param {number} [options.timeout] - The timeout to use for all client requests,
   * in milliseconds. This can be overridden on a per-request basis.
   * @param {module:nil-test.RetryPolicies} [options.retryPolicy=RetryPolicies.Single] - The logic to
   * determine which requests to retry, as well as how many times to retry.
   */
  constructor(options) {
    options = options || {};

    if (options.discovery) {
      try {
        this.address = discovery("nil-test", "http").url();
      } catch (e) {
        this.address = discovery("nil-test", "default").url();
      }
    } else if (options.address) {
      this.address = options.address;
    } else {
      throw new Error("Cannot initialize nil-test without discovery or address");
    }
    if (options.timeout) {
      this.timeout = options.timeout;
    }
    if (options.retryPolicy) {
      this.retryPolicy = options.retryPolicy;
    }
  }

  /**
   * Nil check tests
   * @param {Object} params
   * @param {string} params.id
   * @param {string} [params.query]
   * @param {string} [params.header]
   * @param {string[]} [params.array]
   * @param [params.body]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:nil-test.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:nil-test.Errors.BadRequest}
   * @reject {module:nil-test.Errors.InternalError}
   * @reject {Error}
   */
  nilCheck(params, options, cb) {
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
    headers["header"] = params.header;

    const query = {};
    if (typeof params.query !== "undefined") {
      query["query"] = params.query;
    }

    if (typeof params.array !== "undefined") {
      query["array"] = params.array;
    }


    if (span) {
      opentracing.inject(span, opentracing.FORMAT_TEXT_MAP, headers);
      span.logEvent("POST /v1/check/{id}");
    }

    const requestOptions = {
      method: "POST",
      uri: this.address + "/v1/check/" + params.id + "",
      json: true,
      timeout,
      headers,
      qs: query,
      useQuerystring: true,
    };

    requestOptions.body = params.body;

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

module.exports = NilTest;

/**
 * Retry policies available to use.
 * @alias module:nil-test.RetryPolicies
 */
module.exports.RetryPolicies = {
  Single: singleRetryPolicy,
  Exponential: exponentialRetryPolicy,
  None: noRetryPolicy,
};

/**
 * Errors returned by methods.
 * @alias module:nil-test.Errors
 */
module.exports.Errors = Errors;
