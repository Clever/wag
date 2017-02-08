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
		"strconv"
		"encoding/json"
		"strconv"
		"time"
		"fmt"
		"crypto/md5"

		"{{.PackageName}}/models"
		discovery "github.com/Clever/discovery-go"
		"github.com/afex/hystrix-go/hystrix"
		logger "gopkg.in/Clever/kayvee-go.v5/logger"
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
	logger       *logger.Logger
}

var _ Client = (*WagClient)(nil)

// New creates a new client. The base path and http transport are configurable.
func New(basePath string) *WagClient {
	base := baseDoer{}
	tracing := tracingDoer{d: base}
	// For the short-term don't use the default retry policy since its 5 retries can 5X
	// the traffic. Once we've enabled circuit breakers by default we can turn it on.
	retry := retryDoer{d: tracing, retryPolicy: SingleRetryPolicy{}}
	logger := logger.New("{{.ServiceName}}-wagclient")
	circuit := &circuitBreakerDoer{
		d:     &retry,
		debug: true,
		// one circuit for each service + url pair
		circuitName: fmt.Sprintf("{{.ServiceName}}-%%s", shortHash(basePath)),
		logger: logger,
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

// SetLogger allows for setting a custom logger
func (c *WagClient) SetLogger(logger *logger.Logger) {
	c.logger = logger
	c.circuitDoer.logger = logger
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
			method, err := methodCode(&s, op, s.BasePath, method, path)
			if err != nil {
				return err
			}
			codeTemplate.Methods = append(codeTemplate.Methods, method)
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

			interfaceComment, err := swagger.InterfaceComment(method, pathKey, true, s, pathItemOps[method])
			if err != nil {
				return err
			}
			g.Printf("\t%s\n", interfaceComment)
			g.Printf("\t%s\n\n", swagger.ClientInterface(s, pathItemOps[method]))
		}
	}
	g.Printf("}\n")

	return g.WriteFile("client/interface.go")
}

func methodCode(s *spec.Swagger, op *spec.Operation, basePath, method, methodPath string) (string, error) {
	var buf bytes.Buffer
	capOpID := swagger.Capitalize(op.ID)

	interfaceComment, err := swagger.InterfaceComment(method, methodPath, true, s, op)
	if err != nil {
		return "", err
	}
	buf.WriteString(interfaceComment + "\n")
	buf.WriteString(fmt.Sprintf("func (c *WagClient) %s {\n", swagger.ClientInterface(s, op)))

	buf.WriteString(fmt.Sprintf("\tvar body []byte\n\n"))
	buf.WriteString(fmt.Sprintf(buildPathCode(s, op, basePath, methodPath)))
	buf.WriteString(fmt.Sprintf(buildRequestCode(s, op, method)))

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
`, op.ID, errorMessage(s, op)))

	buf.WriteString(parseResponseCode(s, op, capOpID))

	return buf.String(), nil
}

func buildPathCode(s *spec.Swagger, op *spec.Operation, basePath, methodPath string) string {
	var buf bytes.Buffer
	capOpID := swagger.Capitalize(op.ID)

	if singleParam, _ := swagger.SingleSchemaedBodyParameter(op); len(op.Parameters) == 0 || singleParam {
		buf.WriteString(fmt.Sprintf("\tpath := c.basePath + \"%s\"\n", basePath+methodPath))
	} else if singleParam, singleParamName := swagger.SingleStringPathParameter(op); singleParam {
		buf.WriteString(fmt.Sprintf(
			"\tpath, err := models.%sInputPath(%s)\n",
			capOpID,
			singleParamName,
		))
		buf.WriteString(errorMessage(s, op))
		buf.WriteString(fmt.Sprintf("\tpath = c.basePath + path\n"))
	} else {
		buf.WriteString(fmt.Sprintf("\tpath, err := i.Path()\n"))
		buf.WriteString(errorMessage(s, op))
		buf.WriteString(fmt.Sprintf("\tpath = c.basePath + path\n"))
	}

	return buf.String()
}

// buildRequestCode adds the parameters to the URL, the body, and the headers
func buildRequestCode(s *spec.Swagger, op *spec.Operation, method string) string {
	var buf bytes.Buffer

	for _, param := range op.Parameters {
		if param.In == "body" {
			t := swagger.ParamToTemplate(&param, op)
			if singleParam, _ := swagger.SingleSchemaedBodyParameter(op); singleParam {
				t.AccessString = "i"
			}
			bodyTemplate := bodyParamTemplate{
				ParamTemplate: t,
				ErrorMessage:  errorMessage(s, op),
			}
			str, err := templates.WriteTemplate(bodyParamStr, bodyTemplate)
			if err != nil {
				panic(fmt.Errorf("unexpected error: %s", err))
			}
			buf.WriteString(str)
		}
	}

	buf.WriteString(fmt.Sprintf(`
	client := &http.Client{Transport: c.transport}
	req, err := http.NewRequest("%s", path, bytes.NewBuffer(body))
	%s
`, strings.ToUpper(method), errorMessage(s, op)))

	for _, param := range op.Parameters {
		if param.In == "header" {
			t := swagger.ParamToTemplate(&param, op)
			str, err := templates.WriteTemplate(headerParamStr, t)
			if err != nil {
				panic(fmt.Errorf("unexpected error: %s", err))
			}
			buf.WriteString(str)
		}
	}
	return buf.String()
}

var headerParamStr = `
	{{if .Pointer}}
	if {{.AccessString}} != nil {
	{{end}}
	req.Header.Set("{{.Name}}", {{.ToStringCode}})
	{{if .Pointer}}
	}
	{{end}}
`

type bodyParamTemplate struct {
	swagger.ParamTemplate
	ErrorMessage string
}

var bodyParamStr = `
	{{if .Pointer}}
	if {{.AccessString}} != nil {
	{{end}}
	var err error
	body, err = json.Marshal({{.AccessString}})
	{{.ErrorMessage}}	
	{{if .Pointer}}
	}
	{{end}}
`

func errorMessage(s *spec.Swagger, op *spec.Operation) string {
	str, err := templates.WriteTemplate(errMsgTemplStr, errMsgTmpl{
		NoSuccessType: swagger.SuccessType(s, op) == nil})
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

	successType := swagger.SuccessType(s, op)
	successReturn := "nil, "
	if successType == nil {
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
			NoSuccessType: swagger.SuccessType(s, op) == nil,
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
