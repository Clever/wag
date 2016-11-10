package goclient

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/Clever/wag/swagger"
	"github.com/Clever/wag/templates"
)

// Generate generates a client
func Generate(packageName string, s spec.Swagger) error {
	if err := generateClient(packageName, s); err != nil {
		return err
	}
	if err := generateInterface(packageName, &s, s.Info.InfoProps.Title, s.Paths); err != nil {
		return err
	}
	return nil
}

type clientCodeTemplate struct {
	PackageName          string
	ServiceName          string
	FormattedServiceName string
	Methods              []string
}

var clientCodeTemplateStr = `
package client

import (
		"context"
		"strings"
		"bytes"
		"net/http"
		"net/url"
		"strconv"
		"encoding/json"
		"strconv"
		"time"
		"fmt"
		"crypto/md5"

		"{{.PackageName}}/models"
		discovery "github.com/Clever/discovery-go"
		"github.com/afex/hystrix-go/hystrix"
)

var _ = json.Marshal
var _ = strings.Replace
var _ = strconv.FormatInt
var _ = bytes.Compare

// WagClient is used to make requests to the {{.ServiceName}} service.
type WagClient struct {
	basePath    string
	requestDoer doer
	transport   *http.Transport
	timeout     time.Duration
	// Keep the retry doer around so that we can set the number of retries
	retryDoer *retryDoer
	// Keep the circuit doer around so that we can turn it on / off
	circuitDoer    *circuitBreakerDoer
	defaultTimeout time.Duration
}

var _ Client = (*WagClient)(nil)

// New creates a new client. The base path and http transport are configurable.
func New(basePath string) *WagClient {
	base := baseDoer{}
	tracing := tracingDoer{d: base}
	// For the short-term don't use the default retry policy since its 5 retries can 5X
	// the traffic. Once we've enabled circuit breakers by default we can turn it on.
	retry := retryDoer{d: tracing, retryPolicy: SingleRetryPolicy{}}
	circuit := &circuitBreakerDoer{
		d:     &retry,
		debug: true,
		// one circuit for each service + url pair
		circuitName: fmt.Sprintf("{{.ServiceName}}-%%s", shortHash(basePath)),
	}
	circuit.init()
	client := &WagClient{requestDoer: circuit, retryDoer: &retry, circuitDoer: circuit, defaultTimeout: 10 * time.Second,
 		transport: &http.Transport{}, basePath: basePath}
	client.SetCircuitBreakerSettings(DefaultCircuitBreakerSettings)
	return client
}

// NewFromDiscovery creates a client from the discovery environment variables. This method requires
// the three env vars: SERVICE_{{.FormattedServiceName}}_HTTP_(HOST/PORT/PROTO) to be set. Otherwise it returns an error.
func NewFromDiscovery() (*WagClient, error) {
	url, err := discovery.URL("{{.ServiceName}}", "default")
	if err != nil {
		url, err = discovery.URL("{{.ServiceName}}", "http") // Added fallback to maintain reverse compatibility
		if err != nil {
			return nil, err
		}
	}
	return New(url), nil
}

// SetRetryPolicy sets a the given retry policy for all requests.
func (c *WagClient) SetRetryPolicy(retryPolicy RetryPolicy) {
	c.retryDoer.retryPolicy = retryPolicy
}

// SetCircuitBreakerDebug puts the circuit
func (c *WagClient) SetCircuitBreakerDebug(b bool) {
	c.circuitDoer.debug = b
}

// CircuitBreakerSettings are the parameters that govern the client's circuit breaker.
type CircuitBreakerSettings struct {
	// MaxConcurrentRequests is the maximum number of concurrent requests
	// the client can make at the same time. Default: 100.
	MaxConcurrentRequests int
	// RequestVolumeThreshold is the minimum number of requests needed
	// before a circuit can be tripped due to health. Default: 20.
	RequestVolumeThreshold int
	// SleepWindow how long, in milliseconds, to wait after a circuit opens
	// before testing for recovery. Default: 5000.
	SleepWindow int
	// ErrorPercentThreshold is the threshold to place on the rolling error
	// rate. Once the error rate exceeds this percentage, the circuit opens.
	// Default: 90.
	ErrorPercentThreshold int
}

// DefaultCircuitBreakerSettings describes the default circuit parameters.
var DefaultCircuitBreakerSettings = CircuitBreakerSettings{
	MaxConcurrentRequests:  100,
	RequestVolumeThreshold: 20,
	SleepWindow:            5000,
	ErrorPercentThreshold:  90,
}

// SetCircuitBreakerSettings sets parameters on the circuit breaker. It must be
// called on application startup.
func (c *WagClient) SetCircuitBreakerSettings(settings CircuitBreakerSettings) {
	hystrix.ConfigureCommand(c.circuitDoer.circuitName, hystrix.CommandConfig{
		// redundant, with the timeout we set on the context, so set
		// this to something high and irrelevant
		Timeout:                100 * 1000,
		MaxConcurrentRequests:  settings.MaxConcurrentRequests,
		RequestVolumeThreshold: settings.RequestVolumeThreshold,
		SleepWindow:            settings.SleepWindow,
		ErrorPercentThreshold:  settings.ErrorPercentThreshold,
	})
}

// SetTimeout sets a timeout on all operations for the client. To make a single request
// with a timeout use context.WithTimeout as described here: https://godoc.org/golang.org/x/net/context#WithTimeout.
func (c *WagClient) SetTimeout(timeout time.Duration){
	c.defaultTimeout = timeout
}

{{range $methodCode := .Methods}}
	{{$methodCode}}
{{end}}

func shortHash(s string) string {
	return fmt.Sprintf("%%x", md5.Sum([]byte(s)))[0:6]
}
`

