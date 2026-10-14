package jsclient

// The templates in this file generate the default JS client. The client sends requests
// through sendRequest, a wrapper around Node's built-in fetch that reproduces the behavior
// of the request library the client previously used.

var wrappedIndexJSTmplStr = `const async = require("async");
const discovery = require("clever-discovery");
const kayvee = require("kayvee");
const querystring = require("querystring");
const {commandFactory, circuitFactory, metricsFactory} = require("hystrixjs");
const RollingNumberEvent = require("hystrixjs/lib/metrics/RollingNumberEvent");

const { Errors } = require("./types");

function parseForBaggage(entries) {
  if (!entries) {
    return "";
  }
  // Regular expression for valid characters in keys and values
  const validChars = /^[a-zA-Z0-9!#$%&'*+` + "`\\-.^_`" + `|~]+$/;

  const pairs = [];

  entries.forEach((value, key) => {
    const validKey = key.match(validChars) ? key : encodeURIComponent(key);
    const validValue = value.match(validChars) ? value : encodeURIComponent(value);
    pairs.push(` + "`" + `${validKey}=${validValue}` + "`" + `);
  });

  return pairs.join(",");
}

/**
 * The exponential retry policy will retry five times with an exponential backoff.
 * @alias module:{{.ServiceName}}.RetryPolicies.Exponential
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
 * @alias module:{{.ServiceName}}.RetryPolicies.Single
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
 * @alias module:{{.ServiceName}}.RetryPolicies.None
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
	"backend": "{{.ServiceName}}",
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
 * Converts an error thrown by fetch into the error shape returned by
 * Node's http module, e.g. a connection error with a code like ECONNREFUSED.
 * @private
 */
function toRequestError(err, response) {
  if (err && err.name === "AbortError") {
    const timeoutErr = new Error("ESOCKETTIMEDOUT");
    timeoutErr.code = "ESOCKETTIMEDOUT";
    timeoutErr.connect = false;
    return timeoutErr;
  }
  const cause = err && err.cause instanceof Error ? err.cause : err;
  if (cause && cause.code === "UND_ERR_SOCKET") {
    // The server closed the connection. Node's http module reports this as "socket hang up"
    // before the response headers arrive, and as "aborted" while reading the response body.
    const resetErr = new Error(response ? "aborted" : "socket hang up");
    resetErr.code = "ECONNRESET";
    return resetErr;
  }
  return cause;
}

// Redirects are not followed for these methods, so the 3xx response is returned to the caller.
const noRedirectMethods = ["POST", "PUT", "PATCH", "DELETE"];

/**
 * Sends a JSON request using fetch. The callback is called with (err, response, body),
 * where response contains the statusCode, statusMessage, and lowercased headers, and body
 * is the parsed JSON response body (or the raw text if it is not valid JSON).
 * requestOptions.timeout is an idle timeout: the request fails if the server sends nothing
 * for that many milliseconds while connecting, waiting for headers, or reading the body.
 * @private
 */
function sendRequest(requestOptions, cb) {
  let url = requestOptions.uri;
  if (requestOptions.qs) {
    const qs = querystring.stringify(requestOptions.qs);
    if (qs) {
      url += (url.includes("?") ? "&" : "?") + qs;
    }
  }

  const headers = new Headers();
  Object.keys(requestOptions.headers || {}).forEach((name) => {
    const value = requestOptions.headers[name];
    if (value !== undefined) {
      try {
        headers.set(name, String(value));
      } catch (err) {
        // Matches the old behavior of request: Node's http module rejected invalid header
        // values with code ERR_INVALID_CHAR, which callers check for, but fetch's error has no code.
        const headerErr = new TypeError('Invalid character in header content ["' + name + '"]');
        headerErr.code = "ERR_INVALID_CHAR";
        throw headerErr;
      }
    }
  });
  if (!headers.has("accept")) {
    headers.set("accept", "application/json");
  }

  const controller = new AbortController();
  const init = {
    method: requestOptions.method,
    headers,
    signal: controller.signal,
    redirect: noRedirectMethods.includes(requestOptions.method) ? "manual" : "follow",
  };
  if (requestOptions.body !== undefined) {
    init.body = JSON.stringify(requestOptions.body);
    if (!headers.has("content-type")) {
      headers.set("content-type", "application/json");
    }
  }

  let timer;
  const resetTimer = () => {
    if (requestOptions.timeout) {
      clearTimeout(timer);
      timer = setTimeout(() => controller.abort(), requestOptions.timeout);
    }
  };

  let response;
  (async () => {
    resetTimer();
    const res = await fetch(url, init);
    response = {
      statusCode: res.status,
      statusMessage: res.statusText,
      headers: Object.fromEntries(res.headers.entries()),
    };
    const chunks = [];
    if (res.body) {
      const reader = res.body.getReader();
      for (;;) {
        resetTimer();
        const { done, value } = await reader.read();
        if (done) {
          break;
        }
        chunks.push(value);
      }
    }
    return Buffer.concat(chunks).toString("utf8");
  })().then((text) => {
    clearTimeout(timer);
    let body;
    if (text) {
      try {
        body = JSON.parse(text);
      } catch (e) {
        body = text;
      }
    }
    cb(null, response, body);
  }, (err) => {
    clearTimeout(timer);
    cb(toRequestError(err, response), response);
  });
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
 * @alias module:{{.ServiceName}}.DefaultCircuitOptions
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
 * {{.ServiceName}} client library.
 * @module {{.ServiceName}}
 * @typicalname {{.ClassName}}
 */

/**
 * {{.ServiceName}} client
 * @alias module:{{.ServiceName}}
 */
class {{.ClassName}} {

  /**
   * Create a new client object.
   * @param {Object} options - Options for constructing a client object.
   * @param {string} [options.address] - URL where the server is located. Must provide
   * this or the discovery argument
   * @param {bool} [options.discovery] - Use clever-discovery to locate the server. Must provide
   * this or the address argument
   * @param {number} [options.timeout] - The timeout to use for all client requests,
   * in milliseconds. This can be overridden on a per-request basis. Default is 5000ms.
   * @param {bool} [options.keepalive] - Deprecated: has no effect. Connections are always kept alive
   * by Node's built-in fetch.
   * @param {module:{{.ServiceName}}.RetryPolicies} [options.retryPolicy=RetryPolicies.Single] - The logic to
   * determine which requests to retry, as well as how many times to retry.
   * @param {module:kayvee.Logger} [options.logger=logger.New("{{.ServiceName}}-wagclient")] - The Kayvee
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
   * @param {object} [options.asynclocalstore] a request scoped async store 
   */
  constructor(options) {
    options = options || {};

    if (options.discovery) {
      try {
        this.address = discovery(options.serviceName || "{{.ServiceName}}", "http").url();
      } catch (e) {
        this.address = discovery(options.serviceName || "{{.ServiceName}}", "default").url();
      }
    } else if (options.address) {
      this.address = options.address;
    } else {
      throw new Error("Cannot initialize {{.ServiceName}} without discovery or address");
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
      this.logger = new kayvee.logger((options.serviceName || "{{.ServiceName}}") + "-wagclient");
    }
    if (options.asynclocalstore) {
      this.asynclocalstore = options.asynclocalstore;
    }


    const circuitOptions = Object.assign({}, defaultCircuitOptions, options.circuit);
    // hystrix implements a caching mechanism, we don't want this or we can't trust that clients
    // are initialized with the values passed in. 
    commandFactory.resetCache();
    circuitFactory.resetCache();
    metricsFactory.resetCache();
    this._hystrixCommand = commandFactory.getOrCreate(options.serviceName || "{{.ServiceName}}").
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
    // to avoid counting 4XXs as errors, only count an error if it comes from sending the request
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
    this.logger.infoD("{{.ServiceName}}", {
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
{{range $methodCode := .Methods}}{{$methodCode}}{{end}}};

module.exports = {{.ClassName}};

/**
 * Retry policies available to use.
 * @alias module:{{.ServiceName}}.RetryPolicies
 */
module.exports.RetryPolicies = {
  Single: singleRetryPolicy,
  Exponential: exponentialRetryPolicy,
  None: noRetryPolicy,
};

/**
 * Errors returned by methods.
 * @alias module:{{.ServiceName}}.Errors
 */
module.exports.Errors = Errors;

module.exports.DefaultCircuitOptions = defaultCircuitOptions;

const version = "{{.Version}}";
const versionHeader = "X-Client-Version";
module.exports.Version = version;
module.exports.VersionHeader = versionHeader;
`

