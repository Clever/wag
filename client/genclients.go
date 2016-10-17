package client

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/Clever/wag/swagger"
)

// Generate generates a client
func Generate(packageName string, s spec.Swagger) error {
	if err := generateClient(packageName, s); err != nil {
		return err
	}
	if err := generateInterface(packageName, s.Info.InfoProps.Title, s.Paths); err != nil {
		return err
	}
	return nil
}

func generateClient(packageName string, s spec.Swagger) error {
	g := swagger.Generator{PackageName: packageName}

	g.Printf(`package client

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

		"%[1]s/models"

		discovery "github.com/Clever/discovery-go"
)

var _ = json.Marshal
var _ = strings.Replace
var _ = strconv.FormatInt
var _ = bytes.Compare

// WagClient is used to make requests to the %[2]s service.
type WagClient struct {
	basePath    string
	requestDoer doer
	transport   *http.Transport
	timeout     time.Duration
	// Keep the retry doer around so that we can set the number of retries
	retryDoer *retryDoer
	defaultTimeout time.Duration
}

var _ Client = (*WagClient)(nil)

// New creates a new client. The base path and http transport are configurable.
func New(basePath string) *WagClient {
	base := baseDoer{}
	tracing := tracingDoer{d: base}
	retry := retryDoer{d: tracing, defaultRetries: 1}

	return &WagClient{requestDoer: &retry, retryDoer: &retry, defaultTimeout: 10 * time.Second,
		transport: &http.Transport{}, basePath: basePath}
}

// NewFromDiscovery creates a client from the discovery environment variables. This method requires
// the three env vars: SERVICE_%[3]s_HTTP_(HOST/PORT/PROTO) to be set. Otherwise it returns an error.
func NewFromDiscovery() (*WagClient, error) {
	url, err := discovery.URL("%[2]s", "default")
	if err != nil {
		url, err = discovery.URL("%[2]s", "http") // Added fallback to maintain reverse compatibility
		if err != nil {
			return nil, err
		}
	}
	return New(url), nil
}

// WithRetries returns a new client that retries all GET operations until they either succeed or fail the
// number of times specified.
func (c *WagClient) WithRetries(retries int) *WagClient {
	c.retryDoer.defaultRetries = retries
	return c
}

// WithTimeout returns a new client that has the specified timeout on all operations. To make a single request
// have a timeout use context.WithTimeout as described here: https://godoc.org/golang.org/x/net/context#WithTimeout.
func (c *WagClient) WithTimeout(timeout time.Duration) *WagClient {
	c.defaultTimeout = timeout
	return c
}

`, packageName,
		s.Info.InfoProps.Title,
		strings.ToUpper(strings.Replace(s.Info.InfoProps.Title, "-", "_", -1)))

	g.Printf(swagger.BaseParamToStringCode())

	for _, path := range swagger.SortedPathItemKeys(s.Paths.Paths) {
		pathItem := s.Paths.Paths[path]
		pathItemOps := swagger.PathItemOperations(pathItem)
		for _, method := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[method]
			g.Printf(methodCode(op, s.BasePath, method, path))
		}
	}

	return g.WriteFile("client/client.go")
}

func generateInterface(packageName string, serviceName string, paths *spec.Paths) error {
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
			g.Printf("\t%s\n", swagger.InterfaceComment(method, pathKey, pathItemOps[method]))
			g.Printf("\t%s\n\n", swagger.Interface(pathItemOps[method]))
		}
	}
	g.Printf("}\n")

	return g.WriteFile("client/interface.go")
}

