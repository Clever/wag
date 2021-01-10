package goclient

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/Clever/wag/v6/swagger"
	"github.com/Clever/wag/v6/templates"
	"github.com/Clever/wag/v6/utils"
)

// Generate generates a client
func Generate(packageName, packagePath string, s spec.Swagger) error {
	if err := generateClient(packageName, packagePath, s); err != nil {
		return err
	}
	return generateInterface(packageName, packagePath, &s, s.Info.InfoProps.Title, s.Paths)
}

type clientCodeTemplate struct {
	PackageName          string
	ServiceName          string
	FormattedServiceName string
	Operations           []string
	Version              string
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
		"{{.PackageName}}/tracing"
		discovery "github.com/Clever/discovery-go"
		"github.com/afex/hystrix-go/hystrix"
		logger "gopkg.in/Clever/kayvee-go.v6/logger"
)

var _ = json.Marshal
var _ = strings.Replace
var _ = strconv.FormatInt
var _ = bytes.Compare

// Version of the client.
const Version = "{{ .Version }}"

// VersionHeader is sent with every request.
const VersionHeader = "X-Client-Version"

// WagClient is used to make requests to the {{.ServiceName}} service.
type WagClient struct {
	basePath    string
	requestDoer doer
	client   	*http.Client
	// Keep the retry doer around so that we can set the number of retries
	retryDoer *retryDoer
	// Keep the circuit doer around so that we can turn it on / off
	circuitDoer    *circuitBreakerDoer
	defaultTimeout time.Duration
	logger       logger.KayveeLogger
}

var _ Client = (*WagClient)(nil)