const wrappedPackageJSONTmplStr = `{
  "name": "{{.PackageName}}",
  "version": "{{.Version}}",
  "description": "{{.Description}}",
  "main": "index.js",
  "dependencies": {
    "async": "^2.1.4",
    "clever-discovery": "0.0.8",
    "kayvee": "^3.13.0",
    "hystrixjs": "^0.2.0",
    "rxjs": "^5.4.1"
  },
  "engines": {
    "node": ">=18"
  },
  "devDependencies": {
    "typescript": "^3.3.0"
  }
}
`

const wrappedMethodTmplStr = `
  {{.MethodDefinition}}
    {{if .IterMethod -}}
    const it = (f, saveResults, isAsync) => new Promise((resolve, reject) => {
    {{- else -}}
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
    {{- end}}
      if (!options) {
        options = {};
      }
  
      const optionsBaggage = options.baggage || new Map();

      const storeContext = this.asynclocalstore?.get("context") || new Map();

      const combinedContext = new Map([...storeContext, ...optionsBaggage]);

      const timeout = options.timeout || this.timeout;

      let headers = {};

      // Merge custom headers from options if provided
      headers = {...(options.headers || {})};

      // Convert combinedContext into a string using parseForBaggage
      headers["baggage"] = parseForBaggage(combinedContext);

      headers["Canonical-Resource"] = "{{.Operation}}";
      headers[versionHeader] = version;
      {{- range $param := .PathParams}}
      if (!params.{{$param.JSName}}) {
        reject(new Error("{{$param.JSName}} must be non-empty because it's a path parameter"));
        return;
      }
      {{- end -}}
      {{- range $param := .HeaderParams}}
      headers["{{$param.WagName}}"] = params.{{$param.JSName}};
      {{- end}}

      const query = {};
      {{- range $param := .QueryParams -}}
      {{- if $param.Required }}
      query["{{$param.WagName}}"] = params.{{$param.JSName}};
{{else}}
      if (typeof params.{{$param.JSName}} !== "undefined") {
        query["{{$param.WagName}}"] = params.{{$param.JSName}};
      }
{{end}}{{end}}

      const requestOptions = {
        method: "{{.Method}}",
        uri: this.address + "{{.PathCode}}",
        timeout,
        headers,
        qs: query,
      };
{{ if ne .BodyParam ""}}
      requestOptions.body = params.{{.BodyParam}};
{{ end }}

      const retryPolicy = options.retryPolicy || this.retryPolicy || singleRetryPolicy;
      const backoffs = retryPolicy.backoffs();
      const logger = this.logger;
{{if .IterMethod}}
      let results = [];
      async.whilst(
        () => requestOptions.uri !== "",
        cbW => {
      const address = this.address;
  {{- end}}
      let retries = 0;
      (function requestOnce() {
        sendRequest(requestOptions, {{if $.IterMethod}}async {{end}}(err, response, body) => {
          if (retries < backoffs.length && retryPolicy.retry(requestOptions, err, response, body)) {
            const backoff = backoffs[retries];
            retries += 1;
            setTimeout(requestOnce, backoff);
            return;
          }
          if (err) {
            err._fromRequest = true;
            responseLog(logger, requestOptions, response, err)
            {{- if not .IterMethod}}
            reject(err);
            {{- else}}
            cbW(err);
            {{- end}}
            return;
          }

          switch (response.statusCode) {
            {{ range $response := .Responses }}case {{ $response.StatusCode }}:{{if $response.IsError }}
              var err = new Errors.{{ $response.Name }}(body || {});
              responseLog(logger, requestOptions, response, err);
              {{- if not $.IterMethod}}
              reject(err);
              {{- else}}
              cbW(err);
              {{- end}}
              return;
{{else}}{{if $response.IsNoData}}
              resolve();
              break;
{{else}}
              {{if $.IterMethod -}}
              if (saveResults) {
                results = results.concat(body{{$.IterResourceAccessString}}.map(f));
              } else {
                if (isAsync) {
                  for (let i = 0; i < body{{$.IterResourceAccessString}}.length; i++) {
                    try {
                      await f(body{{$.IterResourceAccessString}}[i], i, body);
                    } catch(err) {
                      reject(err);
                    }
                  }
                } else {
                  body{{$.IterResourceAccessString}}.forEach(f)
                }
              }
              {{- else -}}
              resolve(body);
              {{- end}}
              break;
{{end}}{{end}}
            {{end}}default:
              var err = new Error("Received unexpected statusCode " + response.statusCode);
              responseLog(logger, requestOptions, response, err);
              {{- if not .IterMethod}}
              reject(err);
              {{- else}}
              cbW(err);
              {{- end}}
              return;
          }

          {{- if .IterMethod}}

          requestOptions.qs = null;
          requestOptions.uri = "";
          if (response.headers["x-next-page-path"]) {
            requestOptions.uri = address + response.headers["x-next-page-path"];
          }
          cbW();
          {{- end}}
        });
      }());

      {{- if .IterMethod}}
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
      {{- end}}
    });

    {{- if .IterMethod}}

    return {
      map: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, true, false]), cb),
      toArray: cb => applyCallback(this._hystrixCommand.execute(it, [x => x, true, false]), cb),
      forEach: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, false]), cb),
      forEachAsync: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, true]), cb),
    };
    {{- end}}
  }
`
