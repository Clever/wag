package main

import (
	"io/ioutil"
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
	g.Printf("import \"golang.org/x/net/context\"\n")
	g.Printf("import \"bytes\"\n\n")

	// Whether we use these depends on the parameters, so we do this to prevent unused import
	// error if we don't have the right params
	g.Printf("var _ = json.Marshal\n")
	g.Printf("var _ = strings.Replace\n\n")

	for url, path := range s.Paths {
		for method, op := range path {

			// TODO: Do I really want pointers here and / or in the server?
			g.Printf("func %s(ctx context.Context, i *%sInput) {\n",
				capitalize(op.OperationID), capitalize(op.OperationID))

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
			// TODO: Encode the URL correctly and put in path / query params
			g.Printf("\treq, _ := http.NewRequest(\"%s\", path, bytes.NewBuffer(body))\n", method)

			// TODO: Handle non-required...
			for _, param := range op.Parameters {
				if param.In == "header" {
					g.Printf("\treq.Header.Set(\"%s\", i.%s)\n", param.Name, capitalize(param.Name))
				}
			}

			g.Printf("\tclient.Do(req)\n\n")
			g.Printf("}\n\n")
		}
	}

	return ioutil.WriteFile("generated/client.go", g.buf.Bytes(), 0644)
}
