package client

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/Clever/wag/models"
	"github.com/Clever/wag/swagger"
)

// Generate generates a client
func Generate(packageName string, s spec.Swagger) error {

	g := swagger.Generator{PackageName: packageName}

	g.Printf(`package client

import (
		"golang.org/x/net/context"
		"strings"
		"bytes"
		"net/http"
		"net/url"
		"strconv"
		"encoding/json"
		"strconv"
		"%s/models"
)

var _ = json.Marshal
var _ = strings.Replace

var _ = strconv.FormatInt

type Client struct {
	BasePath    string
	requestDoer doer
	transport   *http.Transport
	// Keep the retry doer around so that we can set the number of retries
	retryDoer *retryDoer
}

// New creates a new client. The base path and http transport are configurable
func New(basePath string) Client {
	base := baseDoer{}
	tracing := tracingDoer{d: base}
	retry := retryDoer{d: tracing, defaultRetries: 1}

	return Client{requestDoer: &retry, retryDoer: &retry, transport: &http.Transport{}, BasePath: basePath}
}

func (c Client) WithRetries(retries int) Client {
	c.retryDoer.defaultRetries = retries
	return c
}

`, packageName)

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

func methodCode(op *spec.Operation, basePath, method, methodPath string) string {
	var buf bytes.Buffer
	capOpID := swagger.Capitalize(op.ID)

	buf.WriteString(fmt.Sprintf("func (c Client) %s(ctx context.Context, i *models.%sInput) (models.%sOutput, error) {\n",
		capOpID, capOpID, capOpID))

	buf.WriteString(fmt.Sprintf("\tpath := c.BasePath + \"%s\"\n", basePath+methodPath))
	buf.WriteString(fmt.Sprintf("\turlVals := url.Values{}\n"))
	buf.WriteString(fmt.Sprintf("\tvar body []byte\n\n"))

	buf.WriteString(fmt.Sprintf(buildRequestCode(op, method)))

	buf.WriteString(fmt.Sprintf(`
	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "%s")
	resp, err := c.requestDoer.Do(client, req.WithContext(ctx))
	if err != nil {
		return nil, models.DefaultInternalError{Msg: err.Error()}
	}
`, op.ID))

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
				buf.WriteString(fmt.Sprintf("\tif i.%s != nil {\n", swagger.Capitalize(param.Name)))
				buf.WriteString(fmt.Sprintf(queryAddCode))
				buf.WriteString(fmt.Sprintf("\t}\n"))
			}
		}
	}

	buf.WriteString(fmt.Sprintf("\tpath = path + \"?\" + urlVals.Encode()\n\n"))

	for _, param := range op.Parameters {
		if param.In == "body" {
			bodyMarshalCode := fmt.Sprintf(`
	var err error
	body, err = json.Marshal(i.%s)
	if err != nil {
		return nil, err
	}
`, swagger.Capitalize(param.Name))
			if param.Required {
				buf.WriteString(fmt.Sprintf(bodyMarshalCode))
			} else {
				buf.WriteString(fmt.Sprintf("\tif i.%s != nil {\n", swagger.Capitalize(param.Name)))
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
	if err != nil {

		return nil, err
	}
`, strings.ToUpper(method)))

	for _, param := range op.Parameters {
		if param.In == "header" {
			headerAddCode := fmt.Sprintf("\treq.Header.Set(\"%s\", %s)\n", param.Name, swagger.ParamToStringCode(param))
			if param.Required {
				buf.WriteString(fmt.Sprintf(headerAddCode))
			} else {
				buf.WriteString(fmt.Sprintf("\tif i.%s != nil {\n", swagger.Capitalize(param.Name)))
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

	for _, statusCode := range swagger.SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
		response := op.Responses.StatusCodeResponses[statusCode]
		outputName := "models." + models.OutputObjectName(op, statusCode)

		buf.WriteString(fmt.Sprintf("\tcase %d:\n", statusCode))

		if response.Schema == nil {
			buf.WriteString(fmt.Sprintf("\t\tvar output %s\n", outputName))
			if statusCode < 400 {
				buf.WriteString(fmt.Sprintf("\t\treturn output, nil\n"))
			} else {
				buf.WriteString(fmt.Sprintf("\t\treturn nil, output\n"))
			}
		} else {
			if statusCode < 400 {
				buf.WriteString(fmt.Sprintf(successResponse(outputName)))
			} else {
				buf.WriteString(fmt.Sprintf("\t\treturn nil, %s{}\n", outputName))
			}
		}
	}

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
	return buf.String()
}

func successResponse(outputName string) string {
	return fmt.Sprintf(`
		var output %s
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return output, nil
`, outputName)
}