func generateClient(packageName string, s spec.Swagger) error {

	codeTemplate := clientCodeTemplate{
		PackageName:          packageName,
		ServiceName:          s.Info.InfoProps.Title,
		FormattedServiceName: strings.ToUpper(strings.Replace(s.Info.InfoProps.Title, "-", "_", -1)),
	}

	for _, path := range swagger.SortedPathItemKeys(s.Paths.Paths) {
		pathItem := s.Paths.Paths[path]
		pathItemOps := swagger.PathItemOperations(pathItem)
		for _, method := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[method]
			if op.Deprecated {
				continue
			}
			codeTemplate.Methods = append(codeTemplate.Methods, methodCode(&s, op, s.BasePath, method, path))
		}
	}

	clientCode, err := templates.WriteTemplate(clientCodeTemplateStr, codeTemplate)
	if err != nil {
		return err
	}

	g := swagger.Generator{PackageName: packageName}
	g.Printf(clientCode)
	return g.WriteFile("client/client.go")
}

func generateInterface(packageName string, s *spec.Swagger, serviceName string, paths *spec.Paths) error {
	g := swagger.Generator{PackageName: packageName}
	g.Printf("package client\n\n")
	g.Printf(swagger.ImportStatements([]string{"context", packageName + "/models"}))
	g.Printf("//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_client.go -package=client\n\n")
	g.Printf("// Client defines the methods available to clients of the %s service.\n", serviceName)
	g.Printf("type Client interface {\n\n")

	for _, pathKey := range swagger.SortedPathItemKeys(paths.Paths) {
		path := paths.Paths[pathKey]
		pathItemOps := swagger.PathItemOperations(path)
		for _, method := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[method]
			if op.Deprecated {
				continue
			}

			g.Printf("\t%s\n", swagger.InterfaceComment(method, pathKey, pathItemOps[method]))
			g.Printf("\t%s\n\n", swagger.Interface(s, pathItemOps[method]))
		}
	}
	g.Printf("}\n")

	return g.WriteFile("client/interface.go")
}