// New creates a new client. The base path and http transport are configurable.
func New(basePath string) *WagClient {
	basePath = strings.TrimSuffix(basePath, "/")
	base := baseDoer{}
	// For the short-term don't use the default retry policy since its 5 retries can 5X
	// the traffic. Once we've enabled circuit breakers by default we can turn it on.
	retry := retryDoer{d: base, retryPolicy: SingleRetryPolicy{}}
	logger := logger.New("{{.ServiceName}}-wagclient")
	circuit := &circuitBreakerDoer{
		d:     &retry,
		// TODO: INFRANG-4404 allow passing circuitBreakerOptions
		debug: true,
		// one circuit for each service + url pair
		circuitName: fmt.Sprintf("{{.ServiceName}}-%%s", shortHash(basePath)),
		logger: logger,
	}
	circuit.init()
	client := &WagClient{
		basePath: basePath,
		requestDoer: circuit,
		client: &http.Client{
			Transport: tracing.NewTransport(http.DefaultTransport, opNameCtx{}),
		},
		retryDoer: &retry,
		circuitDoer: circuit,
		defaultTimeout: 5 * time.Second,
		logger: logger,
	}
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
func (c *WagClient) SetLogger(logger logger.KayveeLogger) {
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

// SetTimeout sets a timeout on all operations for the client. To make a single request with a shorter timeout
// than the default on the client, use context.WithTimeout as described here: https://godoc.org/golang.org/x/net/context#WithTimeout.
func (c *WagClient) SetTimeout(timeout time.Duration){
	c.defaultTimeout = timeout
}

// SetTransport sets the http transport used by the client.
func (c *WagClient) SetTransport(t http.RoundTripper){
	c.client.Transport = tracing.NewTransport(t, opNameCtx{})
}

{{range $operationCode := .Operations}}
	{{$operationCode}}
{{end}}

func shortHash(s string) string {
	return fmt.Sprintf("%%x", md5.Sum([]byte(s)))[0:6]
}
`

func generateClient(packageName, packagePath string, s spec.Swagger) error {

	codeTemplate := clientCodeTemplate{
		PackageName:          packageName,
		ServiceName:          s.Info.InfoProps.Title,
		FormattedServiceName: strings.ToUpper(strings.Replace(s.Info.InfoProps.Title, "-", "_", -1)),
		Version:              s.Info.InfoProps.Version,
	}

	for _, path := range swagger.SortedPathItemKeys(s.Paths.Paths) {
		pathItem := s.Paths.Paths[path]
		pathItemOps := swagger.PathItemOperations(pathItem)
		for _, method := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[method]
			if op.Deprecated {
				continue
			}
			code, err := operationCode(&s, op, s.BasePath, method, path, IsBinaryBody(op, s.Definitions))
			if err != nil {
				return err
			}
			codeTemplate.Operations = append(codeTemplate.Operations, code)
		}
	}

	clientCode, err := templates.WriteTemplate(clientCodeTemplateStr, codeTemplate)
	if err != nil {
		return err
	}

	g := swagger.Generator{PackagePath: packagePath}
	g.Printf(clientCode)
	return g.WriteFile("client/client.go")
}

// IsBinaryBody returns true if the format of the body of the operation is binary
func IsBinaryBody(op *spec.Operation, definitions map[string]spec.Schema) bool {
	for _, param := range op.Parameters {
		if param.In == "body" {
			return IsBinaryParam(param, definitions)
		}
	}
	return false
}

// IsBinaryParam returns true of the format of the parameter is binary
func IsBinaryParam(param spec.Parameter, definitions map[string]spec.Schema) bool {
	definitionName := path.Base(param.Schema.Ref.Ref.GetURL().String())
	return definitions[definitionName].Format == "binary"
}

func generateInterface(packageName, packagePath string, s *spec.Swagger, serviceName string, paths *spec.Paths) error {
	g := swagger.Generator{PackagePath: packagePath}
	g.Printf("package client\n\n")
	g.Printf(swagger.ImportStatements([]string{"context", packageName + "/models"}))
	g.Printf("//go:generate mockgen -source=$GOFILE -destination=mock_client.go -package=client\n\n")

	if err := generateClientInterface(s, &g, serviceName, paths); err != nil {
		return err
	}
	if err := generateIteratorTypes(s, &g, paths); err != nil {
		return err
	}

	return g.WriteFile("client/interface.go")
}

func generateClientInterface(s *spec.Swagger, g *swagger.Generator, serviceName string, paths *spec.Paths) error {
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
			_, hasPaging := swagger.PagingParam(pathItemOps[method])
			if hasPaging {
				g.Printf("\t%s\n\n", swagger.ClientIterInterface(s, pathItemOps[method]))
			}
		}
	}
	g.Printf("}\n\n")
	return nil
}

func generateIteratorTypes(s *spec.Swagger, g *swagger.Generator, paths *spec.Paths) error {
	for _, pathKey := range swagger.SortedPathItemKeys(paths.Paths) {
		path := paths.Paths[pathKey]
		pathItemOps := swagger.PathItemOperations(path)
		for _, method := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[method]
			if op.Deprecated {
				continue
			}
			_, hasPaging := swagger.PagingParam(pathItemOps[method])
			if hasPaging {
				capOpID := swagger.Capitalize(op.ID)
				resourceType, _, err := swagger.PagingResourceType(s, op)
				if err != nil {
					return err
				}

				g.Printf("// %sIter defines the methods available on %s iterators.\n", capOpID, capOpID)
				g.Printf("type %sIter interface {\n", capOpID)
				g.Printf("\tNext(*%s) bool\n", resourceType)
				g.Printf("\tErr() error\n")
				g.Printf("}\n\n")
			}
		}
	}
	return nil
}

func operationCode(s *spec.Swagger, op *spec.Operation, basePath, method, methodPath string, binaryBody bool) (string, error) {
	var buf bytes.Buffer

	generatedMethodCodeString, err := methodCode(s, op, basePath, method, methodPath, binaryBody)
	if err != nil {
		return "", err
	}

	buf.WriteString(generatedMethodCodeString)
	if _, hasPaging := swagger.PagingParam(op); hasPaging {
		iter, err := iterCode(s, op, basePath, methodPath, method)
		if err != nil {
			return "", err
		}
		buf.WriteString(iter)
	}
	buf.WriteString(fmt.Sprintf(methodDoerCode(s, op)))
	return buf.String(), nil
}

func methodCode(s *spec.Swagger, op *spec.Operation, basePath, method, methodPath string, binaryBody bool) (string, error) {
	var buf bytes.Buffer
	capOpID := swagger.Capitalize(op.ID)

	interfaceComment, err := swagger.InterfaceComment(method, methodPath, true, s, op)
	if err != nil {
		return "", err
	}
	buf.WriteString(interfaceComment + "\n")
	buf.WriteString(fmt.Sprintf("func (c *WagClient) %s {\n", swagger.ClientInterface(s, op)))

	buf.WriteString(fmt.Sprintf("\theaders := make(map[string]string)\n\n"))
	if !binaryBody {
		buf.WriteString(fmt.Sprintf("\tvar body []byte\n"))
	}

	buf.WriteString(fmt.Sprintf(buildPathCode(s, op, basePath, methodPath)))
	buf.WriteString(fmt.Sprintf(buildHeadersCode(s, op)))
	buf.WriteString(fmt.Sprintf(buildRequestCode(s, op, method, binaryBody)))

	if _, hasPaging := swagger.PagingParam(op); !hasPaging {
		buf.WriteString(fmt.Sprintf(`
	return c.do%sRequest(ctx, req, headers)
}

`, capOpID))
	} else {
		buf.WriteString(fmt.Sprintf(`
	resp, _, err := c.do%sRequest(ctx, req, headers)
	return resp, err
}

