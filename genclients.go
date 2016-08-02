package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// This is extremely rough code for generating clients...
func generateClients(s Swagger) error {

	var g Generator

	// TODO: Add a general client type... something like type %sClient struct { basePath string }

	g.Printf("package generated\n\n")
	g.Printf("import \"net/http\"\n")
	g.Printf("import \"net/url\"\n")
	g.Printf("import \"encoding/json\"\n")
	g.Printf("import \"strings\"\n")
	g.Printf("import \"errors\"\n")
	g.Printf("import \"golang.org/x/net/context\"\n")
	g.Printf("import \"bytes\"\n")
	g.Printf("import \"fmt\"\n")
	g.Printf("import \"strconv\"\n")
	g.Printf("import opentracing \"github.com/opentracing/opentracing-go\"\n\n")
	// Whether we use these depends on the parameters, so we do this to prevent unused import
	// error if we don't have the right params
	g.Printf("var _ = json.Marshal\n")
	g.Printf("var _ = strings.Replace\n\n")

	for url, path := range s.Paths {
		for method, op := range path {

			// TODO: Do I really want pointers here and / or in the server?
			g.Printf("func %s(ctx context.Context, i *%sInput) (%sOutput, error) {\n",
				capitalize(op.OperationID), capitalize(op.OperationID), capitalize(op.OperationID))

			// TODO: How should I handle required fields... just check for nil pointers???

			// Build the URL
			// TODO: Make the base URL configurable...
			g.Printf("\tpath := \"http://localhost:8080\" + \"%s\"\n", s.BasePath+url)
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

`, capitalize(op.OperationID))

			// TODO: Handle the error
			g.Printf("\tresp, _ := client.Do(req)\n\n")

			// Switch on status code to build the response...
			g.Printf("\tswitch resp.StatusCode {\n")
			for key, _ := range op.Responses {

				if key == "default" {
					// TODO: Fix this... should probably factor along with server codegen
					// Or just let this slip to the default in the case statement
					continue
				}

				code, err := strconv.ParseInt(key, 10, 32)
				if err != nil {
					fmt.Errorf("Response key not valid %s", key)
				}

				g.Printf("\tcase %s:\n", key)
				if code < 400 {
					// TODO: Factor out this common code...
					outputName := fmt.Sprintf("%s%sOutput", capitalize(op.OperationID), capitalize(key))
					g.Printf(`
		var output %s
		if err := json.NewDecoder(resp.Body).Decode(&output.Data); err != nil {
			return nil, err
		}
		return output, nil
`, outputName)

				} else {
					g.Printf("\t\treturn nil, %s%sOutput{}\n", capitalize(op.OperationID), key)
				}
			}
			g.Printf("\tdefault:\n")
			g.Printf("\t\treturn nil, errors.New(\"Unknown response\")\n")
			g.Printf("\t}\n")
			g.Printf("}\n\n")
		}
	}

	return ioutil.WriteFile("generated/client.go", g.buf.Bytes(), 0644)
}

func convertParamToString(p SwaggerParameter) string {
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
