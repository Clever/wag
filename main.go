package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"

	goclient "github.com/Clever/wag/v5/clients/go"
	jsclient "github.com/Clever/wag/v5/clients/js"
	"github.com/Clever/wag/v5/hardcoded"
	"github.com/Clever/wag/v5/models"
	"github.com/Clever/wag/v5/server"
	"github.com/Clever/wag/v5/swagger"
	"github.com/Clever/wag/v5/validation"
)

var version string

func main() {

	swaggerFile := flag.String("file", "swagger.yml", "the spec file to use")
	goPackageName := flag.String("go-package", "", "package of the generated go code")
	outputPath := flag.String("output-path", "", "relative output path of the generated go code")
	jsModulePath := flag.String("js-path", "", "path to put the js client")
	versionFlag := flag.Bool("version", false, "print the wag version and exit")
	flag.Parse()
	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	var goPackagePath string
	// check if the repo uses modules
	if modFile, err := os.Open("go.mod"); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			log.Fatalf("Error checking if the repo contains a go.mod file: %s", err.Error())
		}
		if *goPackageName == "" {
			log.Fatal("go-package is required")
		}
		goPackagePath = *goPackageName
	} else {
		defer modFile.Close()

		goPath := os.Getenv("GOPATH")
		if goPath == "" {
			log.Fatalf("GOPATH must be set")
		}
		if *outputPath == "" {
			log.Fatal("output-path is required")
		}

		*outputPath = path.Clean(*outputPath)
		goPackagePath = getModulePackagePath(goPath, *outputPath)
		*goPackageName = getModulePackageName(modFile, *outputPath)
	}

	if *jsModulePath == "" {
		log.Fatal("js-path is required")
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
		filepath.Join(os.Getenv("GOPATH"), "src", goPackagePath, "server"),
		filepath.Join(os.Getenv("GOPATH"), "src", goPackagePath, "client"),
		filepath.Join(os.Getenv("GOPATH"), "src", goPackagePath, "models"),
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

	if err := models.Generate(goPackagePath, swaggerSpec); err != nil {
		log.Fatalf("Error generating models: %s", err)
	}

	if err := server.Generate(*goPackageName, goPackagePath, swaggerSpec); err != nil {
		log.Fatalf("Failed to generate server: %s", err)
	}

	if err := goclient.Generate(*goPackageName, goPackagePath, swaggerSpec); err != nil {
		log.Fatalf("Failed generating go client %s", err)
	}

	if err := jsclient.Generate(*jsModulePath, swaggerSpec); err != nil {
		log.Fatalf("Failed generating js client %s", err)
	}

	middlewareGenerator := swagger.Generator{PackagePath: goPackagePath}
	middlewareGenerator.Write(hardcoded.MustAsset("../_hardcoded/middleware.go"))
	if err := middlewareGenerator.WriteFile("server/middleware.go"); err != nil {
		log.Fatalf("Failed to copy middleware.go: %s", err)
	}

	doerGenerator := swagger.Generator{PackagePath: goPackagePath}
	doerGenerator.Write(hardcoded.MustAsset("../_hardcoded/doer.go"))
	if err := doerGenerator.WriteFile("client/doer.go"); err != nil {
		log.Fatalf("Failed to copy doer.go: %s", err)
	}

	tracingGenerator := swagger.Generator{PackagePath: goPackagePath}
	tracingGenerator.Write(hardcoded.MustAsset("../_hardcoded/tracing.go"))
	if err := tracingGenerator.WriteFile("tracing/tracing.go"); err != nil {
		log.Fatalf("Failed to copy tracing.go: %s", err)
	}
}

func getModulePackagePath(goPath, outputPath string) string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %s", err.Error())
	}
	goSrcPath := fmt.Sprintf("%v%v", goPath, "/src/")
	return path.Join(strings.TrimPrefix(pwd, goSrcPath), outputPath)
}

// getModulePackageName gets the package name of the generated code
// Example: if packagePath = github.com/Clever/wag/v5/gen-go and the module name is github.com/Clever/wag/v5/v2
// the function will return github.com/Clever/wag/v5/v2/gen-go
// Example: if packagePath = github.com/Clever/wag/v5/gen-go and the module name is github.com/Clever/wag/v5
// the function will return  github.com/Clever/wag/v5/gen-go
func getModulePackageName(modFile *os.File, outputPath string) string {
	// read first line of module file
	r := bufio.NewReader(modFile)
	b, _, err := r.ReadLine()
	if err != nil {
		log.Fatalf("Error checking module name: %s", err.Error())
	}

	// parse module path
	moduleName := strings.TrimPrefix(string(b), "module")
	moduleName = strings.TrimSpace(moduleName)
	return fmt.Sprintf("%v/%v", moduleName, outputPath)
}
