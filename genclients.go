package main

import (
	"fmt"
	"strings"

	"github.com/go-openapi/spec"
)

// This is extremely rough code for generating clients...
func generateClients(packageName string, s spec.Swagger) error {

	var g Generator

	g.Printf("package client\n\n")
	g.Printf(importStatements([]string{"golang.org/x/net/context", "strings", "bytes",
		"net/http", "net/url", "strconv", "encoding/json", "strconv", packageName + "/models"}))

	g.Printf(`var _ = json.Marshal
var _ = strings.Replace

var _ = strconv.FormatInt

type Client struct {
	BasePath    string
	requestDoer doer
	transport   *http.Transport
	// Keep the retry doer around so that we can set the number of retries
	retryDoer retryDoer
}

// New creates a new client. The base path and http transport are configurable
func New(basePath string) Client {
	base := baseDoer{}
	tracing := tracingDoer{d: base}
	retry := retryDoer{d: tracing, defaultRetries: 1}

	return Client{requestDoer: retry, retryDoer: retry, transport: &http.Transport{}, BasePath: basePath}
}

func (c Client) WithRetries(retries int) Client {
	c.retryDoer.defaultRetries = retries
	return c
}

`)

	for _, path := range sortedPathItemKeys(s.Paths.Paths) {
		pathItem := s.Paths.Paths[path]
		pathItemOps := pathItemOperations(pathItem)
		for _, method := range sortedOperationsKeys(pathItemOps) {
			op := pathItemOps[method]

			// TODO: Do I really want pointers here and / or in the server?
			g.Printf("func (c Client) %s(ctx context.Context, i *models.%sInput) (models.%sOutput, error) {\n",
				capitalize(op.ID), capitalize(op.ID), capitalize(op.ID))

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
						"{"+param.Name+"}", convertParamToString(param))
				} else if param.In == "query" {
					g.Printf("\turlVals.Add(\"%s\", %s)\n",
						param.Name, convertParamToString(param))
				}
			}

			g.Printf("\tpath = path + \"?\" + urlVals.Encode()\n\n")

			for _, param := range op.Parameters {
				if param.In == "body" {
					// TODO: Handle errors here. Also, is this syntax quite right???
					g.Printf("\tbody, _ = json.Marshal(i.%s)\n\n", capitalize(param.Name))
				}
			}

			g.Printf("\tclient := &http.Client{Transport: c.transport}\n")
			// TODO: Handle the error
			g.Printf("\treq, _ := http.NewRequest(\"%s\", path, bytes.NewBuffer(body))\n", strings.ToUpper(method))

			for _, param := range op.Parameters {
				if param.In == "header" {
					g.Printf("\treq.Header.Set(\"%s\", %s)\n",
						param.Name, convertParamToString(param))
				}
			}

			g.Printf(`
	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "%s")
	resp, err := c.requestDoer.Do(ctx, client, req)
	if err != nil {
		return nil, models.DefaultInternalError{Msg: err.Error()}
	}
`, op.ID)

			// Switch on status code to build the response...
			g.Printf("\tswitch resp.StatusCode {\n")

			for _, statusCode := range sortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
				response := op.Responses.StatusCodeResponses[statusCode]

				g.Printf("\tcase %d:\n", statusCode)

				if response.Schema == nil {
					g.Printf("\t\tvar output models.%s%dOutput\n", capitalize(op.ID), statusCode)
					if statusCode < 400 {
						g.Printf("\t\treturn output, nil\n")
					} else {
						g.Printf("\t\treturn nil, output\n")
					}
				} else {
					if statusCode < 400 {
						// TODO: Factor out this common code...
						outputName := fmt.Sprintf("models.%s%dOutput", capitalize(op.ID), statusCode)
						g.Printf(successResponse(outputName))
					} else {
						g.Printf("\t\treturn nil, models.%s%dOutput{}\n", capitalize(op.ID), statusCode)
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

	return g.WriteFile("generated/client/client.go")
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

func convertParamToString(p spec.Parameter) string {
	switch p.Type {
	case "string":
		return fmt.Sprintf("i.%s", capitalize(p.Name))
	case "integer":
		return fmt.Sprintf("strconv.FormatInt(i.%s, 10)", capitalize(p.Name))
	case "number":
		return fmt.Sprintf("strconv.FormatFloat(i.%s, 'E', -1, 64)", capitalize(p.Name))
	case "boolean":
		return fmt.Sprintf("strconv.FormatBool(i.%s)", capitalize(p.Name))
	default:
		// Theoretically should have validated before getting here
		panic(fmt.Errorf("Unsupported parameter type %s", p.Type))
	}
}
