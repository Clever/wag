package goclient

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/Clever/wag/v8/swagger"
	"github.com/Clever/wag/v8/templates"
	"github.com/Clever/wag/v8/utils"
)

// Generate generates a client
func Generate(packageName, basePath string, s spec.Swagger) error {
	if err := generateClient(packageName, basePath, s); err != nil {
		return err
	}
	return generateInterface(packageName, basePath, &s, s.Info.InfoProps.Title, s.Paths)
}

type clientCodeTemplate struct {
	PackageName          string
	ServiceName          string
	FormattedServiceName string
	Operations           []string
	Version              string
	VersionSuffix        string
}

var clientCodeTemplateStr = `
package client

// Using Alpha version of WAG Yay!
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
		"os"
		"crypto/md5"

		"github.com/Clever/{{.ServiceName}}/gen-go/models{{.VersionSuffix}}"

		discovery "github.com/Clever/discovery-go"
		wcl "github.com/Clever/wag/logging/wagclientlogger"


		"github.com/afex/hystrix-go/hystrix"

		"go.opentelemetry.io/otel"
		"go.opentelemetry.io/otel/propagation"
		"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
		"go.opentelemetry.io/otel/sdk/resource"
		sdktrace "go.opentelemetry.io/otel/sdk/trace"
		semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
		
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
	logger      wcl.WagClientLogger
}

var _ Client = (*WagClient)(nil)


//This pattern is used instead of using closures for greater transparency and the ability to implement additional interfaces.
type options struct {
	transport    http.RoundTripper
	logger       wcl.WagClientLogger
	instrumentor Instrumentor
	exporter     sdktrace.SpanExporter
}

type Option interface {
	apply(*options)
}


//WithLogger sets client logger option.
func WithLogger(log wcl.WagClientLogger) Option {
	return loggerOption{Log: log}
}

type loggerOption struct {
	Log wcl.WagClientLogger
}

func (l loggerOption) apply(opts *options) {
	opts.logger = l.Log
}


type roundTripperOption struct {
	rt http.RoundTripper
}

func (t roundTripperOption) apply(opts *options) {
	opts.transport = t.rt
}

// WithRoundTripper allows you to pass in intrumented/custom roundtrippers which will then wrap the
// transport roundtripper
func WithRoundTripper(t http.RoundTripper) Option {
	return roundTripperOption{rt: t}
}

// Instrumentor is a function that creates an instrumented round tripper
type Instrumentor func(baseTransport http.RoundTripper, spanNameCtxValue interface{}, tp sdktrace.TracerProvider) http.RoundTripper

// WithInstrumentor sets a instrumenting function that will be used to wrap the roundTripper for tracing.
// For standard instrumentation with tracing use tracing.InstrumentedTransport, default is non-instrumented.

func WithInstrumentor(fn Instrumentor) Option {
	return instrumentorOption{instrumentor: fn}
}

type instrumentorOption struct {
	instrumentor Instrumentor
}

func (i instrumentorOption) apply(opts *options) {
	opts.instrumentor = i.instrumentor
}

// WithExporter sets client span exporter option.
func WithExporter(se sdktrace.SpanExporter) Option {
	return exporterOption{exporter: se}
}

type exporterOption struct {
	exporter sdktrace.SpanExporter
}

func (se exporterOption) apply(opts *options) {
	opts.exporter = se.exporter
}

//----------------------BEGIN LOGGING RELATED FUNCTIONS----------------------


//NewLogger creates a logger for id that produces logs at and below the indicated level.
//Level indicated the level at and below which logs are created.
func NewLogger(id string, level wcl.LogLevel) PrintlnLogger {
	return PrintlnLogger{id: id, level: level}
}

type PrintlnLogger struct {
	level wcl.LogLevel
	id    string
}

func (w PrintlnLogger) Log(level wcl.LogLevel, message string, m map[string]interface{}) {

	if level >= level {
		m["id"] = w.id
		jsonLog, err := json.Marshal(m)
		if err != nil {
			jsonLog, err = json.Marshal(map[string]interface{}{"Error Marshalling Log": err})
		}
		fmt.Println(string(jsonLog))
	}
}

//----------------------END LOGGING RELATED FUNCTIONS------------------------

//----------------------BEGIN TRACING RELATED FUNCTIONS----------------------


// newResource returns a resource describing this application.
// Used for setting up tracer provider
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("{{.ServiceName}}"),
			semconv.ServiceVersionKey.String("{{ .Version }}"),
		),
	)
	return r
}

func newTracerProvider(exporter sdktrace.SpanExporter, samplingProbability float64) *sdktrace.TracerProvider {

	tp:= sdktrace.NewTracerProvider(
		// We use the default ID generator. In order for sampling to work (at least with this sampler)
		// the ID generator must generate trace IDs uniformly at random from the entire space of uint64.
		// For example, the default x-ray ID generator does not do this.
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		// These maximums are to guard against something going wrong and sending a ton of data unexpectedly
		sdktrace.WithSpanLimits(sdktrace.SpanLimits{
			AttributeCountLimit: 100,
			EventCountLimit:     100,
			LinkCountLimit:      100,
		}),
		//Batcher is more efficient, switch to it after testing
		sdktrace.WithSyncer(exporter),
		//sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource()),
		
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func doNothing(baseTransport http.RoundTripper, spanNameCtxValue interface{}, tp sdktrace.TracerProvider) http.RoundTripper {
	return baseTransport
}

func determineSampling() (samplingProbability float64, err error) {

		// If we're running locally, then turn off sampling. Otherwise sample
		// 1%% or whatever TRACING_SAMPLING_PROBABILITY specifies.
		samplingProbability = 0.01
		isLocal := os.Getenv("_IS_LOCAL") == "true"
		if isLocal {
			fmt.Println("Set to Local")
			samplingProbability = 1.0
		} else if v := os.Getenv("TRACING_SAMPLING_PROBABILITY"); v != "" {
			samplingProbabilityFromEnv, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return 0, fmt.Errorf("could not parse '%%s' to float", v)
			}
			samplingProbability = samplingProbabilityFromEnv
		}
		return
	}

//----------------------END TRACING RELATEDFUNCTIONS----------------------

// New creates a new client. The base path and http transport are configurable.
func New(ctx context.Context, basePath string, opts ...Option) *WagClient {

	defaultTransport := http.DefaultTransport
	defaultLogger := NewLogger("{{.ServiceName}}-wagclient", wcl.Info)
	defaultExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		fmt.Println(err)
	}
	defaultInstrumentor := doNothing


	basePath = strings.TrimSuffix(basePath, "/")
	base := baseDoer{}
	// For the short-term don't use the default retry policy since its 5 retries can 5X
	// the traffic. Once we've enabled circuit breakers by default we can turn it on.
	retry := retryDoer{d: base, retryPolicy: SingleRetryPolicy{}}
	options := options{
		transport:    defaultTransport,
		logger:       defaultLogger,
		exporter:     defaultExporter,
		instrumentor: defaultInstrumentor,
	}

	for _, o := range opts {
		o.apply(&options)
	}


	samplingProbability := 1.0 // TODO: Put back logic to set this to 1 for local, 0.1 otherwise etc.
	// samplingProbability := determineSampling()

	tp := newTracerProvider(options.exporter, samplingProbability)
	options.transport = options.instrumentor(options.transport, ctx, *tp)

	circuit := &circuitBreakerDoer{
		d:     &retry,
		// TODO: INFRANG-4404 allow passing circuitBreakerOptions
		debug: true,
		// one circuit for each service + url pair
		circuitName: fmt.Sprintf("{{.ServiceName}}-%%s", shortHash(basePath)),
		logger: options.logger,
	}
	circuit.init()
	client := &WagClient{
		basePath: basePath,
		requestDoer: circuit,
		client: &http.Client{
			Transport: options.transport,
		},
		retryDoer: &retry,
		circuitDoer: circuit,
		defaultTimeout: 5 * time.Second,
		 logger: options.logger,
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
	return New(context.Background(), url), nil
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
func (c *WagClient) SetLogger(l wcl.WagClientLogger) {
	c.logger = l
	c.circuitDoer.logger = l
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
	// c.client.Transport = tracing.NewTransport(t, opNameCtx{})
}

{{range $operationCode := .Operations}}
	{{$operationCode}}
{{end}}

func shortHash(s string) string {
	return fmt.Sprintf("%%x", md5.Sum([]byte(s)))[0:6]
}
`

