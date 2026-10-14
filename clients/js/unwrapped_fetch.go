package jsclient

// The templates in this file generate the JS client when the swagger.yml sets
// `x-unwrapped-fetch: true` in its info section. The client calls Node's built-in fetch
// directly, instead of going through a wrapper that reproduces the behavior of the
// request library. Its public interface (methods, options, errors, retry policies,
// circuit breaker) is the same as the default client, except that:
//   - retry policies receive the fetch Response and the response body text
//   - timeouts cover the entire request and reject with fetch's TimeoutError
//   - request errors are the errors thrown by fetch
//   - redirects are always followed
//   - path parameters are URI-encoded
//   - the keepalive option is not supported

var unwrappedIndexJSTmplStr = `const discovery = require("clever-discovery");
const kayvee = require("kayvee");
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
  retry(request, err, res) {
    return !err && request.method !== "POST" && request.method !== "PATCH" && res.status >= 500;
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
  retry(request, err, res) {
    return !err && request.method !== "POST" && request.method !== "PATCH" && res.status >= 500;
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
 * Logs a failed request. 5XX responses are logged as errors, and everything else as warnings.
 * @private
 */
function logRequestError(logger, request, response, err) {
  const status = response ? response.status : 0;
  const logData = {
    "backend": "{{.ServiceName}}",
    "method": request.method,
    "uri": request.url,
    "message": err,
    "status_code": status,
  };
  if (status >= 500) {
    logger.errorD("client-request-finished", logData);
  } else {
    logger.warnD("client-request-finished", logData);
  }
}

/**
 * Parses a JSON response body, returning undefined if the body is empty.
 * @private
 */
function parseJSON(body) {
  return body ? JSON.parse(body) : undefined;
}

/**
 * Adds a query parameter to a URL, skipping undefined values. Arrays are added
 * as repeated parameters, e.g. ?tag=a&tag=b.
 * @private
 */
function appendQueryParam(url, name, value) {
  if (value === undefined) {
    return;
  }
  for (const v of [].concat(value)) {
    url.searchParams.append(name, v);
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
    this.timeout = options.timeout || 5000;
    this.retryPolicy = options.retryPolicy || singleRetryPolicy;
    this.logger = options.logger || new kayvee.logger((options.serviceName || "{{.ServiceName}}") + "-wagclient");
    this.asynclocalstore = options.asynclocalstore;

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

  /**
   * Returns the headers sent with every request: the caller's headers, plus baggage
   * and the operation and client version.
   * @private
   */
  _headers(operation, options) {
    const storeContext = this.asynclocalstore?.get("context") || new Map();
    const baggage = new Map([...storeContext, ...(options.baggage || new Map())]);

    const headers = new Headers(options.headers);
    headers.set("baggage", parseForBaggage(baggage));
    headers.set("Canonical-Resource", operation);
    headers.set(versionHeader, version);
    return headers;
  }

  /**
   * Sends a request with fetch, retrying it according to the retry policy. Resolves with the
   * response and its body as text. The timeout covers the entire request, including reading
   * the body.
   * @private
   */
  async _send(request, options) {
    const timeout = options.timeout || this.timeout;
    const retryPolicy = options.retryPolicy || this.retryPolicy;
    const backoffs = retryPolicy.backoffs();

    for (let retries = 0; ; retries++) {
      let response, body, err;
      try {
        response = await fetch(request.url, {
          method: request.method,
          headers: request.headers,
          body: request.body,
          signal: AbortSignal.timeout(timeout),
        });
        body = await response.text();
      } catch (e) {
        err = e;
      }

      if (retries < backoffs.length && retryPolicy.retry(request, err, response, body)) {
        await new Promise((resolve) => setTimeout(resolve, backoffs[retries]));
        continue;
      }
      if (err) {
        err._fromRequest = true;
        logRequestError(this.logger, request, response, err);
        throw err;
      }
      return { response, body };
    }
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

var unwrappedPackageJSONTmplStr = `{
  "name": "{{.PackageName}}",
  "version": "{{.Version}}",
  "description": "{{.Description}}",
  "main": "index.js",
  "dependencies": {
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

// unwrappedMethodTmplStr generates the body of a method. MethodDefinition declares the
// method, and for non-paging methods, the async private method that sends the request.
// Paging methods are split the same way: the public method returns the iterator, and a
// private async method fetches each page.
var unwrappedMethodTmplStr = `
  {{.MethodDefinition}}
  {{- if .IterMethod}}
    const it = (f, saveResults, isAsync) => this._{{.MethodName}}(params, options, f, saveResults, isAsync);
    return {
      map: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, true, false]), cb),
      toArray: cb => applyCallback(this._hystrixCommand.execute(it, [x => x, true, false]), cb),
      forEach: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, false]), cb),
      forEachAsync: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false, true]), cb),
    };
  }

  async _{{.MethodName}}(params, options, f, saveResults, isAsync) {
    options = options || {};
  {{- else}}
    if (!options || typeof options === "function") {
      options = {};
    }
  {{- end}}
  {{- range $param := .PathParams}}
    if (!params.{{$param.JSName}}) {
      throw new Error("{{$param.JSName}} must be non-empty because it's a path parameter");
    }
  {{- end}}

    const headers = this._headers("{{.Operation}}", options);
  {{- range $param := .HeaderParams}}
    if (params.{{$param.JSName}} !== undefined) {
      headers.set("{{$param.WagName}}", params.{{$param.JSName}});
    }
  {{- end}}
  {{- if ne .BodyParam ""}}
    headers.set("Content-Type", "application/json");
  {{- end}}

    const url = new URL(this.address + "{{.PathCode}}");
  {{- range $param := .QueryParams}}
    appendQueryParam(url, "{{$param.WagName}}", params.{{$param.JSName}});
  {{- end}}

    const request = {
      method: "{{.Method}}",
      url: url.href,
      headers,
  {{- if ne .BodyParam ""}}
      body: JSON.stringify(params.{{.BodyParam}}),
  {{- end}}
    };
  {{- if .IterMethod}}

    let results = [];
    while (request.url) {
      const { response, body } = await this._send(request, options);
      switch (response.status) {
    {{- range $response := .Responses}}
        case {{$response.StatusCode}}: {
      {{- if $response.IsError}}
          const err = new Errors.{{$response.Name}}(parseJSON(body) || {});
          logRequestError(this.logger, request, response, err);
          throw err;
      {{- else if $response.IsNoData}}
          break;
      {{- else}}
          const page = parseJSON(body);
          const items = page{{$.IterResourceAccessString}};
          if (saveResults) {
            results = results.concat(items.map(f));
          } else if (isAsync) {
            for (let i = 0; i < items.length; i++) {
              await f(items[i], i, page);
            }
          } else {
            items.forEach(f);
          }
          break;
      {{- end}}
        }
    {{- end}}
        default: {
          const err = new Error("Received unexpected statusCode " + response.status);
          logRequestError(this.logger, request, response, err);
          throw err;
        }
      }

      const nextPagePath = response.headers.get("x-next-page-path");
      request.url = nextPagePath ? this.address + nextPagePath : "";
    }
    return saveResults ? results : undefined;
  {{- else}}

    const { response, body } = await this._send(request, options);
    switch (response.status) {
    {{- range $response := .Responses}}
      case {{$response.StatusCode}}: {
      {{- if $response.IsError}}
        const err = new Errors.{{$response.Name}}(parseJSON(body) || {});
        logRequestError(this.logger, request, response, err);
        throw err;
      {{- else if $response.IsNoData}}
        return;
      {{- else}}
        return parseJSON(body);
      {{- end}}
      }
    {{- end}}
      default: {
        const err = new Error("Received unexpected statusCode " + response.status);
        logRequestError(this.logger, request, response, err);
        throw err;
      }
    }
  {{- end}}
  }
`
