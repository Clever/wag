package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

// This is extremely rough code for generating clients...
func generateClients(s Swagger) error {

	var g Generator

	// TODO: Add a general client type... something like type %sClient struct { basePath string }

	g.Printf("package main\n\n")
	g.Printf("import \"net/http\"\n")
	g.Printf("import \"net/url\"\n")
	g.Printf("import \"encoding/json\"\n")
	g.Printf("import \"strings\"\n")
	g.Printf("import \"errors\"\n")
	g.Printf("import \"golang.org/x/net/context\"\n")
	g.Printf("import \"bytes\"\n\n")

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
			g.Printf("\tpath := \"http://localhost:8080\" + \"%s\"\n", url)
			g.Printf("\turlVals := url.Values{}\n")
			g.Printf("\tvar body []byte\n")

			for _, param := range op.Parameters {
				if param.In == "path" {
					// TODO: Should this be done with regex at some point?
					g.Printf("\tpath = strings.Replace(path, \"%s\", i.%s, -1)\n", "{"+param.Name+"}", capitalize(param.Name))
				} else if param.In == "query" {
					g.Printf("\turlVals.Add(\"%s\", i.%s)\n", param.Name, capitalize(param.Name))
				} else if param.In == "body" {
					// TODO: Handle errors here. Also, is this syntax quite right???
					g.Printf("\tbody, _ = json.Marshal(i.%s)\n", capitalize(param.Name))
				}
			}

			g.Printf("\tpath = path + \"?\" + urlVals.Encode()\n")

			g.Printf("\tclient := &http.Client{}\n")
			// TODO: Handle the error
			g.Printf("\treq, _ := http.NewRequest(\"%s\", path, bytes.NewBuffer(body))\n", method)

			for _, param := range op.Parameters {
				if param.In == "header" {
					g.Printf("\treq.Header.Set(\"%s\", i.%s)\n", param.Name, capitalize(param.Name))
				}
			}

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

				// Maybe that should just be the default of the case statement??? A general purpose error?
				g.Printf("\tcase %s:\n", key)
				if code < 400 {
					// TODO: Read the schema out here (should be in the response) and set on the output's
					// data field
					g.Printf("\t\treturn %s%sOutput{}, nil\n", capitalize(op.OperationID), key)
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
