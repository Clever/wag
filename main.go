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

	goclient "github.com/Clever/wag/v6/clients/go"
	jsclient "github.com/Clever/wag/v6/clients/js"
	"github.com/Clever/wag/v6/hardcoded"
	"github.com/Clever/wag/v6/models"
	"github.com/Clever/wag/v6/server"
	"github.com/Clever/wag/v6/swagger"
	"github.com/Clever/wag/v6/validation"
)

// config contains the configuration derived from command line flags
type config struct {
	swaggerFile    *string
	goPackageName  *string
	outputPath     *string
	jsModulePath   *string
	versionFlag    *bool
	clientOnly     *bool
	clientLanguage *string
	goPackagePath  string
}

var version string

func main() {
	conf := &config{
		swaggerFile:    flag.String("file", "swagger.yml", "the spec file to use"),
		goPackageName:  flag.String("go-package", "", "package of the generated go code"),
		outputPath:     flag.String("output-path", "", "relative output path of the generated go code"),
		jsModulePath:   flag.String("js-path", "", "path to put the js client"),
		versionFlag:    flag.Bool("version", false, "print the wag version and exit"),
		clientOnly:     flag.Bool("client-only", false, "only generate client code"),
		clientLanguage: flag.String("client-language", "", "generate client code in specific language [go|js]"),
	}
	flag.Parse()
	if *conf.versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	if err := conf.validate(); err != nil {
		log.Fatalf(err.Error())
	}

	// Check if glide.yaml and glide.lock files are up to date
	// Ignore validation if the files don't yet exist
	glideYAMLFile, err := os.Open("glide.yaml")
	if err == nil {
		defer glideYAMLFile.Close()
		if err = validation.ValidateGlideYAML(glideYAMLFile); err != nil {
			log.Fatal(err)
		}
	}

	glideLockFile, err := os.Open("glide.lock")
	if err == nil {
		defer glideLockFile.Close()
		if err = validation.ValidateGlideLock(glideLockFile); err != nil {
			log.Fatal(err)
		}
	}

	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
	doc, err := loads.Spec(*conf.swaggerFile)
	if err != nil {
		log.Fatalf("Error loading swagger file: %s", err)
	}
	swaggerSpec := *doc.Spec()

	if err := validation.Validate(*doc); err != nil {
		log.Fatalf("Swagger file not valid: %s", err)
	}

	dirs := []string{
		filepath.Join(os.Getenv("GOPATH"), "src", conf.goPackagePath, "server"),
		filepath.Join(os.Getenv("GOPATH"), "src", conf.goPackagePath, "client"),
		filepath.Join(os.Getenv("GOPATH"), "src", conf.goPackagePath, "models"),
		*conf.jsModulePath,
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

	if err := models.Generate(conf.goPackagePath, swaggerSpec); err != nil {
		log.Fatalf("Error generating models: %s", err)
	}

	if err := server.Generate(*conf.goPackageName, conf.goPackagePath, swaggerSpec); err != nil {
		log.Fatalf("Failed to generate server: %s", err)
	}

	if err := goclient.Generate(*conf.goPackageName, conf.goPackagePath, swaggerSpec); err != nil {
		log.Fatalf("Failed generating go client %s", err)
	}

	if err := jsclient.Generate(*conf.jsModulePath, swaggerSpec); err != nil {
		log.Fatalf("Failed generating js client %s", err)
	}

	middlewareGenerator := swagger.Generator{PackagePath: conf.goPackagePath}
	middlewareGenerator.Write(hardcoded.MustAsset("../_hardcoded/middleware.go"))
	if err := middlewareGenerator.WriteFile("server/middleware.go"); err != nil {
		log.Fatalf("Failed to copy middleware.go: %s", err)
	}

	doerGenerator := swagger.Generator{PackagePath: conf.goPackagePath}
	doerGenerator.Write(hardcoded.MustAsset("../_hardcoded/doer.go"))
	if err := doerGenerator.WriteFile("client/doer.go"); err != nil {
		log.Fatalf("Failed to copy doer.go: %s", err)
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
// Example: if packagePath = github.com/Clever/wag/v6/gen-go and the module name is github.com/Clever/wag/v6/v2
// the function will return github.com/Clever/wag/v6/v2/gen-go
// Example: if packagePath = github.com/Clever/wag/v6/gen-go and the module name is github.com/Clever/wag/v6
// the function will return  github.com/Clever/wag/v6/gen-go
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

func (c *config) validate() error {
	// check if the repo uses modules
	if modFile, err := os.Open("go.mod"); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return fmt.Errorf("Error checking if the repo contains a go.mod file: %s", err.Error())
		}
		if *c.goPackageName == "" {
			return fmt.Errorf("go-package is required")
		}
		c.goPackagePath = *c.goPackageName
	} else {
		defer modFile.Close()

		goPath := os.Getenv("GOPATH")
		if goPath == "" {
			return fmt.Errorf("GOPATH must be set")
		}
		if *c.outputPath == "" {
			return fmt.Errorf("output-path is required")
		}

		*c.outputPath = path.Clean(*c.outputPath)
		c.goPackagePath = getModulePackagePath(goPath, *c.outputPath)
		*c.goPackageName = getModulePackageName(modFile, *c.outputPath)
	}

	if *c.jsModulePath == "" {
		return fmt.Errorf("js-path is required")
	}

	return nil
}
