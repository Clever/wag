const async = require("async");
const discovery = require("clever-discovery");
const kayvee = require("kayvee");
const request = require("request");
const opentracing = require("opentracing");
const {commandFactory} = require("hystrixjs");
const RollingNumberEvent = require("hystrixjs/lib/metrics/RollingNumberEvent");

/**
 * @external Span
 * @see {@link https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html}
 */

const { Errors } = require("./types");

/**
 * The exponential retry policy will retry five times with an exponential backoff.
 * @alias module:app-service.RetryPolicies.Exponential
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
 * @alias module:app-service.RetryPolicies.Single
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
 * @alias module:app-service.RetryPolicies.None
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
	"backend": "app-service",
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
 * @alias module:app-service.DefaultCircuitOptions
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
 * app-service client library.
 * @module app-service
 * @typicalname AppService
 */

/**
 * app-service client
 * @alias module:app-service
 */
class AppService {

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
   * @param {module:app-service.RetryPolicies} [options.retryPolicy=RetryPolicies.Single] - The logic to
   * determine which requests to retry, as well as how many times to retry.
   * @param {module:kayvee.Logger} [options.logger=logger.New("app-service-wagclient")] - The Kayvee
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
        this.address = discovery(options.serviceName || "app-service", "http").url();
      } catch (e) {
        this.address = discovery(options.serviceName || "app-service", "default").url();
      }
    } else if (options.address) {
      this.address = options.address;
    } else {
      throw new Error("Cannot initialize app-service without discovery or address");
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
      this.logger = new kayvee.logger((options.serviceName || "app-service") + "-wagclient");
    }
    if (options.tracer) {
      this.tracer = options.tracer;
    } else {
      this.tracer = opentracing.globalTracer();
    }

    const circuitOptions = Object.assign({}, defaultCircuitOptions, options.circuit);
    this._hystrixCommand = commandFactory.getOrCreate(options.serviceName || "app-service").
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
    this.logger.infoD("app-service", {
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
   * Checks if the service is healthy
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.InternalError}
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
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "healthCheck";
      headers[versionHeader] = version;

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /_health"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/_health",
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

  /**
   * @param {Object} params
   * @param {string} [params.email]
   * @param {string} [params.password]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object[]}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getAdmins(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAdmins, arguments), callback);
  }

  _getAdmins(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getAdmins";
      headers[versionHeader] = version;

      const query = {};
      if (typeof params.email !== "undefined") {
        query["email"] = params.email;
      }

      if (typeof params.password !== "undefined") {
        query["password"] = params.password;
      }


      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/admins"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/admins",
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
   * @param {string} adminID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  deleteAdmin(adminID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._deleteAdmin, arguments), callback);
  }

  _deleteAdmin(adminID, options, cb) {
    const params = {};
    params["adminID"] = adminID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "deleteAdmin";
      headers[versionHeader] = version;
      if (!params.adminID) {
        reject(new Error("adminID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "DELETE /v1/admins/{adminID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "DELETE",
        uri: this.address + "/v1/admins/" + params.adminID + "",
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

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {string} adminID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getAdminByID(adminID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAdminByID, arguments), callback);
  }

  _getAdminByID(adminID, options, cb) {
    const params = {};
    params["adminID"] = adminID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getAdminByID";
      headers[versionHeader] = version;
      if (!params.adminID) {
        reject(new Error("adminID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/admins/{adminID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/admins/" + params.adminID + "",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.adminID
   * @param params.admin
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  updateAdmin(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._updateAdmin, arguments), callback);
  }

  _updateAdmin(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "updateAdmin";
      headers[versionHeader] = version;
      if (!params.adminID) {
        reject(new Error("adminID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PATCH /v1/admins/{adminID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PATCH",
        uri: this.address + "/v1/admins/" + params.adminID + "",
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

      requestOptions.body = params.admin;


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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param params.createAdmin
   * @param {string} params.adminID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  createAdmin(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._createAdmin, arguments), callback);
  }

  _createAdmin(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "createAdmin";
      headers[versionHeader] = version;
      if (!params.adminID) {
        reject(new Error("adminID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PUT /v1/admins/{adminID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PUT",
        uri: this.address + "/v1/admins/" + params.adminID + "",
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

      requestOptions.body = params.createAdmin;


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
   * @param {Object} params
   * @param {string} params.code
   * @param {boolean} [params.invalidate]
   * @param {string} params.adminID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  verifyCode(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._verifyCode, arguments), callback);
  }

  _verifyCode(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "verifyCode";
      headers[versionHeader] = version;
      if (!params.adminID) {
        reject(new Error("adminID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};
      query["code"] = params.code;

      if (typeof params.invalidate !== "undefined") {
        query["invalidate"] = params.invalidate;
      }


      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "POST /v1/admins/{adminID}/confirmation_code"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "POST",
        uri: this.address + "/v1/admins/" + params.adminID + "/confirmation_code",
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

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {number} params.duration
   * @param {string} params.adminID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  createVerificationCode(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._createVerificationCode, arguments), callback);
  }

  _createVerificationCode(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "createVerificationCode";
      headers[versionHeader] = version;
      if (!params.adminID) {
        reject(new Error("adminID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};
      query["duration"] = params.duration;


      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PUT /v1/admins/{adminID}/confirmation_code"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PUT",
        uri: this.address + "/v1/admins/" + params.adminID + "/confirmation_code",
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
              var err = new Errors.NotFound(body || {});
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
   * set the verified email of an admin
   * @param {Object} params
   * @param {string} params.adminID
   * @param params.request
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  verifyAdminEmail(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._verifyAdminEmail, arguments), callback);
  }

  _verifyAdminEmail(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "verifyAdminEmail";
      headers[versionHeader] = version;
      if (!params.adminID) {
        reject(new Error("adminID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "POST /v1/admins/{adminID}/verify_email"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "POST",
        uri: this.address + "/v1/admins/" + params.adminID + "/verify_email",
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

      requestOptions.body = params.request;


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

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getAllAnalyticsApps(options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAllAnalyticsApps, arguments), callback);
  }

  _getAllAnalyticsApps(options, cb) {
    const params = {};

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getAllAnalyticsApps";
      headers[versionHeader] = version;

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/analytics/apps"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/analytics/apps",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {string} shortname
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getAnalyticsAppByShortname(shortname, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAnalyticsAppByShortname, arguments), callback);
  }

  _getAnalyticsAppByShortname(shortname, options, cb) {
    const params = {};
    params["shortname"] = shortname;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getAnalyticsAppByShortname";
      headers[versionHeader] = version;
      if (!params.shortname) {
        reject(new Error("shortname must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/analytics/apps/{shortname}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/analytics/apps/" + params.shortname + "",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getAllTrackableApps(options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAllTrackableApps, arguments), callback);
  }

  _getAllTrackableApps(options, cb) {
    const params = {};

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getAllTrackableApps";
      headers[versionHeader] = version;

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/analytics/trackable_apps"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/analytics/trackable_apps",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getAnalyticsUsageUrls(options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAnalyticsUsageUrls, arguments), callback);
  }

  _getAnalyticsUsageUrls(options, cb) {
    const params = {};

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getAnalyticsUsageUrls";
      headers[versionHeader] = version;

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/analytics/usageUrls"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/analytics/usageUrls",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getAllUsageUrls(options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAllUsageUrls, arguments), callback);
  }

  _getAllUsageUrls(options, cb) {
    const params = {};

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getAllUsageUrls";
      headers[versionHeader] = version;

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/appUniverse/usageUrls"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/appUniverse/usageUrls",
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
              var err = new Errors.NotFound(body || {});
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
   * The server takes in the intersection of input parameters
   * @param {Object} params
   * @param {string[]} [params.ids]
   * @param {string} [params.clientId]
   * @param {string} [params.clientSecret]
   * @param {string} [params.shortname]
   * @param {string} [params.businessToken]
   * @param {string[]} [params.tags]
   * @param {string[]} [params.skipTags]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object[]}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getApps(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getApps, arguments), callback);
  }

  _getApps(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getApps";
      headers[versionHeader] = version;

      const query = {};
      if (typeof params.ids !== "undefined") {
        query["ids"] = params.ids;
      }

      if (typeof params.clientId !== "undefined") {
        query["clientId"] = params.clientId;
      }

      if (typeof params.clientSecret !== "undefined") {
        query["clientSecret"] = params.clientSecret;
      }

      if (typeof params.shortname !== "undefined") {
        query["shortname"] = params.shortname;
      }

      if (typeof params.businessToken !== "undefined") {
        query["businessToken"] = params.businessToken;
      }

      if (typeof params.tags !== "undefined") {
        query["tags"] = params.tags;
      }

      if (typeof params.skipTags !== "undefined") {
        query["skipTags"] = params.skipTags;
      }


      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/apps"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/apps",
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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  deleteApp(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._deleteApp, arguments), callback);
  }

  _deleteApp(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "deleteApp";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "DELETE /v1/apps/{appID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "DELETE",
        uri: this.address + "/v1/apps/" + params.appID + "",
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

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getAppByID(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAppByID, arguments), callback);
  }

  _getAppByID(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getAppByID";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/apps/{appID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/apps/" + params.appID + "",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param {boolean} [params.withSchemaPropagation] - If scopes change, then the app schema will be updated. This flag will propagate app schema updates to all connection schemas as well

   * @param params.app
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  updateApp(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._updateApp, arguments), callback);
  }

  _updateApp(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "updateApp";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};
      if (typeof params.withSchemaPropagation !== "undefined") {
        query["withSchemaPropagation"] = params.withSchemaPropagation;
      }


      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PATCH /v1/apps/{appID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PATCH",
        uri: this.address + "/v1/apps/" + params.appID + "",
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

      requestOptions.body = params.app;


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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param [params.app]
   * @param {string} params.appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  createApp(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._createApp, arguments), callback);
  }

  _createApp(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "createApp";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PUT /v1/apps/{appID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PUT",
        uri: this.address + "/v1/apps/" + params.appID + "",
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

      requestOptions.body = params.app;


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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object[]}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getAdminsForApp(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAdminsForApp, arguments), callback);
  }

  _getAdminsForApp(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getAdminsForApp";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/apps/{appID}/admins"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/apps/" + params.appID + "/admins",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param {string} params.adminID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.Forbidden}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  unlinkAppAdmin(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._unlinkAppAdmin, arguments), callback);
  }

  _unlinkAppAdmin(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "unlinkAppAdmin";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }
      if (!params.adminID) {
        reject(new Error("adminID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "DELETE /v1/apps/{appID}/admins/{adminID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "DELETE",
        uri: this.address + "/v1/apps/" + params.appID + "/admins/" + params.adminID + "",
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

            case 403:
              var err = new Errors.Forbidden(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param {string} params.adminID
   * @param params.permissions
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.Forbidden}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  linkAppAdmin(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._linkAppAdmin, arguments), callback);
  }

  _linkAppAdmin(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "linkAppAdmin";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }
      if (!params.adminID) {
        reject(new Error("adminID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PUT /v1/apps/{appID}/admins/{adminID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PUT",
        uri: this.address + "/v1/apps/" + params.appID + "/admins/" + params.adminID + "",
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

      requestOptions.body = params.permissions;


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

            case 403:
              var err = new Errors.Forbidden(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param {string} params.adminID
   * @param {string} params.guideID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.Forbidden}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getGuideConfig(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getGuideConfig, arguments), callback);
  }

  _getGuideConfig(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getGuideConfig";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }
      if (!params.adminID) {
        reject(new Error("adminID must be non-empty because it's a path parameter"));
        return;
      }
      if (!params.guideID) {
        reject(new Error("guideID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/apps/{appID}/admins/{adminID}/guides/{guideID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/apps/" + params.appID + "/admins/" + params.adminID + "/guides/" + params.guideID + "",
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

            case 403:
              var err = new Errors.Forbidden(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param {string} params.adminID
   * @param {string} params.guideID
   * @param params.guideConfig
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.Forbidden}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  setGuideConfig(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._setGuideConfig, arguments), callback);
  }

  _setGuideConfig(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "setGuideConfig";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }
      if (!params.adminID) {
        reject(new Error("adminID must be non-empty because it's a path parameter"));
        return;
      }
      if (!params.guideID) {
        reject(new Error("guideID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PUT /v1/apps/{appID}/admins/{adminID}/guides/{guideID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PUT",
        uri: this.address + "/v1/apps/" + params.appID + "/admins/" + params.adminID + "/guides/" + params.guideID + "",
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

      requestOptions.body = params.guideConfig;


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

            case 403:
              var err = new Errors.Forbidden(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.adminID
   * @param {string} params.appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getPermissionsForAdmin(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getPermissionsForAdmin, arguments), callback);
  }

  _getPermissionsForAdmin(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getPermissionsForAdmin";
      headers[versionHeader] = version;
      if (!params.adminID) {
        reject(new Error("adminID must be non-empty because it's a path parameter"));
        return;
      }
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/apps/{appID}/admins/{adminID}/permissions"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/apps/" + params.appID + "/admins/" + params.adminID + "/permissions",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param {string} params.adminID
   * @param {boolean} params.verified
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.Forbidden}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  verifyAppAdmin(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._verifyAppAdmin, arguments), callback);
  }

  _verifyAppAdmin(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "verifyAppAdmin";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }
      if (!params.adminID) {
        reject(new Error("adminID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};
      query["verified"] = params.verified;


      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "POST /v1/apps/{appID}/admins/{adminID}/verify"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "POST",
        uri: this.address + "/v1/apps/" + params.appID + "/admins/" + params.adminID + "/verify",
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

            case 403:
              var err = new Errors.Forbidden(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  generateNewBusinessToken(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._generateNewBusinessToken, arguments), callback);
  }

  _generateNewBusinessToken(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "generateNewBusinessToken";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "POST /v1/apps/{appID}/business_token"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "POST",
        uri: this.address + "/v1/apps/" + params.appID + "/business_token",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param {number} params.schoolYearStart
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getCertifications(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getCertifications, arguments), callback);
  }

  _getCertifications(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getCertifications";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }
      if (!params.schoolYearStart) {
        reject(new Error("schoolYearStart must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/apps/{appID}/certifications/{schoolYearStart}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/apps/" + params.appID + "/certifications/" + params.schoolYearStart + "",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param {number} params.schoolYearStart
   * @param params.certifications
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  setCertifications(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._setCertifications, arguments), callback);
  }

  _setCertifications(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "setCertifications";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }
      if (!params.schoolYearStart) {
        reject(new Error("schoolYearStart must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "POST /v1/apps/{appID}/certifications/{schoolYearStart}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "POST",
        uri: this.address + "/v1/apps/" + params.appID + "/certifications/" + params.schoolYearStart + "",
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

      requestOptions.body = params.certifications;


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
              var err = new Errors.NotFound(body || {});
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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getSetupStep(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getSetupStep, arguments), callback);
  }

  _getSetupStep(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getSetupStep";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/apps/{appID}/customStep"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/apps/" + params.appID + "/customStep",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param [params.setupStep]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  createSetupStep(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._createSetupStep, arguments), callback);
  }

  _createSetupStep(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "createSetupStep";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PATCH /v1/apps/{appID}/customStep"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PATCH",
        uri: this.address + "/v1/apps/" + params.appID + "/customStep",
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

      requestOptions.body = params.setupStep;


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

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object[]}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getDataRules(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getDataRules, arguments), callback);
  }

  _getDataRules(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getDataRules";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/apps/{appID}/data_rules"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/apps/" + params.appID + "/data_rules",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param [params.rules]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  setDataRules(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._setDataRules, arguments), callback);
  }

  _setDataRules(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "setDataRules";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PUT /v1/apps/{appID}/data_rules"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PUT",
        uri: this.address + "/v1/apps/" + params.appID + "/data_rules",
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

      requestOptions.body = params.rules;


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

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getManagers(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getManagers, arguments), callback);
  }

  _getManagers(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getManagers";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/apps/{appID}/managers"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/apps/" + params.appID + "/managers",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getOnboarding(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getOnboarding, arguments), callback);
  }

  _getOnboarding(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getOnboarding";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/apps/{appID}/onboarding"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/apps/" + params.appID + "/onboarding",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param params.update
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  updateOnboarding(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._updateOnboarding, arguments), callback);
  }

  _updateOnboarding(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "updateOnboarding";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PATCH /v1/apps/{appID}/onboarding"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PATCH",
        uri: this.address + "/v1/apps/" + params.appID + "/onboarding",
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

      requestOptions.body = params.update;


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

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  initializeOnboarding(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._initializeOnboarding, arguments), callback);
  }

  _initializeOnboarding(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "initializeOnboarding";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PUT /v1/apps/{appID}/onboarding"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PUT",
        uri: this.address + "/v1/apps/" + params.appID + "/onboarding",
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

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param {string} params.clientID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  deletePlatform(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._deletePlatform, arguments), callback);
  }

  _deletePlatform(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "deletePlatform";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }
      if (!params.clientID) {
        reject(new Error("clientID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "DELETE /v1/apps/{appID}/platform/{clientID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "DELETE",
        uri: this.address + "/v1/apps/" + params.appID + "/platform/" + params.clientID + "",
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

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param {string} params.clientID
   * @param params.request
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  updatePlatform(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._updatePlatform, arguments), callback);
  }

  _updatePlatform(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "updatePlatform";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }
      if (!params.clientID) {
        reject(new Error("clientID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PATCH /v1/apps/{appID}/platform/{clientID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PATCH",
        uri: this.address + "/v1/apps/" + params.appID + "/platform/" + params.clientID + "",
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

      requestOptions.body = params.request;


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
              var err = new Errors.NotFound(body || {});
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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object[]}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getPlatformsByAppID(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getPlatformsByAppID, arguments), callback);
  }

  _getPlatformsByAppID(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getPlatformsByAppID";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/apps/{appID}/platforms"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/apps/" + params.appID + "/platforms",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param params.request
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  createPlatform(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._createPlatform, arguments), callback);
  }

  _createPlatform(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "createPlatform";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PUT /v1/apps/{appID}/platforms"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PUT",
        uri: this.address + "/v1/apps/" + params.appID + "/platforms",
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

      requestOptions.body = params.request;


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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param {boolean} [params.deleteDataRules] - Delete field setting-style data warnings when app schema is deleted
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  deleteAppSchema(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._deleteAppSchema, arguments), callback);
  }

  _deleteAppSchema(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "deleteAppSchema";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};
      if (typeof params.deleteDataRules !== "undefined") {
        query["deleteDataRules"] = params.deleteDataRules;
      }


      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "DELETE /v1/apps/{appID}/schema"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "DELETE",
        uri: this.address + "/v1/apps/" + params.appID + "/schema",
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

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getAppSchema(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAppSchema, arguments), callback);
  }

  _getAppSchema(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getAppSchema";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/apps/{appID}/schema"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/apps/" + params.appID + "/schema",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param {boolean} [params.skipPropagation] - Skip propagation to connection schemas
   * @param {boolean} [params.updateDataRules] - Update data warnings when app schema changes
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  createAppSchema(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._createAppSchema, arguments), callback);
  }

  _createAppSchema(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "createAppSchema";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};
      if (typeof params.skipPropagation !== "undefined") {
        query["skipPropagation"] = params.skipPropagation;
      }

      if (typeof params.updateDataRules !== "undefined") {
        query["updateDataRules"] = params.updateDataRules;
      }


      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "POST /v1/apps/{appID}/schema"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "POST",
        uri: this.address + "/v1/apps/" + params.appID + "/schema",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param {boolean} [params.skipPropagation] - Skip propagation to connection schemas
   * @param {boolean} [params.updateDataRules] - Update data warnings when app schema changes
   * @param [params.appSchema]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  setAppSchema(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._setAppSchema, arguments), callback);
  }

  _setAppSchema(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "setAppSchema";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};
      if (typeof params.skipPropagation !== "undefined") {
        query["skipPropagation"] = params.skipPropagation;
      }

      if (typeof params.updateDataRules !== "undefined") {
        query["updateDataRules"] = params.updateDataRules;
      }


      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PUT /v1/apps/{appID}/schema"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PUT",
        uri: this.address + "/v1/apps/" + params.appID + "/schema",
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

      requestOptions.body = params.appSchema;


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
              var err = new Errors.NotFound(body || {});
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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getSecrets(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getSecrets, arguments), callback);
  }

  _getSecrets(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getSecrets";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/apps/{appID}/secrets"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/apps/" + params.appID + "/secrets",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  revokeOldClientSecret(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._revokeOldClientSecret, arguments), callback);
  }

  _revokeOldClientSecret(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "revokeOldClientSecret";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PATCH /v1/apps/{appID}/secrets"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PATCH",
        uri: this.address + "/v1/apps/" + params.appID + "/secrets",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  generateNewClientSecret(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._generateNewClientSecret, arguments), callback);
  }

  _generateNewClientSecret(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "generateNewClientSecret";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "POST /v1/apps/{appID}/secrets"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "POST",
        uri: this.address + "/v1/apps/" + params.appID + "/secrets",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  resetClientSecret(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._resetClientSecret, arguments), callback);
  }

  _resetClientSecret(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "resetClientSecret";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PUT /v1/apps/{appID}/secrets"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PUT",
        uri: this.address + "/v1/apps/" + params.appID + "/secrets",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {string} appID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getRecommendedSharing(appID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getRecommendedSharing, arguments), callback);
  }

  _getRecommendedSharing(appID, options, cb) {
    const params = {};
    params["appID"] = appID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getRecommendedSharing";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/apps/{appID}/sharing"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/apps/" + params.appID + "/sharing",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param [params.recommendations]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  setRecommendedSharing(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._setRecommendedSharing, arguments), callback);
  }

  _setRecommendedSharing(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "setRecommendedSharing";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PUT /v1/apps/{appID}/sharing"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PUT",
        uri: this.address + "/v1/apps/" + params.appID + "/sharing",
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

      requestOptions.body = params.recommendations;


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

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.appID
   * @param params.app
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.UnprocessableEntity}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  updateAppIcon(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._updateAppIcon, arguments), callback);
  }

  _updateAppIcon(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "updateAppIcon";
      headers[versionHeader] = version;
      if (!params.appID) {
        reject(new Error("appID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "POST /v1/apps/{appID}/update_icon"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "POST",
        uri: this.address + "/v1/apps/" + params.appID + "/update_icon",
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

      requestOptions.body = params.app;


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
              var err = new Errors.NotFound(body || {});
              responseLog(logger, requestOptions, response, err);
              reject(err);
              return;

            case 422:
              var err = new Errors.UnprocessableEntity(body || {});
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
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getAllCategories(options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAllCategories, arguments), callback);
  }

  _getAllCategories(options, cb) {
    const params = {};

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getAllCategories";
      headers[versionHeader] = version;

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/categories"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/categories",
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
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object[]}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getKnownHosts(options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getKnownHosts, arguments), callback);
  }

  _getKnownHosts(options, cb) {
    const params = {};

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getKnownHosts";
      headers[versionHeader] = version;

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/knownhosts"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/knownhosts",
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
   * @param {Object} params
   * @param {string} [params.category]
   * @param {boolean} [params.includeDevApps]
   * @param {boolean} [params.includeLinks]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getAllLibraryResources(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAllLibraryResources, arguments), callback);
  }

  _getAllLibraryResources(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getAllLibraryResources";
      headers[versionHeader] = version;

      const query = {};
      if (typeof params.category !== "undefined") {
        query["category"] = params.category;
      }

      if (typeof params.includeDevApps !== "undefined") {
        query["includeDevApps"] = params.includeDevApps;
      }

      if (typeof params.includeLinks !== "undefined") {
        query["includeLinks"] = params.includeLinks;
      }


      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/libraryResources"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/libraryResources",
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
   * @param {Object} params
   * @param {string} params.searchTerm
   * @param {boolean} [params.showInLibraryOnly]
   * @param {boolean} [params.includeLinks]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  searchLibraryResource(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._searchLibraryResource, arguments), callback);
  }

  _searchLibraryResource(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "searchLibraryResource";
      headers[versionHeader] = version;

      const query = {};
      query["searchTerm"] = params.searchTerm;

      if (typeof params.showInLibraryOnly !== "undefined") {
        query["showInLibraryOnly"] = params.showInLibraryOnly;
      }

      if (typeof params.includeLinks !== "undefined") {
        query["includeLinks"] = params.includeLinks;
      }


      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/libraryResources/search"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/libraryResources/search",
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
   * @param {Object} params
   * @param {string} params.shortname
   * @param {boolean} [params.includeDevApps]
   * @param {boolean} [params.includeLinks]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getLibraryResourceByShortname(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getLibraryResourceByShortname, arguments), callback);
  }

  _getLibraryResourceByShortname(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getLibraryResourceByShortname";
      headers[versionHeader] = version;
      if (!params.shortname) {
        reject(new Error("shortname must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};
      if (typeof params.includeDevApps !== "undefined") {
        query["includeDevApps"] = params.includeDevApps;
      }

      if (typeof params.includeLinks !== "undefined") {
        query["includeLinks"] = params.includeLinks;
      }


      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/libraryResources/{shortname}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/libraryResources/" + params.shortname + "",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.shortname
   * @param params.libraryResource
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  updateLibraryResourceByShortname(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._updateLibraryResourceByShortname, arguments), callback);
  }

  _updateLibraryResourceByShortname(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "updateLibraryResourceByShortname";
      headers[versionHeader] = version;
      if (!params.shortname) {
        reject(new Error("shortname must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "PATCH /v1/libraryResources/{shortname}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "PATCH",
        uri: this.address + "/v1/libraryResources/" + params.shortname + "",
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

      requestOptions.body = params.libraryResource;


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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.shortname
   * @param params.libraryResource
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  createLibraryResource(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._createLibraryResource, arguments), callback);
  }

  _createLibraryResource(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "createLibraryResource";
      headers[versionHeader] = version;
      if (!params.shortname) {
        reject(new Error("shortname must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "POST /v1/libraryResources/{shortname}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "POST",
        uri: this.address + "/v1/libraryResources/" + params.shortname + "",
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

      requestOptions.body = params.libraryResource;


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
              var err = new Errors.NotFound(body || {});
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
   * @param {string} shortname
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  deleteLibraryResourceLink(shortname, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._deleteLibraryResourceLink, arguments), callback);
  }

  _deleteLibraryResourceLink(shortname, options, cb) {
    const params = {};
    params["shortname"] = shortname;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "deleteLibraryResourceLink";
      headers[versionHeader] = version;
      if (!params.shortname) {
        reject(new Error("shortname must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "DELETE /v1/libraryResources/{shortname}/link"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "DELETE",
        uri: this.address + "/v1/libraryResources/" + params.shortname + "/link",
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

            case 404:
              var err = new Errors.NotFound(body || {});
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
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getValidPermissions(options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getValidPermissions, arguments), callback);
  }

  _getValidPermissions(options, cb) {
    const params = {};

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getValidPermissions";
      headers[versionHeader] = version;

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/permissions"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/permissions",
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
   * The server takes in the intersection of input parameters
   * @param {Object} params
   * @param {string[]} [params.appIds]
   * @param {string} [params.name]
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object[]}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getPlatforms(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getPlatforms, arguments), callback);
  }

  _getPlatforms(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getPlatforms";
      headers[versionHeader] = version;

      const query = {};
      if (typeof params.appIds !== "undefined") {
        query["appIds"] = params.appIds;
      }

      if (typeof params.name !== "undefined") {
        query["name"] = params.name;
      }


      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/platforms"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/platforms",
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
   * @param {string} clientID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getPlatformByClientID(clientID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getPlatformByClientID, arguments), callback);
  }

  _getPlatformByClientID(clientID, options, cb) {
    const params = {};
    params["clientID"] = clientID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getPlatformByClientID";
      headers[versionHeader] = version;
      if (!params.clientID) {
        reject(new Error("clientID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v1/platforms/{clientID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v1/platforms/" + params.clientID + "",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {string} adminID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {Object[]}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  getAppsForAdmin(adminID, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._getAppsForAdmin, arguments), callback);
  }

  _getAppsForAdmin(adminID, options, cb) {
    const params = {};
    params["adminID"] = adminID;

    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "getAppsForAdmin";
      headers[versionHeader] = version;
      if (!params.adminID) {
        reject(new Error("adminID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "GET /v2/admins/{adminID}/apps"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "GET",
        uri: this.address + "/v2/admins/" + params.adminID + "/apps",
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
              var err = new Errors.NotFound(body || {});
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
   * @param {Object} params
   * @param {string} params.srcAppID
   * @param {string} params.destAppID
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:app-service.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill {undefined}
   * @reject {module:app-service.Errors.BadRequest}
   * @reject {module:app-service.Errors.NotFound}
   * @reject {module:app-service.Errors.InternalError}
   * @reject {Error}
   */
  overrideConfig(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._overrideConfig, arguments), callback);
  }

  _overrideConfig(params, options, cb) {
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
      headers["Canonical-Resource"] = "overrideConfig";
      headers[versionHeader] = version;
      if (!params.srcAppID) {
        reject(new Error("srcAppID must be non-empty because it's a path parameter"));
        return;
      }
      if (!params.destAppID) {
        reject(new Error("destAppID must be non-empty because it's a path parameter"));
        return;
      }

      const query = {};

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        span.log({event: "POST /v2/apps/{srcAppID}/override-config/{destAppID}"});
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "POST",
        uri: this.address + "/v2/apps/" + params.srcAppID + "/override-config/" + params.destAppID + "",
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

            case 404:
              var err = new Errors.NotFound(body || {});
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

module.exports = AppService;

/**
 * Retry policies available to use.
 * @alias module:app-service.RetryPolicies
 */
module.exports.RetryPolicies = {
  Single: singleRetryPolicy,
  Exponential: exponentialRetryPolicy,
  None: noRetryPolicy,
};

/**
 * Errors returned by methods.
 * @alias module:app-service.Errors
 */
module.exports.Errors = Errors;

module.exports.DefaultCircuitOptions = defaultCircuitOptions;

const version = "12.15.1";
const versionHeader = "X-Client-Version";
module.exports.Version = version;
module.exports.VersionHeader = versionHeader;
