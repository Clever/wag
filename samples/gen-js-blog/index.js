const async = require("async");
const discovery = require("clever-discovery");
const kayvee = require("kayvee");
const request = require("request");
const opentracing = require("opentracing");
const globalTracing = require("opentracing/lib/global_tracer");
const {commandFactory} = require("hystrixjs");
const RollingNumberEvent = require("hystrixjs/lib/metrics/RollingNumberEvent");

/**
 * @external Span
 * @see {@link https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html}
 */

const { Errors } = require("./types");

/**
 * The exponential retry policy will retry five times with an exponential backoff.
 * @alias module:blog.RetryPolicies.Exponential
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
 * @alias module:blog.RetryPolicies.Single
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
 * @alias module:blog.RetryPolicies.None
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
 */
function responseLog(logger, req, res, err) {
  var res = res || { };
  var req = req || { };
  var logData = {
	"backend": "blog",
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
 * Default circuit breaker options.
 * @alias module:blog.DefaultCircuitOptions
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
 * blog client library.
 * @module blog
 * @typicalname Blog
 */

/**
 * blog client
 * @alias module:blog
 */
class Blog {

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
   * forever: true attribute in request. Defaults to false
   * @param {module:blog.RetryPolicies} [options.retryPolicy=RetryPolicies.Single] - The logic to
   * determine which requests to retry, as well as how many times to retry.
   * @param {module:kayvee.Logger} [options.logger=logger.New("blog-wagclient")] - The Kayvee
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
        this.address = discovery("blog", "http").url();
      } catch (e) {
        this.address = discovery("blog", "default").url();
      }
    } else if (options.address) {
      this.address = options.address;
    } else {
      throw new Error("Cannot initialize blog without discovery or address");
    }
    if (options.keepalive) {
      this.keepalive = options.keepalive;
    } else {
      this.keepalive = false;
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
      this.logger =  new kayvee.logger("blog-wagclient");
    }

    const circuitOptions = Object.assign({}, defaultCircuitOptions, options.circuit);
    this._hystrixCommand = commandFactory.getOrCreate("blog").
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

    setInterval(() => this._logCircuitState(), circuitOptions.logIntervalMs);
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
    this.logger.infoD("blog", {
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
   * Gets the sections for the specified student
   * @param {string} studentID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:blog.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object[]}
   * @reject {module:blog.Errors.BadRequest}
   * @reject {module:blog.Errors.InternalError}
   * @reject {Error}
   */
  getSectionsForStudent(studentID, options, cb) {
    return this._hystrixCommand.execute(this._getSectionsForStudent, arguments);
  }
  _getSectionsForStudent(studentID, options, cb) {
    const params = {};
    params["studentID"] = studentID;

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
      if (!params.studentID) {
        rejecter(new Error("studentID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span) {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        globalTracing.globalTracer().inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.logEvent("GET /students/{student_id}/sections");
        span.setTag("span.kind", "client");
      }

	  const requestOptions = {
        method: "GET",
        uri: this.address + "/students/" + params.studentID + "/sections",
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
            rejecter(err);
            return;
          }

          switch (response.statusCode) {
            case 200:
              resolver(body);
              break;
            
            case 400:
              var err = new Errors.BadRequest(body || {});
              responseLog(logger, requestOptions, response, err);
              rejecter(err);
              return;
            
            case 500:
              var err = new Errors.InternalError(body || {});
              responseLog(logger, requestOptions, response, err);
              rejecter(err);
              return;
            
            default:
              var err = new Error("Received unexpected statusCode " + response.statusCode);
              responseLog(logger, requestOptions, response, err);
              rejecter(err);
              return;
          }
        });
      }());
    });
  }
};

module.exports = Blog;

/**
 * Retry policies available to use.
 * @alias module:blog.RetryPolicies
 */
module.exports.RetryPolicies = {
  Single: singleRetryPolicy,
  Exponential: exponentialRetryPolicy,
  None: noRetryPolicy,
};

/**
 * Errors returned by methods.
 * @alias module:blog.Errors
 */
module.exports.Errors = Errors;

module.exports.DefaultCircuitOptions = defaultCircuitOptions;
