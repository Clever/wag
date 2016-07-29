# wag
sWAGger - Web API Generator

# Running
```
go run main.go
cp hardcoded/main.go generated/
cp hardcoded/middleware.go generated/
cd generated
go run main.go middleware.go router.go handlers.go contexts.go controller.go outputs.go
```

## After Generating

### Main
Two files in the generated directory are hard-coded and independent of your .yml file
- hardcoded/main.go
- hardcoded/middleware.go

The rest are built based on your swagger yml file definition.

### Files to Change
After autogenerating the code you should modify two files:
# TODO Don't copy over the controller file if it already exists
- In controller.go implement the logic of your handlers.
- In middleware.go add any middleware specific to your service


## Swagger Spec

Currently, this repo doesn't implement the entire Swagger Spec. This is a non-exhaustive list of the parts of the Swagger specification that haven't been implemented yet.

### Planning to Implement
All Swagger Data Types (long, float, double, etc...)
Schema
  - parameters
  - tags
$ref in multiple places
Operation
  - tags
Required Fields
Response
  - Headers

### Not Planning to Implement (at least for now)
All Mime Types
Patterned Fields (these are vendor specific extensions anyway)
Multi-File Swagger Definitions
Everything in JSON-Schema (getting a lot of this from the auto-generated go-swagger code)
Schema
  - host
  - basePath
  - scheme (just http / https)
  - consumes
  - produces
  - securityDefinitions
  - security
Consumes
  - produces
  - consumes
  - schemes
  - security
Form parameter type (though if it's easy maybe we should just add...)
Parameter
  - items
  - collectionFormat
  - all the json schema requirements? (uniqueItems, multipleOf, etc...)
Schema object (for now going to try to get this from somewhere else...)
Discriminators
XML Modeling
Security Objects