func methodCode(s *spec.Swagger, op *spec.Operation, basePath, method, methodPath string) string {
	var buf bytes.Buffer
	capOpID := swagger.Capitalize(op.ID)

	buf.WriteString(swagger.InterfaceComment(method, methodPath, op) + "\n")
	buf.WriteString(fmt.Sprintf("func (c *WagClient) %s {\n", swagger.Interface(s, op)))
	buf.WriteString(fmt.Sprintf("\tpath := c.basePath + \"%s\"\n", basePath+methodPath))
	buf.WriteString(fmt.Sprintf("\turlVals := url.Values{}\n"))
	buf.WriteString(fmt.Sprintf("\tvar body []byte\n\n"))

	buf.WriteString(fmt.Sprintf(buildRequestCode(op, method)))

	buf.WriteString(fmt.Sprintf(`

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "%s")
	req = req.WithContext(ctx)
	// Don't add the timeout in a "doer" because we don't want to call "defer.cancel()"
	// until we've finished all the processing of the request object. Otherwise we'll cancel
	// our own request before we've finished it.
	if c.defaultTimeout != 0 {
		ctx, cancel := context.WithTimeout(req.Context(), c.defaultTimeout)
		defer cancel()
	    req = req.WithContext(ctx)
	}
	resp, err := c.requestDoer.Do(client, req)
	%s
	defer resp.Body.Close()
`, op.ID, errorMessage(op)))

	buf.WriteString(parseResponseCode(s, op, capOpID))

	return buf.String()

}

// buildRequestCode adds the parameters to the URL, the body, and the headers
func buildRequestCode(op *spec.Operation, method string) string {
	var buf bytes.Buffer

	singleStringPathParameter, singleParamName := swagger.SingleStringPathParameter(op)
	if singleStringPathParameter {
		buf.WriteString(fmt.Sprintf("\tpath = strings.Replace(path, \"%s\", %s, -1)\n",
			"{"+op.Parameters[0].Name+"}", singleParamName))
	} else {
		for _, param := range op.Parameters {
			if param.In == "path" {
				// TODO: Should this be done with regex at some point?
				buf.WriteString(fmt.Sprintf("\tpath = strings.Replace(path, \"%s\", %s, -1)\n",
					"{"+param.Name+"}", swagger.ParamToStringCode(param)))
			} else if param.In == "query" {
				var queryAddCode string
				if param.Type == "array" {
					queryAddCode = fmt.Sprintf("\tfor _, v := range i.%s {\n\t\turlVals.Add(\"%s\", v)\n\t}\n", swagger.StructParamName(param), param.Name)
				} else {
					queryAddCode = fmt.Sprintf("\turlVals.Add(\"%s\", %s)\n", param.Name, swagger.ParamToStringCode(param))
				}
				if param.Required {
					buf.WriteString(fmt.Sprintf(queryAddCode))
				} else {
					buf.WriteString(fmt.Sprintf("\tif i.%s != nil {\n", swagger.StructParamName(param)))
					buf.WriteString(fmt.Sprintf(queryAddCode))
					buf.WriteString(fmt.Sprintf("\t}\n"))
				}
			}
		}
	}

	buf.WriteString(fmt.Sprintf("\tpath = path + \"?\" + urlVals.Encode()\n\n"))

	for _, param := range op.Parameters {
		if param.In == "body" {
			singleSchemaedBodyParameter, _ := swagger.SingleSchemaedBodyParameter(op)
			var bodyMarshalCode string
			if singleSchemaedBodyParameter {
				// no wrapper struct for single-input methods
				bodyMarshalCode = fmt.Sprintf(`
	var err error
	body, err = json.Marshal(i)
	%s
`, errorMessage(op))
			} else {
				bodyMarshalCode = fmt.Sprintf(`
	var err error
	body, err = json.Marshal(i.%s)
	%s
`, swagger.StructParamName(param), errorMessage(op))
			}

			if param.Required {
				buf.WriteString(fmt.Sprintf(bodyMarshalCode))
			} else {
				if singleSchemaedBodyParameter {
					buf.WriteString("\tif i != nil {\n")
				} else {
					buf.WriteString(fmt.Sprintf("\tif i.%s != nil {\n", swagger.StructParamName(param)))
				}
				buf.WriteString(fmt.Sprintf(bodyMarshalCode))
				buf.WriteString(fmt.Sprintf("\t}\n"))
			}
		}
	}

	buf.WriteString(fmt.Sprintf(`
	client := &http.Client{Transport: c.transport}
	req, err := http.NewRequest("%s", path, bytes.NewBuffer(body))
	%s
`, strings.ToUpper(method), errorMessage(op)))

	for _, param := range op.Parameters {
		if param.In == "header" {
			headerAddCode := fmt.Sprintf("\treq.Header.Set(\"%s\", %s)\n", param.Name, swagger.ParamToStringCode(param))
			if param.Required {
				buf.WriteString(fmt.Sprintf(headerAddCode))
			} else {
				buf.WriteString(fmt.Sprintf("\tif i.%s != nil {\n", swagger.StructParamName(param)))
				buf.WriteString(fmt.Sprintf(headerAddCode))
				buf.WriteString(fmt.Sprintf("\t}\n"))
			}
		}
	}
	return buf.String()
}

