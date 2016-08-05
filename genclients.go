package main

import (
	"fmt"
	"strings"

	"github.com/go-openapi/spec"
)

// This is extremely rough code for generating clients...
func generateClients(packageName string, s spec.Swagger) error {

	var g Generator

	// TODO: Add a general client type... something like type %sClient struct { basePath string }

	g.Printf("package client\n\n")
	g.Printf("import \"net/http\"\n")
	g.Printf("import \"net/url\"\n")
	g.Printf("import \"encoding/json\"\n")
	g.Printf("import \"strings\"\n")
	g.Printf("import \"golang.org/x/net/context\"\n")
	g.Printf("import \"bytes\"\n")
	g.Printf("import \"fmt\"\n")
	g.Printf("import \"strconv\"\n")
	g.Printf("import \"%s/models\"\n\n", packageName)
	g.Printf("import opentracing \"github.com/opentracing/opentracing-go\"\n\n")
	// Whether we use these depends on the parameters, so we do this to prevent unused import
	// error if we don't have the right params
	g.Printf("var _ = json.Marshal\n")
	g.Printf("var _ = strings.Replace\n\n")
	g.Printf("var _ = strconv.FormatInt\n\n")

	for _, path := range sortedPathItemKeys(s.Paths.Paths) {
		pathItem := s.Paths.Paths[path]
		pathItemOps := pathItemOperations(pathItem)
		for _, method := range sortedOperationsKeys(pathItemOps) {
			op := pathItemOps[method]

			// TODO: Do I really want pointers here and / or in the server?
			g.Printf("func %s(ctx context.Context, i *models.%sInput) (models.%sOutput, error) {\n",
				capitalize(op.ID), capitalize(op.ID), capitalize(op.ID))

			// TODO: How should I handle required fields... just check for nil pointers???

			// Build the URL
			// TODO: Make the base URL configurable...
			g.Printf("\tpath := \"http://localhost:8080\" + \"%s\"\n", s.BasePath+path)
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

			g.Printf("\tclient := &http.Client{}\n")
			// TODO: Handle the error
			g.Printf("\treq, _ := http.NewRequest(\"%s\", path, bytes.NewBuffer(body))\n", strings.ToUpper(method))

			for _, param := range op.Parameters {
				if param.In == "header" {
					g.Printf("\treq.Header.Set(\"%s\", %s)\n",
						param.Name, convertParamToString(param))
				}
			}

			// Inject tracing headers
			g.Printf(`
	// Inject tracing headers
	opName := "%s"
	var sp opentracing.Span
	// TODO: add tags relating to input data?
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		sp = opentracing.StartSpan(opName, opentracing.ChildOf(parentSpan.Context()))
	} else {
		sp = opentracing.StartSpan(opName)
	}
	if err := sp.Tracer().Inject(sp.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header)); err != nil {
		return nil, fmt.Errorf("couldn't inject tracing headers (%%v)", err)
	}

`, capitalize(op.ID))

			// TODO: Handle the error
			g.Printf("\tresp, _ := client.Do(req)\n\n")

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
