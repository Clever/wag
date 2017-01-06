package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-openapi/spec"

	"github.com/Clever/wag/clients/go"
	"github.com/Clever/wag/clients/js"
	"github.com/Clever/wag/hardcoded"
	"github.com/Clever/wag/models"
	"github.com/Clever/wag/server"
	"github.com/Clever/wag/swagger"
	"github.com/Clever/wag/validation"
)

func capitalize(input string) string {
	return strings.ToUpper(input[0:1]) + input[1:]
}

var version string

func main() {

	swaggerFile := flag.String("file", "swagger.yml", "the spec file to use")
	goPackageName := flag.String("go-package", "", "package of the generated go code")
	jsModulePath := flag.String("js-path", "", "path to put the js client")
	versionFlag := flag.Bool("version", false, "print the wag version and exit")
	flag.Parse()
	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}
	if *goPackageName == "" {
		log.Fatal("go-package is required")
	}
	if *jsModulePath == "" {
		log.Fatal("js-path is required")
	}

	if glideYMLFile, err := os.Open("glide.yaml"); err == nil {
		if err := validation.ValidateGlideYML(glideYMLFile); err != nil {
			log.Fatal(err)
		}
	}

	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
	doc, err := loads.Spec(*swaggerFile)
	if err != nil {
		log.Fatalf("Error loading swagger file: %s", err)
	}
	swaggerSpec := *doc.Spec()

	if err := validation.Validate(*doc); err != nil {
		log.Fatalf("Swagger file not valid: %s", err)
	}

	dirs := []string{
		filepath.Join(os.Getenv("GOPATH"), "src", *goPackageName, "server"),
		filepath.Join(os.Getenv("GOPATH"), "src", *goPackageName, "client"),
		filepath.Join(os.Getenv("GOPATH"), "src", *goPackageName, "models"),
		*jsModulePath,
	}

	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			log.Fatalf("Could not remove directory: %s, error :%s", dir, err)
		}

		if err := os.MkdirAll(dir, 0700); err != nil {
			if !os.IsExist(err.(*os.PathError)) {
				log.Fatalf("Could not create directory: %s, error: %s", dir, err)
			}
		}
	}

	err = swagger.ValidateResponses(swaggerSpec)
	if err != nil {
		log.Fatalf("Failed processing the swagger spec: %s", err)
	}

	if err := models.Generate(*goPackageName, swaggerSpec); err != nil {
		log.Fatalf("Error generating models: %s", err)
	}

	if err := server.Generate(*goPackageName, swaggerSpec); err != nil {
		log.Fatalf("Failed to generate server: %s", err)
	}

	if err := goclient.Generate(*goPackageName, swaggerSpec); err != nil {
		log.Fatalf("Failed generating go client %s", err)
	}

	if err := jsclient.Generate(*jsModulePath, swaggerSpec); err != nil {
		log.Fatalf("Failed generating js client %s", err)
	}

	middlewareGenerator := swagger.Generator{PackageName: *goPackageName}
	middlewareGenerator.Write(hardcoded.MustAsset("_hardcoded/middleware.go"))
	if err := middlewareGenerator.WriteFile("server/middleware.go"); err != nil {
		log.Fatalf("Failed to copy middleware.go: %s", err)
	}

	doerGenerator := swagger.Generator{PackageName: *goPackageName}
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
