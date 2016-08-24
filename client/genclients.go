package client

import (
	"fmt"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/Clever/wag/swagger"
)

// Generate generates a client
func Generate(packageName string, s spec.Swagger) error {

	g := swagger.Generator{PackageName: packageName}

	g.Printf("package client\n\n")
	g.Printf(swagger.ImportStatements([]string{"golang.org/x/net/context", "strings", "bytes",
		"net/http", "net/url", "strconv", "encoding/json", "strconv", packageName + "/models"}))

	g.Printf(`var _ = json.Marshal
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

`)

	g.Printf(swagger.BaseParamToStringCode())

	for _, path := range swagger.SortedPathItemKeys(s.Paths.Paths) {
		pathItem := s.Paths.Paths[path]
		pathItemOps := swagger.PathItemOperations(pathItem)
		for _, method := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[method]
			capOpID := swagger.Capitalize(op.ID)

			// TODO: Do I really want pointers here and / or in the server?
			g.Printf("func (c Client) %s(ctx context.Context, i *models.%sInput) (models.%sOutput, error) {\n",
				capOpID, capOpID, capOpID)

			// TODO: How should I handle required fields... just check for nil pointers???

			// Build the URL
			// TODO: Make the base URL configurable...
			g.Printf("\tpath := c.BasePath + \"%s\"\n", s.BasePath+path)
			g.Printf("\turlVals := url.Values{}\n")
			g.Printf("\tvar body []byte\n\n")

			for _, param := range op.Parameters {
				if param.In == "path" {
					// TODO: Should this be done with regex at some point?
					g.Printf("\tpath = strings.Replace(path, \"%s\", %s, -1)\n",
						"{"+param.Name+"}", swagger.ParamToStringCode(param))

				} else if param.In == "query" {
					queryAddCode := fmt.Sprintf("\turlVals.Add(\"%s\", %s)\n", param.Name, swagger.ParamToStringCode(param))
					if param.Required {
						g.Printf(queryAddCode)
					} else {
						g.Printf("\tif i.%s != nil {\n", swagger.Capitalize(param.Name))
						g.Printf(queryAddCode)
						g.Printf("\t}\n")
					}
				}
			}

			g.Printf("\tpath = path + \"?\" + urlVals.Encode()\n\n")

			for _, param := range op.Parameters {
				if param.In == "body" {
					// TODO: Handle errors here. Also, is this syntax quite right???
					bodyMarshalCode := fmt.Sprintf("\tbody, _ = json.Marshal(i.%s)\n\n", swagger.Capitalize(param.Name))
					if param.Required {
						g.Printf(bodyMarshalCode)
					} else {
						g.Printf("\tif i.%s != nil {\n", swagger.Capitalize(param.Name))
						g.Printf(bodyMarshalCode)
						g.Printf("\t}\n")
					}
				}
			}

			g.Printf("\tclient := &http.Client{Transport: c.transport}\n")
			// TODO: Handle the error
			g.Printf("\treq, _ := http.NewRequest(\"%s\", path, bytes.NewBuffer(body))\n", strings.ToUpper(method))

			for _, param := range op.Parameters {
				if param.In == "header" {
					headerAddCode := fmt.Sprintf("\treq.Header.Set(\"%s\", %s)\n", param.Name, swagger.ParamToStringCode(param))
					if param.Required {
						g.Printf(headerAddCode)
					} else {
						g.Printf("\tif i.%s != nil {\n", swagger.Capitalize(param.Name))
						g.Printf(headerAddCode)
						g.Printf("\t}\n")
					}
				}
			}

			g.Printf(`
	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "%s")
	resp, err := c.requestDoer.Do(client, req.WithContext(ctx))
	if err != nil {
		return nil, models.DefaultInternalError{Msg: err.Error()}
	}
`, op.ID)

			// Switch on status code to build the response...
			g.Printf("\tswitch resp.StatusCode {\n")

			for _, statusCode := range swagger.SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
				response := op.Responses.StatusCodeResponses[statusCode]

				g.Printf("\tcase %d:\n", statusCode)

				if response.Schema == nil {
					g.Printf("\t\tvar output models.%s%dOutput\n", capOpID, statusCode)
					if statusCode < 400 {
						g.Printf("\t\treturn output, nil\n")
					} else {
						g.Printf("\t\treturn nil, output\n")
					}
				} else {
					if statusCode < 400 {
						// TODO: Factor out this common code...
						outputName := fmt.Sprintf("models.%s%dOutput", capOpID, statusCode)
						g.Printf(successResponse(outputName))
					} else {
						g.Printf("\t\treturn nil, models.%s%dOutput{}\n", capOpID, statusCode)
					}
				}
			}

			// Add in the default 400, 500 responses
			g.Printf("\tcase 400:\n")
			g.Printf(badRequestCode)
			g.Printf("\tcase 500:\n")
			g.Printf(internalErrorCode)

			g.Printf("\tdefault:\n")
			g.Printf("\t\treturn nil, models.DefaultInternalError{Msg: \"Unknown response\"}\n")
			g.Printf("\t}\n")
			g.Printf("}\n\n")
		}
	}

	return g.WriteFile("client/client.go")
}

var badRequestCode = `
		var output models.DefaultBadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return nil, output
`

var internalErrorCode = `
		var output models.DefaultInternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return nil, output
`

func successResponse(outputName string) string {
	return fmt.Sprintf(`
		var output %s
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return output, nil
`, outputName)
}
