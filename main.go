package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"

	goclient "github.com/Clever/wag/v9/clients/go"
	jsclient "github.com/Clever/wag/v9/clients/js"
	"github.com/Clever/wag/v9/hardcoded"
	"github.com/Clever/wag/v9/models"
	"github.com/Clever/wag/v9/server"
	"github.com/Clever/wag/v9/server/gendb"
	"github.com/Clever/wag/v9/swagger"
	"github.com/Clever/wag/v9/validation"
)

// config contains the configuration of command line flags and configuration derived from command line flags
type config struct {
	clientLanguage     *string
	clientOnly         *bool
	dynamoOnly         *bool
	outputPath         *string
	versionFlag        *bool
	swaggerFile        *string
	relativeDynamoPath *string
	jsModulePath       *string
	goPackageName      *string

	dynamoPath            string
	goAbsolutePackagePath string
	jsClientPath          string

	generateDynamo   bool
	generateGoClient bool
	generateGoModels bool
	generateJSClient bool
	generateServer   bool
	generateTracing  bool
}

var version string

func main() {
	conf := config{
		swaggerFile:        flag.String("file", "swagger.yml", "the spec file to use"),
		goPackageName:      flag.String("go-package", "", "package of the generated go code"),
		outputPath:         flag.String("output-path", "", "relative output path of the generated go code"),
		jsModulePath:       flag.String("js-path", "", "path to put the js client"),
		versionFlag:        flag.Bool("version", false, "print the wag version and exit"),
		clientOnly:         flag.Bool("client-only", false, "only generate client code"),
		clientLanguage:     flag.String("client-language", "", "generate client code in specific language [go|js]"),
		dynamoOnly:         flag.Bool("dynamo-only", false, "only generate dynamo code"),
		relativeDynamoPath: flag.String("dynamo-path", "", "path to generate dynamo code relative to go package path"),
	}
	flag.Parse()
	if *conf.versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	if err := conf.parse(); err != nil {
		log.Fatalf(err.Error())
	}

	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
	doc, err := loads.Spec(*conf.swaggerFile)
	if err != nil {
		log.Fatalf("Error loading swagger file: %s", err)
	}
	swaggerSpec := *doc.Spec()

	if err := validation.Validate(*doc, conf.generateJSClient); err != nil {
		log.Fatalf("Swagger file not valid: %s", err)
	}

	err = swagger.ValidateResponses(swaggerSpec)
	if err != nil {
		log.Fatalf("Failed processing the swagger spec: %s", err)
	}

	if conf.generateGoModels {
		if err := generateGoModels(*conf.goPackageName, conf.goAbsolutePackagePath, swaggerSpec); err != nil {
			log.Fatal(err.Error())
		}
	}

	if conf.generateServer {
		if err := generateServer(*conf.goPackageName, conf.goAbsolutePackagePath, swaggerSpec); err != nil {
			log.Fatal(err.Error())
		}
	}

	if conf.generateTracing {
		if err := generateTracing(conf.goAbsolutePackagePath); err != nil {
			log.Fatal(err.Error())
		}
	}

	if conf.generateDynamo {
		if err := generateDynamo(conf.dynamoPath, *conf.goPackageName, conf.goAbsolutePackagePath, *conf.relativeDynamoPath, swaggerSpec); err != nil {
			log.Fatal(err.Error())
		}
	}

	if conf.generateGoClient {
		if err := generateGoClient(*conf.goPackageName, conf.goAbsolutePackagePath, swaggerSpec); err != nil {
			log.Fatal(err.Error())
		}
	}

	if conf.generateJSClient {
		if err := generateJSClient(*conf.jsModulePath, swaggerSpec); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func generateGoModels(packageName, basePath string, swaggerSpec spec.Swagger) error {
	if err := prepareDir(filepath.Join(basePath, "models")); err != nil {
		return err
	}
	if err := models.Generate(packageName, basePath, swaggerSpec); err != nil {

		return fmt.Errorf("Error generating models: %s", err)
	}
	return nil
}

func generateServer(goPackageName, basePath string, swaggerSpec spec.Swagger) error {
	if err := prepareDir(filepath.Join(basePath, "server")); err != nil {
		return err
	}
	if err := server.Generate(goPackageName, basePath, swaggerSpec); err != nil {
		return fmt.Errorf("Failed to generate server: %s", err)
	}
	middlewareGenerator := swagger.Generator{BasePath: basePath}
	middlewareGenerator.Write(hardcoded.MustAsset("../_hardcoded/middleware.go"))
	if err := middlewareGenerator.WriteFile("server/middleware.go"); err != nil {
		return fmt.Errorf("Failed to copy middleware.go: %s", err)
	}

	return nil
}

func generateTracing(basePath string) error {
	if err := prepareDir(filepath.Join(basePath, "servertracing")); err != nil {

		return err
	}

	tracingGenerator := swagger.Generator{BasePath: basePath}
	tracingGenerator.Write(hardcoded.MustAsset("../_hardcoded/tracing.go"))
	if err := tracingGenerator.WriteFile("servertracing/tracing.go"); err != nil {
		log.Fatalf("Failed to copy tracing.go: %s", err)
	}

	return nil
}

func generateDynamo(dynamoPath, goPackageName, basePath, outputPath string, swaggerSpec spec.Swagger) error {
	if err := prepareDir(dynamoPath); err != nil {
		return err
	}
	if err := gendb.GenerateDB(goPackageName, basePath, &swaggerSpec, outputPath); err != nil {
		return fmt.Errorf("Failed to generate database: %s", err)
	}
	return nil
}

func generateGoClient(goPackageName, basePath string, swaggerSpec spec.Swagger) error {
	if err := prepareDir(filepath.Join(basePath, "client")); err != nil {
		return err
	}
	if err := goclient.Generate(goPackageName, basePath, swaggerSpec); err != nil {
		return fmt.Errorf("Failed generating go client %s", err)
	}
	doerGenerator := swagger.Generator{BasePath: basePath}
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

// parseCmdConfig determines the which code is generated and the location of the generated code
// from the command line arguments
func (c *config) parse() error {
	if err := c.setGoPaths(swag.StringValue(c.outputPath), swag.StringValue(c.goPackageName)); err != nil {
		return err
	}

	clientLanguage, jsModulePath := swag.StringValue(c.clientLanguage), swag.StringValue(c.jsModulePath)
	if err := c.setGenerateFlags(clientLanguage, jsModulePath); err != nil {
		return err
	}

	c.setGeneratedFilePaths()

	return nil
}

func (c *config) setGenerateFlags(clientLanguage, jsModulePath string) error {
	if swag.BoolValue(c.clientOnly) && swag.BoolValue(c.dynamoOnly) {
		return fmt.Errorf("Cannot use -dynamo-only and -client-only together")
	} else if swag.BoolValue(c.clientOnly) {
		if err := c.setClientLanguage(clientLanguage, jsModulePath); err != nil {
			return err
		}
		c.generateGoModels = c.generateGoClient
		c.generateTracing = c.generateGoClient
	} else if swag.BoolValue(c.dynamoOnly) {
		c.generateDynamo = true
		c.generateGoModels = true
		c.generateGoClient = false
	} else {
		c.generateServer = true
		c.generateDynamo = true
		c.generateGoModels = true
		c.generateTracing = true
		if err := c.setClientLanguage(clientLanguage, jsModulePath); err != nil {
			return err
		}
	}
	return nil
}

// setGoPaths sets the golang package path and package name. Go repos using modules may have a
// package name differing from its package path.
func (c *config) setGoPaths(outputPath, goPackageName string) error {
	// check if the repo uses modules
	if modFile, err := os.Open("go.mod"); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return fmt.Errorf("Error checking if the repo contains a go.mod file: %s", err.Error())
		}
		if goPackageName == "" {
			return fmt.Errorf("go-package is required")
		}
		// if the repo does not use modules, the package name is equivalent to the package path
		c.goAbsolutePackagePath = filepath.Join(os.Getenv("GOPATH"), "src", goPackageName)
	} else {
		defer modFile.Close()
		if outputPath == "" {
			return fmt.Errorf("output-path is required")
		}

		absolutePath, err := filepath.Abs(outputPath)
		if err != nil {
			return fmt.Errorf("converting output-path to absolute path: %v", err)
		}
		c.goAbsolutePackagePath = absolutePath

		*c.goPackageName = getModulePackageName(modFile, path.Clean(outputPath))

	}
	return nil
}

// setClientLanguage determines in which langues to generate the server client
func (c *config) setClientLanguage(clientLanguage, jsModulePath string) error {
	if clientLanguage != "" {
		if clientLanguage != "go" && clientLanguage != "js" {
			return fmt.Errorf("client-language must be one of \"go\" or \"js\"")
		}
		switch clientLanguage {
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

	if c.generateJSClient && jsModulePath == "" {
		return fmt.Errorf("js-path is required")
	}

	return nil
}

// setGeneratedFilePaths determines where to output the generated files.
func (c *config) setGeneratedFilePaths() {
	const serverDir = "server"

	if c.generateDynamo {
		// set path of generated dynamo code if none specified
		if swag.StringValue(c.relativeDynamoPath) == "" {
			if c.generateServer {
				c.relativeDynamoPath = swag.String(path.Join(serverDir, "db"))
			} else {
				c.relativeDynamoPath = swag.String("db")
			}
		}
		c.dynamoPath = path.Join(c.goAbsolutePackagePath, *c.relativeDynamoPath)
	}
}

// getModulePackageName gets the package name of the generated code
// Example: if packagePath = github.com/Clever/wag/v8/gen-go and the module name is github.com/Clever/wag/v8/v2
// the function will return github.com/Clever/wag/v8/v2/gen-go
// Example: if packagePath = github.com/Clever/wag/v8/gen-go and the module name is github.com/Clever/wag/v8
// the function will return  github.com/Clever/wag/v8/gen-go
func getSubModulePackageName(modFile *os.File, outputPath string) string {
	r := bufio.NewReader(modFile)
	b, _, err := r.ReadLine()
	if err != nil {
		log.Fatalf("Error checking module name: %s", err.Error())
	}

	// parse module path
	moduleName := strings.TrimPrefix(string(b), "module")
	moduleName = strings.TrimSpace(moduleName)

	//Remove /v<version> from end of path.
	regex, err := regexp.Compile("/v[0-9]$|/v[0-9][0-9]")
	if err != nil {
		log.Fatalf("Error checking module name: %s", err.Error())
	}

	moduleName = regex.ReplaceAllString(moduleName, "")

	return fmt.Sprintf("%v/%v", moduleName, outputPath)

}

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
