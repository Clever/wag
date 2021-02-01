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
	"github.com/go-openapi/spec"

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
	// flag values
	swaggerFile    *string
	goPackageName  *string
	outputPath     *string
	jsModulePath   *string
	versionFlag    *bool
	clientOnly     *bool
	clientLanguage *string

	// dervied values
	goPackagePath    string
	generateServer   bool
	generateGoClient bool
	generateGoModels bool
	generateJSClient bool
	modelsPath       string
	serverPath       string
	goClientPath     string
	jsClientPath     string
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

	err = swagger.ValidateResponses(swaggerSpec)
	if err != nil {
		log.Fatalf("Failed processing the swagger spec: %s", err)
	}

	if conf.generateGoModels {
		if err := generateGoModels(conf.modelsPath, conf.goPackagePath, swaggerSpec); err != nil {
			log.Fatal(err.Error())
		}
	}

	if conf.generateServer {
		if err := generateServer(conf.serverPath, *conf.goPackageName, conf.goPackagePath, swaggerSpec); err != nil {
			log.Fatal(err.Error())
		}
	}

	if conf.generateGoClient {
		if err := generateGoClient(conf.goClientPath, *conf.goPackageName, conf.goPackagePath, swaggerSpec); err != nil {
			log.Fatal(err.Error())
		}
	}

	if conf.generateJSClient {
		if err := generateJSClient(*conf.jsModulePath, swaggerSpec); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func generateGoModels(modelsPath, goPackagePath string, swaggerSpec spec.Swagger) error {
	if err := prepareDir(modelsPath); err != nil {
		return err
	}
	if err := models.Generate(goPackagePath, swaggerSpec); err != nil {
		return fmt.Errorf("Error generating models: %s", err)
	}
	return nil
}

func generateServer(serverPath, goPackageName, goPackagePath string, swaggerSpec spec.Swagger) error {
	if err := prepareDir(serverPath); err != nil {
		return err
	}
	if err := server.Generate(goPackageName, goPackagePath, swaggerSpec); err != nil {
		return fmt.Errorf("Failed to generate server: %s", err)
	}
	middlewareGenerator := swagger.Generator{PackagePath: goPackagePath}
	middlewareGenerator.Write(hardcoded.MustAsset("../_hardcoded/middleware.go"))
	if err := middlewareGenerator.WriteFile("server/middleware.go"); err != nil {
		return fmt.Errorf("Failed to copy middleware.go: %s", err)
	}
	return nil
}

func generateGoClient(goClientPath, goPackageName, goPackagePath string, swaggerSpec spec.Swagger) error {
	if err := prepareDir(goClientPath); err != nil {
		return err
	}
	if err := goclient.Generate(goPackageName, goPackagePath, swaggerSpec); err != nil {
		return fmt.Errorf("Failed generating go client %s", err)
	}
	doerGenerator := swagger.Generator{PackagePath: goPackagePath}
	doerGenerator.Write(hardcoded.MustAsset("../_hardcoded/doer.go"))
	if err := doerGenerator.WriteFile("client/doer.go"); err != nil {
		return fmt.Errorf("Failed to copy doer.go: %s", err)
	}
	return nil
}

func generateJSClient(jsModulePath string, swaggerSpec spec.Swagger) error {
	if err := prepareDir(jsModulePath); err != nil {
		return err
	}
	if err := jsclient.Generate(jsModulePath, swaggerSpec); err != nil {
		return fmt.Errorf("Failed generating js client %s", err)
	}
	return nil
}

func prepareDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("Could not remove directory: %s, error :%s", dir, err)
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		if !os.IsExist(err.(*os.PathError)) {
			return fmt.Errorf("Could not create directory: %s, error: %s", dir, err)
		}
	}
	return nil
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

// validate validates the configuration and determines the goPackagePath for the generated code
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

	// parse flags determing if client/server code is generated
	c.generateServer = true
	c.generateGoClient = true
	c.generateJSClient = true
	// only generate client
	if c.clientOnly != nil && *c.clientOnly == true {
		c.generateServer = false
		// only generate client of specific language
		if c.clientLanguage != nil && len(*c.clientLanguage) > 0 {
			if *c.clientLanguage != "go" && *c.clientLanguage != "js" {
				return fmt.Errorf("client-language must be one of \"go\" or \"js\"")
			}
			switch *c.clientLanguage {
			case "go":
				c.generateGoClient = true
				c.generateJSClient = false
			case "js":
				c.generateGoClient = false
				c.generateJSClient = true
			default:
				return fmt.Errorf("client-language must be one of \"go\" or \"js\"")
			}
		} else {
			c.generateGoClient = true
			c.generateJSClient = true
		}
	}
	c.generateGoModels = c.generateServer || c.generateGoClient
	if c.generateJSClient && *c.jsModulePath == "" {
		return fmt.Errorf("js-path is required")
	}

	// determine paths for generated files
	c.modelsPath = filepath.Join(os.Getenv("GOPATH"), "src", c.goPackagePath, "models")
	c.serverPath = filepath.Join(os.Getenv("GOPATH"), "src", c.goPackagePath, "server")
	c.goClientPath = filepath.Join(os.Getenv("GOPATH"), "src", c.goPackagePath, "client")

	return nil
}
