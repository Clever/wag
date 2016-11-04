package jsclient

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"

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
			methodCode, err := methodCode(op, method, s.BasePath, path)
			if err != nil {
				return err
			}
			tmplInfo.Methods = append(tmplInfo.Methods, methodCode)
		}
	}

	indexJS, err := templates.WriteTemplate(indexJSTmplStr, tmplInfo)
	if err != nil {
		return err
	}

	packageJSON, err := templates.WriteTemplate(packageJSONTmplStr, tmplInfo)
	if err != nil {
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

var indexJSTmplStr = `const discovery = require("@clever/discovery");
const request = require("request");
const url = require("url");
const opentracing = require("opentracing");

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

module.exports = class {{.ClassName}} {

  constructor(options) {
    options = options || {};

    if (options.discovery) {
      try {
        this.address = discovery("{{.ServiceName}}", "http").url();
      } catch (e) {
        this.address = discovery("{{.ServiceName}}", "default").url();
      };
    } else if (options.address) {
      this.address = options.address;
    } else {
      throw new Error("Cannot initialize {{.ServiceName}} without discovery or address");
    }
    if (options.timeout) {
      this.timeout = options.timeout
    }
    if (options.retryPolicy) {
      this.retryPolicy = options.retryPolicy;
    }
  }
{{range $methodCode := .Methods}}{{$methodCode}}{{end}}}

module.exports.RetryPolicies = {
  Default: defaultRetryPolicy,
  None: noRetryPolicy,
};
`

var packageJSONTmplStr = `{
  "name": "{{.PackageName}}",
  "version": "{{.Version}}",
  "description": "{{.Description}}",
  "main": "index.js",
  "dependencies": {
    "@clever/discovery": "0.0.8",
    "opentracing": "^0.11.1",
    "request": "^2.75.0"
  }
}
`

var methodTmplStr = `
  {{.MethodDefinition}}
    if (!cb && typeof options === "function") {
      cb = options;
      options = undefined;
    }

    if (!options) {
      options = {};
    }

    const timeout = options.timeout || this.timeout;
    const span = options.span;

    const headers = {};{{range $param := .HeaderParams}}
    headers["{{$param.WagName}}"] = params.{{$param.JSName}};{{end}}

    const query = {};{{range $param := .QueryParams}}{{ if $param.Required }}
    query["{{$param.WagName}}"] = params.{{$param.JSName}};
{{else}}
    if (typeof params.{{$param.JSName}} !== "undefined") {
      query["{{$param.WagName}}"] = params.{{$param.JSName}};
    }
{{end}}{{end}}

    if (span) {
      opentracing.inject(span, opentracing.FORMAT_TEXT_MAP, headers);
      span.logEvent("{{.Method}} {{.Path}}");
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
{{ if ne .BodyParam ""}}
    requestOptions.body = params.{{.BodyParam}};
{{ end }}
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
`

var singleParamMethodDefinionTemplateString = `{{.MethodName}}({{range $param := .Params}}{{$param.JSName}}, {{end}}options, cb) {
    const params = {};{{range $param := .Params}}
    params["{{$param.JSName}}"] = {{$param.JSName}};{{end}}
`

var pluralParamMethodDefinionTemplateString = `{{.MethodName}}(params, options, cb) {`

type paramMapping struct {
	JSName   string
	WagName  string
	Required bool
}

type methodTemplate struct {
	MethodName       string
	MethodDefinition string
	Params           []paramMapping
	Method           string
	PathCode         string
	Path             string
	HeaderParams     []paramMapping
	QueryParams      []paramMapping
	BodyParam        string
}

// This function takes in a swagger path such as "/path/goes/to/{location}/and/to/{other_Location}"
// and returns a string of javacript code such as "/path/goes/to/" + location + "/and/to/" + otherLocation
func fillOutPath(path string) string {
	paramRegex := regexp.MustCompile("({.+?})")
	paramNameRegex := regexp.MustCompile("{(.+?)}")
	return paramRegex.ReplaceAllStringFunc(path, func(param string) string {
		return paramNameRegex.ReplaceAllStringFunc(param, func(paramName string) string {
			return "\" + params." + utils.CamelCase(paramName, false) + " + \""
		})
	})
}

func methodCode(op *spec.Operation, method, basePath, path string) (string, error) {

	tmplInfo := methodTemplate{
		MethodName: op.ID,
		Method:     method,
		PathCode:   basePath + fillOutPath(path),
		Path:       basePath + path,
	}

	for _, wagParam := range op.Parameters {
		param := paramMapping{
			JSName:   utils.CamelCase(wagParam.Name, false),
			WagName:  wagParam.Name,
			Required: wagParam.Required,
		}

		tmplInfo.Params = append(tmplInfo.Params, param)
		switch wagParam.In {
		case "header":
			tmplInfo.HeaderParams = append(tmplInfo.HeaderParams, param)
		case "body": // Will only ever be a single bodyParam so we can just set here
			tmplInfo.BodyParam = param.JSName
		case "query":
			tmplInfo.QueryParams = append(tmplInfo.QueryParams, param)
		}
	}

	var err error
	var methodDefinition string
	if len(op.Parameters) <= 1 {
		methodDefinition, err = templates.WriteTemplate(singleParamMethodDefinionTemplateString, tmplInfo)
	} else {
		methodDefinition, err = templates.WriteTemplate(pluralParamMethodDefinionTemplateString, tmplInfo)
	}
	if err != nil {
		return "", err
	}
	tmplInfo.MethodDefinition = methodDefinition
	return templates.WriteTemplate(methodTmplStr, tmplInfo)
}
