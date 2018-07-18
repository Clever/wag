package jsclient

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
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
			methodCode, err := methodCode(s, op, method, s.BasePath, path)
			if err != nil {
				return err
			}
			tmplInfo.Methods = append(tmplInfo.Methods, methodCode)
		}
	}

	typeFileCode, err := generateTypesFile(s)
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

	if err = ioutil.WriteFile(filepath.Join(modulePath, "types.js"), []byte(typeFileCode), 0644); err != nil {
		return err
	}

	if err = ioutil.WriteFile(filepath.Join(modulePath, "index.js"), []byte(indexJS), 0644); err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(modulePath, "package.json"), []byte(packageJSON), 0644)
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
   * forever: true attribute in request. Defaults to false
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
        this.address = discovery("{{.ServiceName}}", "http").url();
      } catch (e) {
        this.address = discovery("{{.ServiceName}}", "default").url();
      }
    } else if (options.address) {
      this.address = options.address;
    } else {
      throw new Error("Cannot initialize {{.ServiceName}} without discovery or address");
    }
    if (options.keepalive) {
      this.keepalive = options.keepalive
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
      this.logger =  new kayvee.logger("{{.ServiceName}}-wagclient");
    }

    const circuitOptions = Object.assign({}, defaultCircuitOptions, options.circuit);
    this._hystrixCommand = commandFactory.getOrCreate("{{.ServiceName}}").
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
`

var packageJSONTmplStr = `{
  "name": "{{.PackageName}}",
  "version": "{{.Version}}",
  "description": "{{.Description}}",
  "main": "index.js",
  "dependencies": {
    "async": "^2.1.4",
    "clever-discovery": "0.0.8",
    "opentracing": "^0.11.1",
    "request": "^2.87.0",
    "kayvee": "^3.8.2",
    "hystrixjs": "^0.2.0",
    "rxjs": "^5.4.1"
  }
}
`

var methodTmplStr = `
  {{.MethodDefinition}}
    {{if .IterMethod -}}
    const it = (f, saveResults, cb) => new Promise((resolve, reject) => {
    {{- else -}}
    if (!cb && typeof options === "function") {
      cb = options;
      options = undefined;
    }

    return new Promise((resolve, reject) => {
    {{- end}}
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
      {{- range $param := .PathParams}}
      if (!params.{{$param.JSName}}) {
        rejecter(new Error("{{$param.JSName}} must be non-empty because it's a path parameter"));
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

      if (span) {
        opentracing.inject(span, opentracing.FORMAT_TEXT_MAP, headers);
        {{- if not .IterMethod}}
        span.logEvent("{{.Method}} {{.Path}}");
        {{- end}}
        span.setTag("span.kind", "client");
      }

	  const requestOptions = {
        method: "{{.Method}}",
        uri: this.address + "{{.PathCode}}",
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
          if (span) {
            span.logEvent("{{.Method}} {{.Path}}");
          }
      const address = this.address;
  {{- end}}
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
            {{- if not .IterMethod}}
            rejecter(err);
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
              rejecter(err);
              {{- else}}
              cbW(err);
              {{- end}}
              return;
            {{else}}{{if $response.IsNoData}}
              resolver();
              break;
            {{else}}
              {{if $.IterMethod -}}
              if (saveResults) {
                results = results.concat(body{{$.IterResourceAccessString}}.map(f));
              } else {
                body{{$.IterResourceAccessString}}.forEach(f);
              }
              {{- else -}}
              resolver(body);
              {{- end}}
              break;
            {{end}}{{end}}
            {{end}}default:
              var err = new Error("Received unexpected statusCode " + response.statusCode);
              responseLog(logger, requestOptions, response, err);
              {{- if not .IterMethod}}
              rejecter(err);
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
            rejecter(err);
            return;
          }
          if (saveResults) {
            resolver(results);
          } else {
            resolver();
          }
        }
      );
      {{- end}}
    });

    {{- if .IterMethod}}

    return {
      map: (f, cb) => this._hystrixCommand.execute(it, [f, true, cb]),
      toArray: cb => this._hystrixCommand.execute(it, [x => x, true, cb]),
      forEach: (f, cb) => this._hystrixCommand.execute(it, [f, false, cb]),
    };
    {{- end}}
  }
`

var singleParamMethodDefinitionTemplateString = `/**{{if .Description}}
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
    return this._hystrixCommand.execute(this._{{.MethodName}}, arguments);
  }
  _{{.MethodName}}({{range $param := .Params}}{{$param.JSName}}, {{end}}options, cb) {
  {{- end}}
    const params = {};{{range $param := .Params}}
    params["{{$param.JSName}}"] = {{$param.JSName}};{{end}}
`

var pluralParamMethodDefinitionTemplateString = `/**{{if .Description}}
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
    return this._hystrixCommand.execute(this._{{.MethodName}}, arguments);
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

func methodCode(s spec.Swagger, op *spec.Operation, method, basePath, path string) (string, error) {

	tmplInfo := methodTemplate{
		ServiceName: s.Info.InfoProps.Title,
		MethodName:  op.ID,
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

func generateTypesFile(s spec.Swagger) (string, error) {
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