func getVersionSuffix(version string) string {
	num, err := strconv.Atoi(strings.TrimSuffix(string(version[0:2]), "."))
	if err != nil {
		return ""
	}
	if num <= 1 {
		return ""
	}

	return "/v" + strconv.Itoa(num)

}

func generateClient(packageName, basePath string, s spec.Swagger) error {

	codeTemplate := clientCodeTemplate{
		PackageName:          packageName,
		ServiceName:          s.Info.InfoProps.Title,
		FormattedServiceName: strings.ToUpper(strings.Replace(s.Info.InfoProps.Title, "-", "_", -1)),
		Version:              s.Info.InfoProps.Version,
		VersionSuffix:        getVersionSuffix(s.Info.InfoProps.Version),
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

	g := swagger.Generator{BasePath: basePath}
	g.Printf(clientCode)
	err = g.WriteFile("client/client.go")
	if err != nil {
		return err
	}

	return CreateModFile("client/go.mod", basePath, codeTemplate)

}

//CreateModFile creates a go.mod file for the client module.
func CreateModFile(path string, basePath string, codeTemplate clientCodeTemplate) error {

	absPath := basePath + "/" + path
	f, err := os.Create(absPath)

	if err != nil {
		return err
	}

	defer f.Close()
	modFileString := `
module github.com/Clever/` + codeTemplate.ServiceName + `/gen-go/client` + codeTemplate.VersionSuffix + `

go 1.16

require (
	//removed this because it can never get the right version unless I tag it first. Adding with: go get github.com/Clever/dapple/gen-go/models@INFRANG-5015
	//github.com/Clever/` + codeTemplate.ServiceName + `/gen-go/models` + codeTemplate.VersionSuffix + ` v` + codeTemplate.Version + `
	github.com/Clever/discovery-go v1.8.1
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/donovanhide/eventsource v0.0.0-20171031113327-3ed64d21fb0b
)

require (
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/go-openapi/analysis v0.21.2 // indirect
	github.com/go-openapi/errors v0.20.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/loads v0.21.1 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/strfmt v0.21.2 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/go-openapi/validate v0.22.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/smartystreets/goconvey v1.7.2 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.1-0.20200424115421-065759f9c3d7 // indirect
	go.mongodb.org/mongo-driver v1.7.5 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect

)
//Replace directives will work locally but mess up imports.
//replace github.com/Clever/` + codeTemplate.ServiceName + `/gen-go/models` + codeTemplate.VersionSuffix + ` v` + codeTemplate.Version + ` => ../models `

	_, err = f.WriteString(modFileString)

	if err != nil {
		return err
	}

	return nil
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
func extractModuleNameAndVersionSuffix(packageName string) (moduleName string, versionSuffix string) {
	regex, err := regexp.Compile("/v[0-9]/gen-go$|/v[0-9][0-9]/gen-go")
	if err != nil {
		log.Fatalf("Error getting module name from packageName: %s", err.Error())
	}
	versionSuffix = strings.TrimSuffix(regex.FindString(packageName), "/gen-go")
	if bool(regex.MatchString(packageName)) {
		moduleName = regex.ReplaceAllString(packageName, "")
	} else {
		moduleName = strings.TrimSuffix(packageName, "/gen-go")
	}
	return

}
func generateInterface(packageName, basePath string, s *spec.Swagger, serviceName string, paths *spec.Paths) error {
	g := swagger.Generator{BasePath: basePath}
	g.Printf("package client\n\n")
	moduleName, versionSuffix := extractModuleNameAndVersionSuffix(packageName)
	g.Printf(swagger.ImportStatements([]string{"context", moduleName + "/gen-go/models" + versionSuffix}))
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
	req.Done=true

	resp, err := c.requestDoer.Do(c.client, req)
	retCode := 0
	if resp != nil {
	  retCode = resp.StatusCode
	}

	// log all client failures and non-successful HT
	logData := map[string]interface{}{
		"backend": "%s",
		"method": req.Method,
		"uri": req.URL,
		"status_code": retCode,
	}
	if err == nil && retCode > 399 {
		logData["message"] = resp.Status 
		c.logger.Log(wcl.Error, "client-request-finished", logData)
	}
	if err != nil {
		logData["message"] = err.Error()
		c.logger.Log(wcl.Error, "client-request-finished", logData)
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
