package jsclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Clever/go-utils/stringset"
	"github.com/Clever/wag/swagger"
	"github.com/Clever/wag/templates"
	"github.com/Clever/wag/utils"
	"github.com/go-openapi/spec"
)

// Generate generates a client
func Generate(modulePath string, s spec.Swagger) error {
	pkgName, ok := s.Info.Extensions.GetString("x-npm-package")
	if !ok {
		return errors.New("must provide 'x-npm-package' in the 'info' section of the swagger.yml")
	}

	tmplInfo := clientCodeTemplate{
		ClassName:   utils.CamelCase(s.Info.InfoProps.Title, true),
		PackageName: pkgName,
		ServiceName: s.Info.InfoProps.Title,
		Version:     s.Info.InfoProps.Version,
		Description: s.Info.InfoProps.Description,
	}

	for _, path := range swagger.SortedPathItemKeys(s.Paths.Paths) {
		pathItem := s.Paths.Paths[path]
		pathItemOps := swagger.PathItemOperations(pathItem)
		for _, method := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[method]
			if op.Deprecated {
				continue
			}
			methodCode, err := methodCode(s, op, method, path)
			if err != nil {
				return err
			}
			tmplInfo.Methods = append(tmplInfo.Methods, methodCode)
		}
	}

	indexDTS, err := generateTypescriptTypes(s)
	if err != nil {
		return err
	}

	errorsJS, err := generateErrorsFile(s)
	if err != nil {
		return err
	}

	indexJS, err := templates.WriteTemplate(indexJSTmplStr, tmplInfo)
	if err != nil {
		return err
	}

	packageJSON, err := templates.WriteTemplate(packageJSONTmplStr, tmplInfo)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(filepath.Join(modulePath, "types.js"), []byte(errorsJS), 0644); err != nil {
		return err
	}

	if err = ioutil.WriteFile(filepath.Join(modulePath, "index.js"), []byte(indexJS), 0644); err != nil {
		return err
	}

	if ioutil.WriteFile(filepath.Join(modulePath, "package.json"), []byte(packageJSON), 0644); err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(modulePath, "index.d.ts"), []byte(indexDTS), 0644)
}

type clientCodeTemplate struct {
	PackageName string
	ClassName   string
	ServiceName string
	Version     string
	Description string
	Methods     []string
}

