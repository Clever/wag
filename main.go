package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-openapi/spec"

	"github.com/Clever/wag/client"
	"github.com/Clever/wag/models"
	"github.com/Clever/wag/server"
)

func pathItemOperations(item spec.PathItem) map[string]*spec.Operation {
	ops := make(map[string]*spec.Operation)
	if item.Get != nil {
		ops["GET"] = item.Get
	}
	if item.Put != nil {
		ops["PUT"] = item.Put
	}
	if item.Post != nil {
		ops["POST"] = item.Post
	}
	if item.Delete != nil {
		ops["DELETE"] = item.Delete
	}
	if item.Options != nil {
		ops["OPTIONS"] = item.Options
	}
	if item.Head != nil {
		ops["HEAD"] = item.Head
	}
	if item.Patch != nil {
		ops["PATCH"] = item.Patch
	}
	return ops
}

func sliceContains(slice []string, key string) bool {
	for _, val := range slice {
		if val == key {
			return true
		}
	}
	return false
}

func capitalize(input string) string {
	return strings.ToUpper(input[0:1]) + input[1:]
}

func main() {

	swaggerFile := flag.String("file", "swagger.yml", "the spec file to use")
	packageName := flag.String("package", "", "package of the generated code")
	flag.Parse()
	if *packageName == "" {
		log.Fatal("package is required")
	}

	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
	doc, err := loads.Spec(*swaggerFile)
	if err != nil {
		log.Fatalf("Error loading swagger file: %s", err)
	}
	swagger := *doc.Spec()

	if err := validate(swagger); err != nil {
		log.Fatalf("Swagger file not valid: %s", err)
	}

	for _, dir := range []string{"server", "client", "models"} {
		dirName := os.Getenv("GOPATH") + "/src/" + *packageName + "/" + dir
		if err := os.Mkdir(dirName, 0700); err != nil {
			if !os.IsExist(err.(*os.PathError)) {
				log.Fatalf("Could not create directory: %s, error: %s", dirName, err)
			}
		}
	}

	if err := models.Generate(*packageName, *swaggerFile, swagger); err != nil {
		log.Fatalf("Error generating models: %s", err)
	}

	if err := server.Generate(*packageName, swagger); err != nil {
		log.Fatalf("Failed to generate server: %s", err)
	}

	if err := client.Generate(*packageName, swagger); err != nil {
		log.Fatalf("Failed generating clients %s", err)
	}

	if err := ioutil.WriteFile("generated/server/middleware.go", MustAsset("hardcoded/_middleware.go"), 0644); err != nil {
		log.Fatalf("Failed to copy middleware.go: %s", err)
	}

	if err := ioutil.WriteFile("generated/client/doer.go", MustAsset("hardcoded/_doer.go"), 0644); err != nil {
		log.Fatalf("Failed to copy doer.go: %s", err)
	}
}

type Generator struct {
	buf bytes.Buffer
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) WriteFile(path string) error {
	fileBytes, err := format.Source(g.buf.Bytes())
	if err != nil {
		// This will error if the code isn't valid so let's print it to make it
		// easier to debug
		fmt.Printf("BAD CODE\n%s\n", string(g.buf.Bytes()))
		return err
	}
	return ioutil.WriteFile(path, fileBytes, 0644)
}

func sortedPathItemKeys(m map[string]spec.PathItem) []string {
	sortedKeys := []string{}
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

func sortedOperationsKeys(m map[string]*spec.Operation) []string {
	sortedKeys := []string{}
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

func sortedStatusCodeKeys(m map[int]spec.Response) []int {
	sortedKeys := []int{}
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)
	return sortedKeys
}

// importStatements takes a list of import strings and converts them to a formatted imports
func importStatements(imports []string) string {
	if len(imports) == 0 {
		return ""
	}
	output := "import (\n"
	for _, importStr := range imports {
		output += fmt.Sprintf("\t\"%s\"\n", importStr)
	}
	output += ")\n\n"
	return output
}

func typeFromSchema(schema *spec.Schema) (string, error) {
	// We support three types of schemas
	// 1. An empty schema
	// 2. A schema with one element, the $ref key
	// 3. A schema with two elements. One a type with value 'array' and another items field
	// referencing the $ref
	// TODO: The error messages here aren't great...
	if schema == nil {
		// represent this as a string, which is empty by default
		return "string", nil
	} else if schema.Ref.String() != "" {
		ref := schema.Ref.String()
		if !strings.HasPrefix(ref, "#/definitions/") {
			return "", fmt.Errorf("schema.$ref has undefined reference type. Must be #/definitions")
		}
		return ref[len("#/definitions/"):], nil
	} else {
		schemaType := schema.Type
		if len(schemaType) != 1 || schemaType[0] != "array" {
			return "", fmt.Errorf("Two element schemas must have a 'type' field with the value 'array'")
		}
		items := schema.Items
		if items == nil || items.Schema == nil {
			return "", fmt.Errorf("Two element schemas must have an '$ref' field in the 'items' descriptions")
		}
		ref := items.Schema.Ref.String()
		if !strings.HasPrefix(ref, "#/definitions/") {
			return "", fmt.Errorf("schema.$ref has undefined reference type. Must be #/definitions")
		}
		return "[]" + ref[len("#/definitions/"):], nil
	}
}