`, capOpID))
	}

	return buf.String(), nil
}

func methodDoerCode(s *spec.Swagger, op *spec.Operation) string {
	var buf bytes.Buffer
	capOpID := swagger.Capitalize(op.ID)

	errReturn := ""
	returnType := ""
	if successType := swagger.SuccessType(s, op); successType != nil {
		errReturn += "nil, "
		returnType += *successType + ", "
	}
	if _, hasPaging := swagger.PagingParam(op); hasPaging {
		errReturn += "\"\", "
		returnType += "string, "
	}
	returnType += "error"

	if len(returnType) != len("error") {
		returnType = "(" + returnType + ")"
	}

	buf.WriteString(fmt.Sprintf(`
func (c *WagClient) do%sRequest(ctx context.Context, req *http.Request, headers map[string]string) %s {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "%s")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

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
	resp, err := c.requestDoer.Do(c.client, req)
	retCode := 0
	if resp != nil {
	  retCode = resp.StatusCode
	}

	// log all client failures and non-successful HT
	logData := logger.M{
		"backend": "%s",
		"method": req.Method,
		"uri": req.URL,
		"status_code": retCode,
	}
	if err == nil && retCode > 399 {
		logData["message"] = resp.Status
		c.logger.ErrorD("client-request-finished", logData)
	}
	if err != nil {
		logData["message"] = err.Error()
		c.logger.ErrorD("client-request-finished", logData)
		return %serr
	}
	defer resp.Body.Close()
`, capOpID, returnType, op.ID, op.ID, s.Info.InfoProps.Title, errReturn))

	buf.WriteString(parseResponseCode(s, op, capOpID))

	return buf.String()
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

func buildBodyCode(s *spec.Swagger, op *spec.Operation, method string) string {
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
			return str
		}
	}
	return ""
}

func getAccessString(op *spec.Operation) string {
	for _, param := range op.Parameters {
		if param.In == "body" {
			t := swagger.ParamToTemplate(&param, op)
			if singleParam, _ := swagger.SingleSchemaedBodyParameter(op); singleParam {
				t.AccessString = "i"
			}
			return t.AccessString
		}
	}
	panic("unexpected error: no body in request")
}

// buildRequestCode adds the body and makes the request
func buildRequestCode(s *spec.Swagger, op *spec.Operation, method string, binaryBody bool) string {
	var buf bytes.Buffer

	// binary bodies are io.ReadCloser and do not need to be transformed
	if binaryBody {
		buf.WriteString(fmt.Sprintf(`
	req, err := http.NewRequestWithContext(ctx, "%s", path, *%s)
	%s
`, strings.ToUpper(method), getAccessString(op), errorMessage(s, op)))
	} else {
		buf.WriteString(buildBodyCode(s, op, method))
		buf.WriteString(fmt.Sprintf(`
	req, err := http.NewRequestWithContext(ctx, "%s", path, bytes.NewBuffer(body))
	%s
`, strings.ToUpper(method), errorMessage(s, op)))
	}

	return buf.String()
}

// buildHeadersCode adds the parameters to the header
func buildHeadersCode(s *spec.Swagger, op *spec.Operation) string {
	var buf bytes.Buffer

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
	headers["{{.Name}}"] = {{.ToStringCode}}
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

// buildSuccessReturn builds the zero values of the success portion of an op's
// return (so that an error can be appended).
func buildSuccessReturn(s *spec.Swagger, op *spec.Operation) string {
	ret := ""
	if successType := swagger.SuccessType(s, op); successType != nil {
		ret = ret + "nil, "
	}
	if _, hasPaging := swagger.PagingParam(op); hasPaging {
		ret = ret + "\"\", "
	}
	return ret
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

	successReturn := buildSuccessReturn(s, op)

	// TODO: at some point should encapsulate this behind an interface on the operation
	errorType, _ := swagger.OutputType(s, op, 500)
	buf.WriteString(fmt.Sprintf(`
	default:
		return %s&%s{Message: fmt.Sprintf("Unknown status code %%%%%%%%v", resp.StatusCode)}
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

	_, hasPaging := swagger.PagingParam(op)
	return templates.WriteTemplate(codeDetectorTmplStr,
		codeDetectorTmpl{
			StatusCode:    statusCode,
			NoSuccessType: swagger.SuccessType(s, op) == nil,
			SuccessReturn: buildSuccessReturn(s, op),
			HasPaging:     hasPaging,
			ErrorType:     statusCode >= 400,
			TypeName:      outputName,
			OutputType:    outputType,
		})
}