func methodCode(op *spec.Operation, basePath, method, methodPath string) string {
	var buf bytes.Buffer
	capOpID := swagger.Capitalize(op.ID)

	buf.WriteString(swagger.InterfaceComment(method, methodPath, op) + "\n")
	buf.WriteString(fmt.Sprintf("func (c *WagClient) %s {\n", swagger.Interface(op)))
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
`, op.ID, errorMessage("models.DefaultInternalError{Msg: err.Error()}", op)))

	buf.WriteString(parseResponseCode(op, capOpID))

	return buf.String()

}

// buildRequestCode adds the parameters to the URL, the body, and the headers
func buildRequestCode(op *spec.Operation, method string) string {
	var buf bytes.Buffer

	for _, param := range op.Parameters {
		if param.In == "path" {
			// TODO: Should this be done with regex at some point?
			buf.WriteString(fmt.Sprintf("\tpath = strings.Replace(path, \"%s\", %s, -1)\n",
				"{"+param.Name+"}", swagger.ParamToStringCode(param)))

		} else if param.In == "query" {
			queryAddCode := fmt.Sprintf("\turlVals.Add(\"%s\", %s)\n", param.Name, swagger.ParamToStringCode(param))
			if param.Required {
				buf.WriteString(fmt.Sprintf(queryAddCode))
			} else {
				buf.WriteString(fmt.Sprintf("\tif i.%s != nil {\n", swagger.StructParamName(param)))
				buf.WriteString(fmt.Sprintf(queryAddCode))
				buf.WriteString(fmt.Sprintf("\t}\n"))
			}
		}
	}

	buf.WriteString(fmt.Sprintf("\tpath = path + \"?\" + urlVals.Encode()\n\n"))

	for _, param := range op.Parameters {
		if param.In == "body" {
			singleInputSchema, _ := swagger.SingleSchemaedBodyParameter(op)
			var bodyMarshalCode string
			if singleInputSchema {
				// no wrapper struct for single-input methods
				bodyMarshalCode = fmt.Sprintf(`
	var err error
	body, err = json.Marshal(i)
	%s
`, errorMessage("err", op))
			} else {
				bodyMarshalCode = fmt.Sprintf(`
	var err error
	body, err = json.Marshal(i.%s)
	%s
`, swagger.StructParamName(param), errorMessage("err", op))
			}

			if param.Required {
				buf.WriteString(fmt.Sprintf(bodyMarshalCode))
			} else {
				if singleInputSchema {
					buf.WriteString("\tif i != nil {\n")
				} else {
					buf.WriteString(fmt.Sprintf("\tif i.%s != nil {\n", swagger.StructParamName(param)))
				}
				buf.WriteString(fmt.Sprintf(bodyMarshalCode))
				buf.WriteString(fmt.Sprintf("\t}\n"))
			}
		}
	}

	// TODO: decide how to represent this error (should it be model.DefaultInternalError?)
	// and whether the client should retry
	buf.WriteString(fmt.Sprintf(`
	client := &http.Client{Transport: c.transport}
	req, err := http.NewRequest("%s", path, bytes.NewBuffer(body))
	%s
`, strings.ToUpper(method), errorMessage("err", op)))

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

// parseResponseCode generates the code for handling the http response.
// In the client code we want to return a different object depending on the status code, so
// let's generate code that switches on the status code and returns the right object in each
// case. This includes the default 400 and 500 cases.
func parseResponseCode(op *spec.Operation, capOpID string) string {
	var buf bytes.Buffer

	buf.WriteString("\tswitch resp.StatusCode {\n")

	noSuccessType := swagger.NoSuccessType(op)
	for _, statusCode := range swagger.SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
		response := op.Responses.StatusCodeResponses[statusCode]
		outputName := swagger.OutputType(op, statusCode)

		buf.WriteString(fmt.Sprintf("\tcase %d:\n", statusCode))

		if noSuccessType && statusCode < 400 {
			buf.WriteString("\t\treturn nil\n")
		} else if response.Schema == nil {
			buf.WriteString(fmt.Sprintf("\t\tvar output %s\n", outputName))
			if statusCode < 400 {
				buf.WriteString("\t\treturn output, nil\n")
			} else {
				if noSuccessType {
					buf.WriteString("\t\treturn output\n")
				} else {
					buf.WriteString("\t\treturn nil, output\n")
				}
			}
		} else {
			if statusCode < 400 {
				pointer := "&"
				// No pointer for array types (TODO: consider factoring this out...)
				if response.Schema.Ref.String() == "" {
					pointer = ""
				}
				buf.WriteString(fmt.Sprintf(successResponse(outputName, pointer)))
			} else {
				buf.WriteString(fmt.Sprintf("\t\treturn nil, %s{}\n", outputName))
			}
		}
	}

	if !swagger.NoSuccessType(op) {
		buf.WriteString(`
	case 400:
		var output models.DefaultBadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return nil, output

	case 500:
		var output models.DefaultInternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return nil, output

	default:
		return nil, models.DefaultInternalError{Msg: "Unknown response"}
	}
}

`)
	} else {
		buf.WriteString(`
	case 400:
		var output models.DefaultBadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return models.DefaultInternalError{Msg: err.Error()}
		}
		return output

	case 500:
		var output models.DefaultInternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return models.DefaultInternalError{Msg: err.Error()}
		}
		return output

	default:
		return models.DefaultInternalError{Msg: "Unknown response"}
	}
}

`)
	}

	return buf.String()
}

func errorMessage(err string, op *spec.Operation) string {
	if swagger.NoSuccessType(op) {
		return fmt.Sprintf(`
	if err != nil {
		return %s
	}
`, err)
	}
	return fmt.Sprintf(`
	if err != nil {
		return nil, %s
	}
`, err)
}

func successResponse(outputName, pointer string) string {
	return fmt.Sprintf(`
		var output %s
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return %soutput, nil
`, outputName, pointer)
}