func errorMessage(op *spec.Operation) string {
	str, err := templates.WriteTemplate(errMsgTemplStr, errMsgTmpl{NoSuccessType: swagger.NoSuccessType(op)})
	if err != nil {
		panic("internal error generating client")
	}
	return str
}

type errMsgTmpl struct {
	NoSuccessType bool
}

var errMsgTemplStr = `
	{{if .NoSuccessType}}
		if err != nil {
			return err
		}
	{{else}}	
		if err != nil {
			return nil, err
		}
	{{end}}
`

type statusCodeReturn struct {
	responseTypes []string
	// unclear if we need this decode param
	decode      bool
	makePointer bool
}

// parseResponseCode generates the code for handling the http response.
// In the client code we want to return a different object depending on the status code, so
// let's generate code that switches on the status code and returns the right object in each
// case.
func parseResponseCode(s *spec.Swagger, op *spec.Operation, capOpID string) string {
	var buf bytes.Buffer

	buf.WriteString("\tswitch resp.StatusCode {\n")

	for _, statusCode := range swagger.SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
		statusCodeDecoder, err := writeStatusCodeDecoder(s, op, statusCode)
		if err != nil {
			// TODO: move this up???
			panic(fmt.Errorf("error parsing response code: %s", err))
		}
		buf.WriteString(statusCodeDecoder)
	}

	// It would be nice if we could remove this too
	noSuccessType := swagger.NoSuccessType(op)
	successReturn := "nil, "
	if noSuccessType {
		successReturn = ""
	}

	// TODO: at some point should encapsulate this behind an interface on the operation
	errorType, _ := swagger.OutputType(s, op, 500)
	buf.WriteString(fmt.Sprintf(`
	default:
		return %s&%s{Message: "Unknown response"}
	}
}

`, successReturn, errorType))

	return buf.String()
}

func writeStatusCodeDecoder(s *spec.Swagger, op *spec.Operation, statusCode int) (string, error) {
	outputName, makePointer := swagger.OutputType(s, op, statusCode)

	// TODO: Need makePointer to handle arrays... not sure if there's a better way to do this...
	outputType := "output"
	if makePointer {
		outputType = "&output"
	}

	return templates.WriteTemplate(codeDetectorTmplStr,
		codeDetectorTmpl{
			StatusCode:    statusCode,
			NoSuccessType: swagger.NoSuccessType(op),
			ErrorType:     statusCode >= 400,
			TypeName:      outputName,
			OutputType:    outputType,
		})
}

type codeDetectorTmpl struct {
	StatusCode    int
	NoSuccessType bool
	ErrorType     bool
	TypeName      string
	OutputType    string
}

var codeDetectorTmplStr = `
	case {{.StatusCode}}:

	{{if .NoSuccessType}}
		{{if .ErrorType}}
		var output {{.TypeName}}
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return err
		}
		return {{.OutputType}}
		{{else}}
		return nil
		{{end}}
	{{else}}
		{{if .ErrorType}}
		var output {{.TypeName}}
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return nil, {{.OutputType}}
		{{else}}
		var output {{.TypeName}}
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return {{.OutputType}}, nil
		{{end}}
	{{end}}
`