type codeDetectorTmpl struct {
	StatusCode    int
	NoSuccessType bool
	SuccessReturn string
	HasPaging     bool
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
			return {{.SuccessReturn}}err
		}
		return {{.SuccessReturn}}{{.OutputType}}
		{{else}}
		return {{.SuccessReturn}}nil
		{{end}}
	{{else}}
		{{if .ErrorType}}
		var output {{.TypeName}}
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return {{.SuccessReturn}}err
		}
		return {{.SuccessReturn}}{{.OutputType}}
		{{else}}
		var output {{.TypeName}}
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return {{.SuccessReturn}}err
		}
		{{if .HasPaging}}
		return {{.OutputType}}, resp.Header.Get("X-Next-Page-Path"), nil
		{{else}}
		return {{.OutputType}}, nil
		{{end}}
		{{end}}
	{{end}}
`

func iterCode(s *spec.Swagger, op *spec.Operation, basePath, methodPath, method string) (string, error) {
	capOpID := swagger.Capitalize(op.ID)
	resourceType, needsPointer, err := swagger.PagingResourceType(s, op)
	if err != nil {
		return "", err
	}

	var responseType string
	if needsPointer {
		responseType = fmt.Sprintf("[]*%s", resourceType)
	} else {
		responseType = fmt.Sprintf("[]%s", resourceType)
	}

	resourceAccessString := ""
	for _, pathComponent := range swagger.PagingResourcePath(op) {
		resourceAccessString = resourceAccessString + "." + utils.CamelCase(pathComponent, true)
	}

	return templates.WriteTemplate(
		iterTmplStr,
		iterTmpl{
			OpID:                 op.ID,
			CapOpID:              capOpID,
			Input:                swagger.OperationInput(op),
			BuildPathCode:        buildPathCode(s, op, basePath, methodPath),
			BuildHeadersCode:     buildHeadersCode(s, op),
			BuildBodyCode:        buildBodyCode(s, op, method),
			Method:               method,
			ResponseType:         responseType,
			ResourceType:         resourceType,
			ResponseAccessString: resourceAccessString,
			PointerArray:         needsPointer,
		},
	)
}

type iterTmpl struct {
	OpID                 string
	CapOpID              string
	Input                string
	BuildPathCode        string
	BuildHeadersCode     string
	BuildBodyCode        string
	Method               string
	ResponseType         string
	ResourceType         string
	ResponseAccessString string
	PointerArray         bool
}

var iterTmplStr = `
type {{.OpID}}IterImpl struct {
	c            *WagClient
	ctx          context.Context
	lastResponse {{.ResponseType}}
	index        int
	err          error
	nextURL      string
	headers      map[string]string
	body         []byte
}

// New{{.OpID}}Iter constructs an iterator that makes calls to {{.OpID}} for
// each page.
func (c *WagClient) New{{.CapOpID}}Iter(ctx context.Context, {{.Input}}) ({{.CapOpID}}Iter, error) {
	{{.BuildPathCode}}

	headers := make(map[string]string)
	{{.BuildHeadersCode}}

	var body []byte
	{{.BuildBodyCode}}

	return &{{.OpID}}IterImpl{
		c:            c,
		ctx:          ctx,
		lastResponse: {{.ResponseType}}{},
		nextURL:      path,
		headers:      headers,
		body:         body,
	}, nil
}

func (i *{{.OpID}}IterImpl) refresh() error {
	req, err := http.NewRequestWithContext(i.ctx, "{{.Method}}", i.nextURL, bytes.NewBuffer(i.body))

	if err != nil {
		i.err = err
		return err
	}

	resp, nextPage, err := i.c.do{{.CapOpID}}Request(i.ctx, req, i.headers)
	if err != nil {
		i.err = err
		return err
	}

	i.lastResponse = resp{{.ResponseAccessString}}
	i.index = 0
	if nextPage != "" {
		i.nextURL = i.c.basePath + nextPage
	} else {
		i.nextURL = ""
	}
	return nil
}

// Next retrieves the next resource from the iterator and assigns it to the
// provided pointer, fetching a new page if necessary. Returns true if it
// successfully retrieves a new resource.
func (i *{{.OpID}}IterImpl) Next(v *{{.ResourceType}}) bool {
	if i.err != nil {
		return false
	} else if i.index < len(i.lastResponse) {
		*v = {{if .PointerArray}}*{{end}}i.lastResponse[i.index]
		i.index++
		return true
	} else if i.nextURL == "" {
		return false
	}

	if err := i.refresh(); err != nil {
		return false
	}
	return i.Next(v)
}

// Err returns an error if one occurred when .Next was called.
func (i *{{.OpID}}IterImpl) Err() error {
	return i.err
}
`