var indexJSTmplStr = `const async = require("async");
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
   * @param {bool} [options.keepalive] - Set keepalive to true for client requests. This sets the
   * forever: true attribute in request. Defaults to true.
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
    if (options.tracer) {
      this.tracer = options.tracer;
    } else {
      this.tracer = opentracing.globalTracer();
    }

    const circuitOptions = Object.assign({}, defaultCircuitOptions, options.circuit);
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

const packageJSONTmplStr = `{
  "name": "{{.PackageName}}",
  "version": "{{.Version}}",
  "description": "{{.Description}}",
  "main": "index.js",
  "dependencies": {
    "async": "^2.1.4",
    "clever-discovery": "0.0.8",
    "opentracing": "^0.14.0",
    "request": "^2.87.0",
    "kayvee": "^3.13.0",
    "hystrixjs": "^0.2.0",
    "rxjs": "^5.4.1"
  },
  "devDependencies": {
    "typescript": "^3.3.0"
  }
}
`

const methodTmplStr = `
  {{.MethodDefinition}}
    {{if .IterMethod -}}
    const it = (f, saveResults) => new Promise((resolve, reject) => {
    {{- else -}}
    if (!cb && typeof options === "function") {
      options = undefined;
    }

    return new Promise((resolve, reject) => {
    {{- end}}
      if (!options) {
        options = {};
      }

      const timeout = options.timeout || this.timeout;
      const tracer = options.tracer || this.tracer;
      const span = options.span;

      const headers = {};
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

      if (span && typeof span.log === "function") {
        // Need to get tracer to inject. Use HTTP headers format so we can properly escape special characters
        tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers);
        {{- if not .IterMethod}}
        span.log({event: "{{.Method}} {{.Path}}"});
        {{- end}}
        span.setTag("span.kind", "client");
      }

      const requestOptions = {
        method: "{{.Method}}",
        uri: this.address + "{{.PathCode}}",
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
          if (span && typeof span.log === "function") {
            span.log({event: "{{.Method}} {{.Path}}"});
          }
      const address = this.address;
  {{- end}}
      let retries = 0;
      (function requestOnce() {
        request(requestOptions, async (err, response, body) => {
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
								for (let i = 0; i < body{{$.IterResourceAccessString}}.length; i++) {
									try {
										await f(body{{$.IterResourceAccessString}}[i], i, body);
									} catch(err) {
										reject(err);
									}
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
          requestOptions.useQuerystring = false;
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
      map: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, true]), cb),
      toArray: cb => applyCallback(this._hystrixCommand.execute(it, [x => x, true]), cb),
      forEach: (f, cb) => applyCallback(this._hystrixCommand.execute(it, [f, false]), cb),
    };
    {{- end}}
  }
`

const singleParamMethodDefinitionTemplateString = `/**{{if .Description}}
   * {{.Description}}{{end}}{{range $param := .Params}}
   * @param {{if $param.JSDocType}}{{.JSDocType}} {{end}}{{$param.JSName}}{{if $param.Default}}={{$param.Default}}{{end}}{{if $param.Description}} - {{.Description}}{{end}}{{end}}
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:{{.ServiceName}}.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   {{- if .IterMethod}}
   * @returns {Object} iter
   * @returns {function} iter.map - takes in a function, applies it to each resource, and returns a promise to the result as an array
   * @returns {function} iter.toArray - returns a promise to the resources as an array
   * @returns {function} iter.forEach - takes in a function, applies it to each resource
   {{- else}}
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill{{if .JSDocSuccessReturnType}} {{.JSDocSuccessReturnType}}{{else}} {*}{{end}}{{$ServiceName := .ServiceName}}{{range $response := .Responses}}{{if $response.IsError}}
   * @reject {module:{{$ServiceName}}.Errors.{{$response.Name}}}{{end}}{{end}}
   * @reject {Error}{{end}}
   */
  {{- if .IterMethod}}
  {{.MethodName}}({{range $param := .Params}}{{$param.JSName}}, {{end}}options) {
  {{- else}}
  {{.MethodName}}({{range $param := .Params}}{{$param.JSName}}, {{end}}options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._{{.MethodName}}, arguments), callback);
  }

  _{{.MethodName}}({{range $param := .Params}}{{$param.JSName}}, {{end}}options, cb) {
  {{- end}}
    const params = {};{{range $param := .Params}}
    params["{{$param.JSName}}"] = {{$param.JSName}};{{end}}
`

const pluralParamMethodDefinitionTemplateString = `/**{{if .Description}}
   * {{.Description}}{{end}}
   * @param {Object} params{{range $param := .Params}}
   * @param {{if $param.JSDocType}}{{.JSDocType}} {{end}}{{if not $param.Required}}[{{end}}params.{{$param.JSName}}{{if $param.Default}}={{$param.Default}}{{end}}{{if not $param.Required}}]{{end}}{{if $param.Description}} - {{.Description}}{{end}}{{end}}
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {external:Span} [options.span] - An OpenTracing span - For example from the parent request
   * @param {module:{{.ServiceName}}.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   {{- if .IterMethod}}
   * @returns {Object} iter
   * @returns {function} iter.map - takes in a function, applies it to each resource, and returns a promise to the result as an array
   * @returns {function} iter.toArray - returns a promise to the resources as an array
   * @returns {function} iter.forEach - takes in a function, applies it to each resource
   {{- else}}
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill{{if .JSDocSuccessReturnType}} {{.JSDocSuccessReturnType}}{{else}} {*}{{end}}{{$ServiceName := .ServiceName}}{{range $response := .Responses}}{{if $response.IsError}}
   * @reject {module:{{$ServiceName}}.Errors.{{$response.Name}}}{{end}}{{end}}
   * @reject {Error}{{end}}
   */
  {{- if .IterMethod}}
  {{.MethodName}}(params, options) {
  {{- else}}
  {{.MethodName}}(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._{{.MethodName}}, arguments), callback);
  }

  _{{.MethodName}}(params, options, cb) {
  {{- end}}`

type paramMapping struct {
	JSName      string
	WagName     string
	Required    bool
	JSDocType   string
	Default     interface{}
	Description string
}

type responseMapping struct {
	StatusCode int
	Name       string
	IsError    bool
	IsNoData   bool
}

type methodTemplate struct {
	ServiceName              string
	Operation                string
	MethodName               string
	IterMethod               bool
	IterResourceAccessString string
	Description              string
	MethodDefinition         string
	Params                   []paramMapping
	Method                   string
	PathCode                 string
	Path                     string
	HeaderParams             []paramMapping
	PathParams               []paramMapping
	QueryParams              []paramMapping
	BodyParam                string
	Responses                []responseMapping
	JSDocSuccessReturnType   string
}

// This function takes in a swagger path such as "/path/goes/to/{location}/and/to/{other_Location}"
// and returns a string of javacript code such as "/path/goes/to/" + location + "/and/to/" + otherLocation.
func fillOutPath(path string) string {
	paramRegex := regexp.MustCompile("({.+?})")
	paramNameRegex := regexp.MustCompile("{(.+?)}")
	return paramRegex.ReplaceAllStringFunc(path, func(param string) string {
		return paramNameRegex.ReplaceAllStringFunc(param, func(paramName string) string {
			return "\" + params." + utils.CamelCase(paramName, false) + " + \""
		})
	})
}

func methodCode(s spec.Swagger, op *spec.Operation, method, path string) (string, error) {
	basePath := s.BasePath
	tmplInfo := methodTemplate{
		ServiceName: s.Info.InfoProps.Title,
		Operation:   op.ID,
		MethodName:  op.ID, // might mutate to op.ID + "Iter" for paging methods
		Description: op.Description,
		Method:      method,
		PathCode:    basePath + fillOutPath(path),
		Path:        basePath + path,
	}

	var successResponse *spec.Response
	for _, statusCode := range swagger.SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
		if successResponse == nil && statusCode >= 200 && statusCode < 400 {
			r := op.Responses.StatusCodeResponses[statusCode]
			successResponse = &r
		}
		response := responseMapping{
			StatusCode: statusCode,
			IsError:    statusCode >= 400,
		}
		typeName, _ := swagger.OutputType(&s, op, statusCode)
		if strings.HasPrefix(typeName, "models.") {
			typeName = typeName[7:] // models.ResponseType -> ResponseType
		}
		response.Name = typeName
		if typeName == "" {
			response.IsNoData = true
		}
		tmplInfo.Responses = append(tmplInfo.Responses, response)
	}
	tmplInfo.JSDocSuccessReturnType = responseToJSDocReturnType(successResponse)

	for _, wagParam := range op.Parameters {
		param := paramMapping{
			JSName:      utils.CamelCase(wagParam.Name, false),
			WagName:     wagParam.Name,
			Required:    wagParam.Required,
			JSDocType:   paramToJSDocType(wagParam),
			Description: wagParam.Description,
			Default:     wagParam.Default,
		}

		tmplInfo.Params = append(tmplInfo.Params, param)
		switch wagParam.In {
		case "path":
			tmplInfo.PathParams = append(tmplInfo.PathParams, param)
		case "header":
			tmplInfo.HeaderParams = append(tmplInfo.HeaderParams, param)
		case "body": // Will only ever be a single bodyParam so we can just set here
			tmplInfo.BodyParam = param.JSName
		case "query":
			tmplInfo.QueryParams = append(tmplInfo.QueryParams, param)
		}
	}

	if err := fillMethodDefinition(op, &tmplInfo); err != nil {
		return "", err
	}

	res, err := templates.WriteTemplate(methodTmplStr, tmplInfo)
	if err != nil {
		return "", err
	}

	if _, hasPaging := swagger.PagingParam(op); hasPaging {
		tmplInfo.IterMethod = true
		tmplInfo.MethodName += "Iter"

		if err := fillMethodDefinition(op, &tmplInfo); err != nil {
			return "", err
		}

		resourcePath := swagger.PagingResourcePath(op)
		if len(resourcePath) > 0 {
			tmplInfo.IterResourceAccessString = "." + strings.Join(resourcePath, ".")
		}

		iterMethodCode, err := templates.WriteTemplate(methodTmplStr, tmplInfo)
		if err != nil {
			return "", err
		}
		res += "\n" + iterMethodCode
	}

	return res, nil
}

func fillMethodDefinition(op *spec.Operation, tmplInfo *methodTemplate) error {
	var err error
	var methodDefinition string
	if len(op.Parameters) <= 1 {
		methodDefinition, err = templates.WriteTemplate(singleParamMethodDefinitionTemplateString, tmplInfo)
	} else {
		methodDefinition, err = templates.WriteTemplate(pluralParamMethodDefinitionTemplateString, tmplInfo)
	}
	if err != nil {
		return err
	}
	tmplInfo.MethodDefinition = methodDefinition
	return err
}

// paramToJSDocType returns the JSDoc type to assign a parameter.
// It returns empty string if it cannot discern a type for the parameter.
func paramToJSDocType(param spec.Parameter) string {
	if param.Type == "string" || param.Type == "number" || param.Type == "boolean" {
		return "{" + param.Type + "}"
	} else if param.Type == "integer" {
		return "{number}"
	} else if param.Type == "array" && param.Items != nil &&
		(param.Items.Type == "string" || param.Items.Type == "number" || param.Items.Type == "boolean") {
		return fmt.Sprintf("{%s[]}", param.Items.Type)
	}
	log.Printf("TODO: unhandled param name=%s. Documentation will be incomplete for this parameter.", param.Name)
	return ""
}

// schemaToJSDocType returns the JSDoc type to assign a schema.
// It returns empty string if it cannot discern a type for the schema.
func schemaToJSDocType(schema *spec.Schema) string {
	if len(schema.Type) == 1 {
		typ := schema.Type[0]
		if typ == "string" {
			return "{string}"
		} else if typ == "integer" {
			return "{number}"
		}
	}
	log.Printf("TODO: unhandled schema type %v.", schema.Type)
	return ""
}

// responseToJSDocReturnType returns the JS Doc type for a response.
// It returns empty string if it cannot determine the type.
func responseToJSDocReturnType(r *spec.Response) string {
	if r.Schema == nil {
		return "{undefined}"
	} else if r.Schema.Type != nil && len(r.Schema.Type) == 1 && r.Schema.Type[0] == "array" &&
		r.Schema.Items != nil && r.Schema.Items.Schema != nil {
		return "{Object[]}" // in the future, a more specific type would be nice
	} else if r.Schema.Ref.String() != "" {
		return "{Object}" // in the future, a more specific type would be nice
	}
	log.Printf("TODO: unhandled response: %#v. Documentation will be incomplete for this response type.", *r)
	return ""
}

var typeTmplString = `module.exports.Errors = {};
{{$ServiceName := .ServiceName}}
{{range .ErrorTypes}}/**
 * {{.Name}}
 * @extends Error
 * @memberof module:{{$ServiceName}}
 * @alias module:{{$ServiceName}}.Errors.{{.Name}}{{range .JSDocProperties}}
 * @property {{.Type}} {{.Name}}{{end}}
 */
module.exports.Errors.{{ .Name }} = class extends Error {
  constructor(body) {
    super(body.message);
    for (const k of Object.keys(body)) {
      this[k] = body[k];
    }
  }
};

{{ end }}`

type typesTemplate struct {
	ServiceName string
	ErrorTypes  []errorType
}

type errorType struct {
	StatusCode      int
	Name            string
	JSDocProperties []jsDocProperty
}

type jsDocProperty struct {
	Type string
	Name string
}

func jsDocPropertyFromSchema(name string, schema *spec.Schema) jsDocProperty {
	return jsDocProperty{
		Name: name,
		Type: schemaToJSDocType(schema),
	}
}

func generateErrorsFile(s spec.Swagger) (string, error) {
	typesTmpl := typesTemplate{
		ServiceName: s.Info.InfoProps.Title,
	}

	typeNames := stringset.New()

	for _, pathKey := range swagger.SortedPathItemKeys(s.Paths.Paths) {
		path := s.Paths.Paths[pathKey]
		pathItemOps := swagger.PathItemOperations(path)
		for _, opKey := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[opKey]
			for _, statusCode := range swagger.SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
				if statusCode < 400 {
					continue
				}
				typeName, _ := swagger.OutputType(&s, op, statusCode)
				if strings.HasPrefix(typeName, "models.") {
					typeName = typeName[7:]
				}
				if typeNames.Contains(typeName) {
					continue
				}
				typeNames.Add(typeName)

				etype := errorType{
					StatusCode: statusCode,
					Name:       typeName,
				}

				if schema, ok := s.Definitions[typeName]; !ok {
					log.Printf("TODO: could not find schema for %s, JS documentation will be incomplete", typeName)
				} else if len(schema.Properties) > 0 {
					for _, name := range swagger.SortedSchemaProperties(schema) {
						propertySchema := schema.Properties[name]
						etype.JSDocProperties = append(etype.JSDocProperties, jsDocPropertyFromSchema(name, &propertySchema))
					}
				}

				typesTmpl.ErrorTypes = append(typesTmpl.ErrorTypes, etype)
			}
		}
	}

	return templates.WriteTemplate(typeTmplString, typesTmpl)
}

// JSType is a wrapper type for things that represent a TypeScript type
type JSType string

// JSTypeMap is a representation of a TypeScript object
type JSTypeMap map[string]JSType

type typescriptTypes struct {
	ServiceName   string
	IncludedTypes []string
	MethodDecls   []string
	ErrorTypes    []string
}

var isDefaultIncludedType = map[string]bool{
	"BadRequest":    true,
	"InternalError": true,
	"NotFound":      true,
}

var primitiveTypes = map[string]string{
	"string":  "string",
	"integer": "number",
	"number":  "number",
	"boolean": "boolean",
}

func generateTypescriptTypes(s spec.Swagger) (string, error) {
	tt := typescriptTypes{
		ServiceName:   utils.CamelCase(s.Info.InfoProps.Title, true),
		IncludedTypes: []string{},
		MethodDecls:   []string{},
	}
	includedTypeMap := JSTypeMap{}

	for _, path := range swagger.SortedPathItemKeys(s.Paths.Paths) {
		pathItem := s.Paths.Paths[path]
		pathItemOps := swagger.PathItemOperations(pathItem)
		for _, method := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[method]
			if op.Deprecated {
				continue
			}
			methodDecl, err := methodDecl(s, op, path, method)
			if err != nil {
				return "", err
			}
			tt.MethodDecls = append(tt.MethodDecls, methodDecl)
			err = addInputType(&includedTypeMap, op)
			if err != nil {
				return "", err
			}
		}
	}

	for name, schema := range s.Definitions {
		if !isDefaultIncludedType[name] {
			theType, err := asJSType(&schema, "")
			if err != nil {
				return "", err
			}
			includedTypeMap[name] = theType
		}
	}

	var keys []string
	for k := range includedTypeMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	includedTypes := []string{}
	for _, typeName := range keys {
		includedTypes = append(includedTypes, fmt.Sprintf("type %s = %s;", typeName, includedTypeMap[typeName]))
	}
	tt.IncludedTypes = includedTypes

	errorTypes, err := getErrorTypes(s)
	if err != nil {
		return "", err
	}
	tt.ErrorTypes = errorTypes

	types, err := templates.WriteTemplate(typescriptTmplStr, tt)
	if err != nil {
		return "", err
	}

	return types, nil
}

func getErrorTypes(s spec.Swagger) ([]string, error) {
	typeNames := map[string]struct{}{}
	errorTypes := []string{}

	for _, pathKey := range swagger.SortedPathItemKeys(s.Paths.Paths) {
		path := s.Paths.Paths[pathKey]
		pathItemOps := swagger.PathItemOperations(path)
		for _, opKey := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[opKey]
			for _, statusCode := range swagger.SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
				if statusCode < 400 {
					continue
				}
				typeName, _ := swagger.OutputType(&s, op, statusCode)
				if strings.HasPrefix(typeName, "models.") {
					typeName = typeName[7:]
				}

				if _, exists := typeNames[typeName]; exists {
					continue
				}
				typeNames[typeName] = struct{}{}

				if schema, ok := s.Definitions[typeName]; !ok {
					errorTypes = append(errorTypes, fmt.Sprintf("class %s {}", typeName))
				} else if len(schema.Properties) > 0 {
					declaration, err := generateErrorDeclaration(&schema, typeName, "models.")
					if err != nil {
						return errorTypes, err
					}
					errorTypes = append(errorTypes, declaration)
				}
			}
		}
	}
	return errorTypes, nil
}

func methodDecl(s spec.Swagger, op *spec.Operation, path, method string) (string, error) {
	returnType, err := ReturnType(s, op)
	if returnType != "void" && returnType != "never" {
		returnType = JSType(fmt.Sprintf("models.%s", returnType))
	}
	if err != nil {
		return "", err
	}
	methodName := op.ID
	var params string
	var methodDecl string
	if len(op.Parameters) == 0 {
		params = ""
	} else if len(op.Parameters) == 1 {
		paramName := op.Parameters[0].Name
		var paramType JSType
		if op.Parameters[0].ParamProps.Schema != nil {
			paramType, err = asJSType(op.Parameters[0].ParamProps.Schema, "")
			paramType = JSType(fmt.Sprintf("models.%s", paramType))
		} else {
			paramType, err = asJSTypeSimple(op.Parameters[0].SimpleSchema)
		}
		if err != nil {
			return "", err
		}
		if op.Parameters[0].Required {
			params = fmt.Sprintf("%s: %s, ", paramName, paramType)
		} else {
			params = fmt.Sprintf("%s?: %s, ", paramName, paramType)
		}
	} else {
		paramType := fmt.Sprintf("%sParams", utils.CamelCase(methodName, true))
		params = fmt.Sprintf("params: models.%s, ", paramType)
	}
	methodDecl = fmt.Sprintf("%s(%soptions?: RequestOptions, cb?: Callback<%s>): Promise<%s>",
		methodName, params, returnType, returnType)

	if _, hasPaging := swagger.PagingParam(op); hasPaging {
		resourcePath := swagger.PagingResourcePath(op)
		for _, ident := range resourcePath {
			returnType += JSType(fmt.Sprintf("[\"%s\"]", ident))
		}
		returnType = "ArrayInner<" + returnType + ">"
		methodDecl += fmt.Sprintf("\n  %s(%soptions?: RequestOptions): IterResult<%s>", methodName+"Iter", params, returnType)
	}
	return methodDecl, nil
}

type paramDeclTmpl struct {
	TypeName string
	Fields   string
}

func addInputType(jsTypeMap *JSTypeMap, op *spec.Operation) error {
	if len(op.Parameters) <= 1 {
		return nil
	}
	typeName := utils.CamelCase(op.ID+"Params", true)
	paramNames := []string{}
	fields := JSTypeMap{}
	for _, param := range op.Parameters {
		if param.In == "formData" {
			return fmt.Errorf("input parameters with 'In' formData are not supported")
		}
		paramType, err := paramToJSType(param)
		if err != nil {
			return err
		}
		paramType = ": " + paramType
		if !param.Required {
			paramType = "?" + paramType
		}
		paramNames = append(paramNames, param.Name)
		fields[param.Name] = paramType
	}
	fieldsStrings := []string{}
	for _, paramName := range paramNames {
		jsName := utils.CamelCase(paramName, false)
		fieldsStrings = append(fieldsStrings, fmt.Sprintf("%s%s;", jsName, fields[paramName]))
	}
	(*jsTypeMap)[typeName] = JSType(fmt.Sprintf("{\n  %s\n}", strings.Join(fieldsStrings, "\n  ")))
	return nil
}

const paramDeclTmlpStr = `interface {{.TypeName}} {
  {{.Fields}}
}`

const methodDeclTmplStr = `{{.Name}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param}}{{end}}): {{.ReturnType}}`

// ReturnType returns the methods return type
func ReturnType(s spec.Swagger, op *spec.Operation) (JSType, error) {
	successCodes := []int{}
	for statusCode := range op.Responses.StatusCodeResponses {
		if statusCode < 400 {
			successCodes = append(successCodes, statusCode)
		}
	}
	if len(successCodes) == 0 {
		return "never", nil
	} else if len(successCodes) == 1 {
		return typeOf(s, op, successCodes[0])
	}
	return "", fmt.Errorf("Operation %s has more than one possible success return type", op.ID)
}

func typeOf(s spec.Swagger, op *spec.Operation, statusCode int) (JSType, error) {
	schema := swagger.OutputSchema(&s, op, statusCode)
	if schema == nil {
		return "void", nil
	}
	return asJSType(schema, "")
}

func paramToJSType(param spec.Parameter) (JSType, error) {
	if param.In == "body" {
		typeName, err := asJSType(param.Schema, "")
		if err != nil {
			return "", err
		}
		return typeName, nil
	}

	if len(param.Enum) > 0 {
		return typeFromEnum(param.Enum)
	}

	var typeName string
	switch param.Type {
	case "string":
		typeName = "string"
	case "integer", "number":
		typeName = "number"
	case "boolean":
		typeName = "boolean"
	case "array":
		if param.Items.Type != "string" {
			return JSType(""), fmt.Errorf("array parameters must have string sub-types")
		}
		typeName = "string[]"
	default:
		// Note. We don't support 'array' or 'file' types even though they're in the
		// Swagger spec.
		return JSType(""), fmt.Errorf("unsupported param type: \"%s\"", param.Type)
	}
	return JSType(typeName), nil
}

func asJSTypeSimple(simpleSchema spec.SimpleSchema) (JSType, error) {
	if jsType, ok := primitiveTypes[simpleSchema.Type]; ok {
		return JSType(jsType), nil
	}

	if simpleSchema.Type == "array" {
		return JSType("any[]"), nil
	}

	if simpleSchema.Type == "object" {
		return JSType("any"), nil
	}

	return JSType(""), fmt.Errorf("Unknown type '%v'", simpleSchema.Type)
}

func asJSType(schema *spec.Schema, refPrefix string) (JSType, error) {
	if schema == nil {
		return JSType(""), fmt.Errorf("No schema")
	}

	if schema.Ref.String() != "" {
		def, err := defFromRef(schema.Ref.String())
		if err != nil {
			return JSType(""), err
		}
		if refPrefix != "" {
			return JSType(fmt.Sprintf("%s%s", refPrefix, def)), nil
		}
		return JSType(def), nil
	}

	if len(schema.Type) == 0 {
		if schema.AdditionalProperties != nil {
			if schema.AdditionalProperties.Schema != nil {
				innerType, err := asJSType(schema.AdditionalProperties.Schema, refPrefix)
				if err != nil {
					return JSType(""), err
				}
				return "{ [key: string]: " + innerType + " }", nil
			}
			if schema.AdditionalProperties.Allows {
				return "{ [key: string]: any }", nil
			}
		}
		return JSType("any"), nil
	}

	if len(schema.Enum) > 0 {
		return typeFromEnum(schema.Enum)
	}

	if len(schema.Type) > 1 {
		return JSType(""), fmt.Errorf("Having mltiple types in schema is not supported")
	}
	if jsType, ok := primitiveTypes[schema.Type[0]]; ok {
		return JSType(jsType), nil
	}

	if schema.Type[0] == "array" {
		innerType, err := asJSType(schema.Items.Schema, refPrefix)
		return innerType + "[]", err
	}

	if schema.Type[0] == "object" || len(schema.Properties) > 0 {
		fieldsStrings, err := generatePropertyDeclarations(schema, refPrefix)
		if err != nil {
			return JSType(""), err
		}

		if schema.AdditionalProperties != nil {
			if schema.AdditionalProperties.Schema != nil {
				innerType, err := asJSType(schema.AdditionalProperties.Schema, refPrefix)
				if err != nil {
					return JSType(""), err
				}
				fieldsStrings = append(fieldsStrings, "[key: string]: "+string(innerType)+";")
			} else if schema.AdditionalProperties.Allows {
				fieldsStrings = append(fieldsStrings, "[key: string]: any;")
			}
		}
		return JSType(fmt.Sprintf("{\n  %s\n}", strings.Join(fieldsStrings, "\n  "))), nil
	}

	return JSType(""), fmt.Errorf("Unknown type '%v'", schema.Type[0])
}

func typeFromEnum(enum []interface{}) (JSType, error) {
	enums := []string{}
	for _, enum := range enum {
		e, err := json.Marshal(enum)
		if err != nil {
			return JSType(""), err
		}
		enums = append(enums, string(e))
	}
	return JSType("(" + strings.Join(enums, " | ") + ")"), nil
}

func defFromRef(ref string) (string, error) {
	if strings.HasPrefix(ref, "#/definitions/") {
		return ref[len("#/definitions/"):], nil
	}
	return "", fmt.Errorf("schema.$ref has undefined reference type \"%s\". "+
		"Must start with #/definitions or #/responses.", ref)
}

const typescriptTmplStr = `import { Span, Tracer } from "opentracing";
import { Logger } from "kayvee";

type Callback<R> = (err: Error, result: R) => void;
type ArrayInner<R> = R extends (infer T)[] ? T : never;

interface RetryPolicy {
  backoffs(): number[];
  retry(requestOptions: {method: string}, err: Error, res: {statusCode: number}): boolean;
}

interface RequestOptions {
  timeout?: number;
  span?: Span;
  retryPolicy?: RetryPolicy;
}

interface IterResult<R> {
  map<T>(f: (r: R) => T, cb?: Callback<T[]>): Promise<T[]>;
  toArray(cb?: Callback<R[]>): Promise<R[]>;
  forEach(f: (r: R) => void, cb?: Callback<void>): Promise<void>;
}

interface CircuitOptions {
  forceClosed?: boolean;
  maxConcurrentRequests?: number;
  requestVolumeThreshold?: number;
  sleepWindow?: number;
  errorPercentThreshold?: number;
}

interface GenericOptions {
  timeout?: number;
  keepalive?: boolean;
  retryPolicy?: RetryPolicy;
  logger?: Logger;
  tracer?: Tracer;
  circuit?: CircuitOptions;
  serviceName?: string;
}

interface DiscoveryOptions {
  discovery: true;
  address?: undefined;
}

interface AddressOptions {
  discovery?: false;
  address: string;
}

type {{.ServiceName}}Options = (DiscoveryOptions | AddressOptions) & GenericOptions;

import models = {{.ServiceName}}.Models

declare class {{.ServiceName}} {
  constructor(options: {{.ServiceName}}Options);

  {{range .MethodDecls}}
  {{.}}
  {{end}}
}

declare namespace {{.ServiceName}} {
  const RetryPolicies: {
    Single: RetryPolicy;
    Exponential: RetryPolicy;
    None: RetryPolicy;
  }

  const DefaultCircuitOptions: CircuitOptions;

  namespace Errors {
    interface ErrorBody {
      message: string;
      [key: string]: any;
    }

    {{range .ErrorTypes}}
    {{.}}
    {{end}}
  }

  namespace Models {
    {{range .IncludedTypes}}
    {{.}}
    {{end}}
  }
}

export = {{.ServiceName}};
`
