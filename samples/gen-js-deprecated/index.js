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
