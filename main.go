package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-openapi/spec"

	"github.com/Clever/wag/client"
	"github.com/Clever/wag/hardcoded"
	"github.com/Clever/wag/models"
	"github.com/Clever/wag/server"
	"github.com/Clever/wag/swagger"
	"github.com/Clever/wag/validation"
)

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
	swaggerSpec := *doc.Spec()

	if err := validation.Validate(swaggerSpec); err != nil {
		log.Fatalf("Swagger file not valid: %s", err)
	}

	for _, dir := range []string{"server", "client", "models"} {
		dirName := os.Getenv("GOPATH") + "/src/" + *packageName + "/" + dir
		if err := os.RemoveAll(dirName); err != nil {
			log.Fatalf("Could not remove directory: %s, error :%s", dirName, err)
		}

		if err := os.MkdirAll(dirName, 0700); err != nil {
			if !os.IsExist(err.(*os.PathError)) {
				log.Fatalf("Could not create directory: %s, error: %s", dirName, err)
			}
		}
	}

	if err := models.Generate(*packageName, swaggerSpec); err != nil {
		log.Fatalf("Error generating models: %s", err)
	}

	if err := server.Generate(*packageName, swaggerSpec); err != nil {
		log.Fatalf("Failed to generate server: %s", err)
	}

	if err := client.Generate(*packageName, swaggerSpec); err != nil {
		log.Fatalf("Failed generating clients %s", err)
	}

	middlewareGenerator := swagger.Generator{PackageName: *packageName}
	middlewareGenerator.Write(hardcoded.MustAsset("_hardcoded/middleware.go"))
	if err := middlewareGenerator.WriteFile("server/middleware.go"); err != nil {
		log.Fatalf("Failed to copy middleware.go: %s", err)
	}

	doerGenerator := swagger.Generator{PackageName: *packageName}
	doerGenerator.Write(hardcoded.MustAsset("_hardcoded/doer.go"))
	if err := doerGenerator.WriteFile("client/doer.go"); err != nil {
		log.Fatalf("Failed to copy doer.go: %s", err)
	}
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
