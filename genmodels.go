package main

import (
	"fmt"

	"github.com/go-openapi/spec"

	"github.com/go-swagger/go-swagger/generator"
)

func generateModels(packageName, swaggerFile string, swagger spec.Swagger) error {

	// generate models with go-swagger
	if err := generator.GenerateServer("", []string{}, []string{}, generator.GenOpts{
		Spec:           swaggerFile,
		ModelPackage:   "models",
		Target:         "./generated/",
		IncludeModel:   true,
		IncludeHandler: false,
		IncludeSupport: false,
	}); err != nil {
		return err
	}

	if err := generateOutputs(packageName, swagger.Paths); err != nil {
		return err
	}
	return generateInputs(swagger.Paths)
}

func generateInputs(paths *spec.Paths) error {

	var g Generator

	g.Printf(`
package models

import(
		"encoding/json"
		"strconv"
)

// These imports may not be used depending on the input parameters
var _ = json.Marshal
var _ = strconv.FormatInt
`)

	for _, pathKey := range sortedPathItemKeys(paths.Paths) {
		path := paths.Paths[pathKey]
		pathItemOps := pathItemOperations(path)
		for _, opKey := range sortedOperationsKeys(pathItemOps) {
			op := pathItemOps[opKey]
			if err := printInputStruct(&g, op); err != nil {
				return err
			}
			if err := printInputValidation(&g, op); err != nil {
				return err
			}
		}
	}

	return g.WriteFile("generated/models/inputs.go")
}

func printInputStruct(g *Generator, op *spec.Operation) error {
	g.Printf("type %sInput struct {\n", capitalize(op.ID))

	for _, param := range op.Parameters {
		if param.In == "formData" {
			return fmt.Errorf("Input parameters with 'In' formData are not supported")
		}

		var typeName string
		if param.In != "body" {
			switch param.Type {
			case "string":
				typeName = "string"
			case "integer":
				typeName = "int64"
			case "boolean":
				typeName = "bool"
			case "number":
				typeName = "float64"
			default:
				// Note. We don't support 'array' or 'file' types even though they're in the
				// Swagger spec.
				return fmt.Errorf("Unsupported param type")
			}
		} else {
			var err error
			typeName, err = typeFromSchema(param.Schema)
			if err != nil {
				return err
			}
		}
		g.Printf("\t%s %s\n", capitalize(param.Name), typeName)
	}
	g.Printf("}\n\n")

	return nil
}

func printInputValidation(g *Generator, op *spec.Operation) error {
	g.Printf("func (i %sInput) Validate() error{\n", capitalize(op.ID))

	// TODO: Right now we only support validation on complex types (schemas)
	for _, param := range op.Parameters {
		if param.In == "body" {
			g.Printf("\tif err := i.%s.Validate(nil); err != nil {\n", capitalize(param.Name))
			g.Printf("\t\treturn err\n")
			g.Printf("\t}\n\n")
		}
	}
	g.Printf("\treturn nil\n")
	g.Printf("}\n\n")

	return nil
}

func generateOutputs(packageName string, paths *spec.Paths) error {
	var g Generator

	g.Printf("package models\n\n")

	g.Printf(defaultOutputTypes())

	for _, pathKey := range sortedPathItemKeys(paths.Paths) {
		path := paths.Paths[pathKey]
		pathItemOps := pathItemOperations(path)
		for _, opKey := range sortedOperationsKeys(pathItemOps) {
			op := pathItemOps[opKey]
			// We classify response keys into three types:
			// 1. 200-399 - these are "success" responses and implement the Output interface
			// 	defined above
			// 2. 400-599 - these are "failure" responses and implement the error interface
			// 3. Default - this is defined as a 500 (TODO: decide if we want to keep this...)

			// Define the success interface
			g.Printf("type %sOutput interface {\n", capitalize(op.ID))
			g.Printf("\t%sStatus() int\n", capitalize(op.ID))
			g.Printf("\t// Data is the underlying model object. We know it is json serializable\n")
			g.Printf("\t%sData() interface{}\n", capitalize(op.ID))
			g.Printf("}\n\n")

			// Define the error interface
			g.Printf("type %sError interface {\n", capitalize(op.ID))
			g.Printf("\terror // Extend the error interface\n")
			g.Printf("\t%sStatusCode() int\n", capitalize(op.ID))
			g.Printf("}\n\n")

			for _, statusCode := range sortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
				response := op.Responses.StatusCodeResponses[statusCode]

				if statusCode < 200 || statusCode > 599 {
					// TODO: Write a test for this...
					return fmt.Errorf("Response map key must be an integer between 200 and 599 or "+
						"the string 'default'. Was %d", statusCode)
				}
				if statusCode == 400 {
					return fmt.Errorf("Use the pre-defined default 400 response 'DefaultBadRequest' " +
						"instead of defining your own")
				} else if statusCode == 500 {
					return fmt.Errorf("Use the pre-defined default 500 response `DefaultInternalError` " +
						"instead of defining your own")
				}

				outputName := fmt.Sprintf("%s%dOutput", capitalize(op.ID), statusCode)
				typeName, err := typeFromSchema(response.Schema)
				if err != nil {
					return err
				}

				g.Printf("type %s %s\n\n", outputName, typeName)

				g.Printf("func (o %s) %sData() interface{} {\n", outputName, capitalize(op.ID))
				g.Printf("\treturn o\n")
				g.Printf("}\n\n")

				if statusCode < 400 {

					// TODO: Do we really want to have that as part of the interface?
					g.Printf("func (o %s) %sStatus() int {\n", outputName, capitalize(op.ID))
					// TODO: Use the right status code...
					g.Printf("\treturn %d\n", statusCode)
					g.Printf("}\n\n")

				} else {

					g.Printf("func (o %s) Error() string {\n", outputName)
					// TODO: Would it make sense to give this a constructor so we can have a more detailed
					// error message?
					g.Printf("\treturn \"Status Code: \" + \"%d\"\n", statusCode)
					g.Printf("}\n\n")

					g.Printf("func (o %s) %sStatusCode() int {\n", outputName, capitalize(op.ID))
					g.Printf("\treturn %d\n", statusCode)
					g.Printf("}\n\n")
				}
			}
		}
	}
	return g.WriteFile("generated/models/outputs.go")
}
