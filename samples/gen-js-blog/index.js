const async = require("async");
const discovery = require("clever-discovery");
const kayvee = require("kayvee");
const request = require("request");
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
 * @private
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
	if (logData.status_code <= 499){
		logger.warnD("client-request-finished", logData);
	}else{
		logger.errorD("client-request-finished", logData);
	}
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
   * forever: true attribute in request. Defaults to true.
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
        this.address = discovery(options.serviceName || "blog", "http").url();
      } catch (e) {
        this.address = discovery(options.serviceName || "blog", "default").url();
      }
    } else if (options.address) {
      this.address = options.address;
    } else {
      throw new Error("Cannot initialize blog without discovery or address");
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
      this.logger = new kayvee.logger((options.serviceName || "blog") + "-wagclient");
    }

    const circuitOptions = Object.assign({}, defaultCircuitOptions, options.circuit);
    // hystrix implements a caching mechanism, we don't want this or we can't trust that clients
    // are initialized with the values passed in. 
    commandFactory.resetCache();
    circuitFactory.resetCache();
    metricsFactory.resetCache();
    this._hystrixCommand = commandFactory.getOrCreate(options.serviceName || "blog").
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
   * Posts the grade file for the specified student
   * @param {Object} params
   * @param {string} params.studentID
   * @param [params.file]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:blog.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:blog.Errors.BadRequest}
   * @reject {module:blog.Errors.InternalError}
   * @reject {Error}
   */
  postGradeFileForStudent(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._postGradeFileForStudent, arguments), callback);
  }

  _postGradeFileForStudent(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }
  
      const optionsBaggage = options.baggage || new Map();

      const timeout = options.timeout || this.timeout;

      let headers = {};
      
      // Convert optionsBaggage into a string using parseForBaggage
      headers["baggage"] = parseForBaggage(optionsBaggage);
      
      headers["Canonical-Resource"] = "postGradeFileForStudent";
      headers[versionHeader] = version;
      if (!params.studentID) {
        reject(new Error("studentID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      const requestOptions = {
        method: "POST",
        uri: this.address + "/students/" + params.studentID + "/gradeFile",
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

      requestOptions.body = params.file;


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

  /**
   * Gets the sections for the specified student
   * @param {string} studentID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:blog.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object[]}
   * @reject {module:blog.Errors.BadRequest}
   * @reject {module:blog.Errors.InternalError}
   * @reject {Error}
   */
  getSectionsForStudent(studentID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getSectionsForStudent, arguments), callback);
  }

  _getSectionsForStudent(studentID, options, cb) {
    const params = {};
    params["studentID"] = studentID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }
  
      const optionsBaggage = options.baggage || new Map();

      const timeout = options.timeout || this.timeout;

      let headers = {};
      
      // Convert optionsBaggage into a string using parseForBaggage
      headers["baggage"] = parseForBaggage(optionsBaggage);
      
      headers["Canonical-Resource"] = "getSectionsForStudent";
      headers[versionHeader] = version;
      if (!params.studentID) {
        reject(new Error("studentID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      const requestOptions = {
        method: "GET",
        uri: this.address + "/students/" + params.studentID + "/sections",
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
   * Posts the sections for the specified student
   * @param {Object} params
   * @param {string} params.studentID
   * @param {string} params.sections
   * @param {string} params.userType
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:blog.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object[]}
   * @reject {module:blog.Errors.BadRequest}
   * @reject {module:blog.Errors.InternalError}
   * @reject {Error}
   */
  postSectionsForStudent(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._postSectionsForStudent, arguments), callback);
  }

  _postSectionsForStudent(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }
  
      const optionsBaggage = options.baggage || new Map();

      const timeout = options.timeout || this.timeout;

      let headers = {};
      
      // Convert optionsBaggage into a string using parseForBaggage
      headers["baggage"] = parseForBaggage(optionsBaggage);
      
      headers["Canonical-Resource"] = "postSectionsForStudent";
      headers[versionHeader] = version;
      if (!params.studentID) {
        reject(new Error("studentID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};
      query["sections"] = params.sections;

      query["userType"] = params.userType;


      const requestOptions = {
        method: "POST",
        uri: this.address + "/students/" + params.studentID + "/sections",
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

const version = "9.0.0";
const versionHeader = "X-Client-Version";
module.exports.Version = version;
module.exports.VersionHeader = versionHeader;
