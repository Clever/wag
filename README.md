# wag
sWAGger - Web API Generator

## Usage
Wag requires Go 1.7+.
### Generating Code
Create a swagger.yml file with your [service definition](http://editor.swagger.io/#/). Wag supports a [subset](https://github.com/Clever/wag#swagger-spec) of the Swagger spec.
Copy the latest `wag.mk` from the [dev-handbook](https://github.com/Clever/dev-handbook/blob/master/make/wag.mk).
Set up a `generate` target in your `Makefile` that will generate server and client code:

```
include wag.mk

WAG_VERSION := latest

generate: wag-generate-deps
	$(call wag-generate,./swagger.yml, $(PKG))
```

Define global BadRequest and InternalError response types. These are used internally for validation errors and unknown errors respectively. They must reference a definition with a message field. For example:

```
responses:
  BadRequest:
    description: Bad Request
    schema:
      $ref: "#/definitions/BadRequest"
  InternalError:
    description: Internal Error
    schema:
      $ref: "#/definitions/InternalError"

definitions:
  BadRequest:
    type: object
    properties:
      message:
        type: string
  InternalError:
    type: object
    properties:
      message:
        type: string
```

For more information on error definitions see the Errors section below.

Then generate your code:
```
make generate
```

This generates four directories. You should not have to modify any of the generated code:
- gen-go/models: contains all the definitions in your Swagger file as well as the API input / output definitions
- gen-go/server: contains the router, middleware, and handler logic
- gen-go/client: contains the Go client library
- gen-js: contains the javascript client library

## Implementing and Running the Server
To implement and run the generated server you need to:
- Implement the controller interface defined in `gen-go/server/interface.go`
- Pass the controller into the Server constructor. For example:
```
  s := server.New(myController, ":8000")
  // Serve should not return
  log.Fatal(s.Serve())
```
Or, with custom middleware:
```
  s := server.NewWithMiddleware(myController, ":8000", []func(http.Handler) http.Handler{
    myFirstMiddlware, mySecondMiddlware})
  // Serve should not return
  log.Fatal(s.Serve())
```

### Interface

The server interface defined in `gen-go/server/interface.go` has one method for each operation defined in the swagger.yml. We generate the interface based on the following rules:

#### Input Parameters
  * The first argument to each Wag operation is a `context.Context`. See below for more details on how Wag uses contexts.
  * If one parameter is defined then Wag uses that input directly in the function definition.
  `func F(ctx context.Context, input string) error`
  * If more than one parameter is defined then Wag generates a input struct with all the parameters:
  `func F(ctx context.Context, input *models.{{OperationID}}Input) error`
    * Optional parameters that don't have defaults are pointers in the input struct so that the server can distinguish between parameters that aren't set and parameters that are set to the zero value.


#### Response Parameters
  * Wag only supports defining a single 2XX response status code and doesn't support 3XX status codes.
    * If the success response type has a data type associated with it then Wag generates an interface that takes a pointer to that data type as the first argument.
    `func(...) (*SuccessType, error)`
    * If the operation uses `x-paging`, then Wag generates the an interface that takes the type of the page ID parameter as the second return type.
    `func(...) (*SuccessType, PageParamType, error)`
    * If the success response type doesn't define a data type then Wag generates an interface with only an error response. A nil error tells the client that the request succeeded.
    `func(...) error`


### Logging
  The [kayvee middleware logger](https://godoc.org/gopkg.in/Clever/kayvee-go.v6/middleware) is automatically added to the context object.
  It can be pulled out of the context object and used via the kayvee `FromContext` method:

```go
import "gopkg.in/Clever/kayvee-go.v6/logger"
...
logger.FromContext(ctx).Info(...)
```

  You should use this logger for all logging within your controller implementation.

#### Application Log Routing

**Note**: This is an internal Clever feature. Ignore this if you are using `wag` outside of Clever (also, hi!)

`wag` is already set up for log routing if you so wish. To set up [application log routing](https://clever.atlassian.net/wiki/display/ENG/Application+Log+Routing) in your service, add your `kvconfig.yml` file to the same directory as your service executable. e.g. in your Dockerfile:

```
COPY kvconfig.yml /usr/bin/kvconfig.yml
COPY bin/my-wag-service /usr/bin/my-wag-service
```

### Errors
  * Wag supports three types of errors
    * Global error response types
    * Response types for a specific operation
    * Unexpected errors

  * Any of these can be returned from a controller. To return a global or response specific error type return a pointer to the model defintion for that error type. To return an unexpected error return any Go error. Wag automatically converts errors not defined swagger yml into the default 500 response.

  * All error responses defined in the swagger yml must have a `Message` field. The field is used as the return value of the `Error()` for the corresponding Go error type.

  * Wag has two built-in errors: `#/definitions/BadRequest` (400) and '#/responses/InternalError' (500). Any operation that doesn't explicitly define a 400 and/or 500 response gets these automatically so Wag can use them to return validation and internal errors respectively.

  * Errors returned from your controller are logged by the
  autogenerated handler code, so there is no need to separately log errors
  yourself. If you use the `github.com/go-errors/errors` package, the
  stacktrace will also be logged, making debugging easier.

  For undefined error types the best practices are:
    * If you receive an error from an external dependency, use
      `errors.WrapPrefix(err, "foopackage.func", 0)` to return an error with a
     stacktrace and prefix the source of the error (in this example, we received
     an error from the function `func` in the package `foopackage`).
   * If you generate a new error, use `errors.Errorf` or `errors.New` to build
     the error.
   * If you receive an error from an internal function, just return the error
    directly since it should already have stacktrace information (either it is
      a wrapped external error or a `go-errors`-generated internal error).

### Input Parameters
  * Wag supports four types of parameters
    * Path parameters
      * Must be required
      * Must be a simple type (e.g. string, integer)
      * Will not be pointers
    * Body parameters
      * Must reference a 'definition' schema
      * Will be pointers
      * Cannot have defaults
    * Query parameters
      * Must be a simple or array type. If an array must be an array of strings
      * If the type is 'simple' and the parameter is not required, the type will be a pointer
      * If the type is 'array' it won't be a pointer. Query parameters can't distinguish between an empty array and an nil array so it converts both these cases to a nil array. If you need to distinguish between the two use a body parameter
      * In other cases the parameter is not a pointer
    * Header parameters
      * Must be simple types
      * If marked required will ensure that the input isn't the nil value. Headers cannot have pointer types since HTTP doesn't distinguish between empty and missing headers.
      * If it doesn't have a default value specified, the default value will be the nil value for the type

### Paging
  * Wag can help implement paging on endpoints if you use the `x-paging`
    configuration:
    ```
    /books:
      operationId: getBooks
      x-paging:
        pageParameter: startingAfter
    ```
  * `pageParameter` should be set to the name of a parameter you define on the
    operation that specifies a page ID. The parameter can be a query, header,
    or path parameter (it may not be the only path parameter).
  * You can also specify `resourcePath` if the array you want iterate over
    isn't the top-level success return for the operation. (See `/authors` in
    samples/swagger.yml.)
  * The server interface will be amended with a 2nd success return type so that
    your controller can provide the next page ID the client should fetch.
  * The autogenerated Go client will include a `New<OperationID>Iter` function
    that returns an iterator object that exposes a Next() function that will
    return successive resources, requesting new pages as needed.
  * The autogenerated JS client will include an `<operationID>Iter` function
    that exposes `map`, `forEach`, and `toArray` functions to iterate over the
    results, again requesting new pages as needed.

### Contexts
  * The first argument to every Wag function is a `context.Context` (https://blog.golang.org/context). Contexts play a few important roles in Wag.
    * They can be used to set request specific behavior like a retry policy in client libraries. This includes timeouts and cancellation.
    * They are used to pass metadata, like tracing information, across API requests.

  * To get these benefits pass the context object to any subsequent network requests you make.
    Many client libraries accept the context object, e.g.:
      * **net/http**: If you're making HTTP requests, use the [golang.org/x/net/context/ctxhttp](https://godoc.org/golang.org/x/net/context/ctxhttp) package.
      * **wag** If your handler consumes a `wag`-generated client, then pass the context object to these client methods.

  * If you don't have a context to pass to a Wag function you have two options
    * context.Background() - use this when this is the creator of the request chain, like a test or a top-level service.
    * context.TODO() - use this when you haven't been passed a context from a caller yet, but you expect the caller to send you one at some point.


### DynamoDB Codegen

  * Wag can auto-generate server code to save models to DynamoDB if you specify the `x-db` extension on a schema:
  ```yaml
definitions:
  Thing:
    x-db:
      AllowOverwrites: false
      DynamoDB:
        KeySchema:
          - AttributeName: name
            KeyType: HASH
          - AttributeName: version
            KeyType: RANGE
    type: object
    properties:
      name:
        type: string
      version:
        type: integer
  ```
  The above will generate a `db` package with code to load/persist `Thing` objects from/to DynamoDB.
  * `AllowOverwrites` specifies whether the auto-generated `Save` method should succeed if the object already exists.
  * `DynamoDB` specifies the configuration for a DyanmoDB table for the schema.
     It follows the format of the [`AWS::DynamoDB::Table`](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-dynamodb-table.html) CloudFormation resource.
     Currently it supports a subset of the configuration allowed there.


### Tracing
  `wag` instruments the context object with tracing-related metadata.
  This is done via [opentracing](http://opentracing.io/).
  In order for it to work, you are required to do two things:

  * Configure your `main()` function to report tracing data to [LightStep](http://lightstep.com/).
     We are testing Lightstep as a way to view tracing-related data:
```go
package main

import (
	lightstep "github.com/lightstep/lightstep-tracer-go"
	opentracing "github.com/opentracing/opentracing-go"
)

func main() {
	tags := make(map[string]interface{})
	tags[lightstep.ComponentNameKey] = "<name of the repo>"
	lightstepTracer := lightstep.NewTracer(lightstep.Options{
	    AccessToken: os.Getenv("LIGHTSTEP_ACCESS_TOKEN"),
	    Tags:        tags,
	})
	defer lightstep.FlushLightStepTracer(lightstepTracer)
	opentracing.InitGlobalTracer(lightstepTracer)
	...
}
```
  * By default we tag traces with `http.method`, `span.kind`, `http.url`, `http.status_code`, and `error`. For more information about what these tags mean see: https://github.com/opentracing/opentracing.io/blob/95b966bd6a6b2cf0f231260e3e1fa6206ede2151/_docs/pages/api/data-conventions.md#component-identification

## Using the Go Client
Initialize the client with `New`
```
c := client.New("https://url_of_your_service:port")
```

Make an API call
```
books, err := c.GetBooks(ctx, GetBookByIDInput{Authors: []string{"Twain"}})
if err != nil {
  // Do something with the error
}
```

If you're using the client from another WAG-ified service you should pass in the `ctx` object you get in your server handler. Otherwise you can use `context.Background()`

### Custom String Validation
We've added custom string validation for mongo-ids to avoid repeating: "^[0-9a-f]{24}$"` throughout the swagger.yml. To use it you have must:

- Change you swagger.yml file to have the `mongo-id` format. For example:
```
authorID:
        type: string
        format: mongo-id
```

- Import `github.com/Clever/wag/swagger` and call `swagger.InitCustomFormats()` in your server code.

Note that custom string validation only applies to input parameters and does not have any impact on objects defined in '#/definitions'.

Right now we do not allow user-defined custom strings, but this is something we may add if there's sufficient demand.


## Using the Javascript Client
You can initialize the client by either passing a url or by using [discovery](https://github.com/Clever/discovery-node).

```javascript
import * as SampleClientLib from '@clever/sample-client-lib-js';

const sampleClient = new SampleClientLib({address: "https://url_of_your_service:port"}); // Explicit url
// OR
const sampleClient = new SampleClientLib({discovery: true}); // Using discovery
```

You may also configure a global timeout for requests when initalizing the client.

```javascript
const sampleClient = new SampleClientLib({discovery: true, timeout: 1000}); // Timeout any requests taking longer than 1 second
```

You may then call methods on the client. Methods support callbacks and promises.

```javascript
// Promises
sampleClient.getBookById("bookID").then((book) => {
  // ...
}).catch((err) => {
  // ...
});

// Callbacks
sampleClient.getBookById("bookID", (err, book) => {
  // ...
});
```

You can also pass an optional options argument. This can have the following options
- `timeout` - overide the global timeout for this specific call
- `span` - Pass an opentracing span to instrument with the call - More on this below

```javascript
const options = {
  timeout: 5000 // Timeout after 5 seconds
}

sampleClient.getBookById("bookID", options, (err, book) => {
  // ...
});
```

#### Tracing

To utilize the `span` option above you need to pass an opentracing span into the request. The below
example shows you how to setup an express app to track requests and any calls made via a wag client.

```bash
npm install --save tracing-middleware
```


```javascript
import * as express from "express";
import middleware from "tracing-middleware";
import * as SampleClientLib from 'sample-client-lib-js';

const app = express();
const sampleClient = new SampleClientLib({discovery: true}); // Using discovery

const LIGHTSTEP_ACCESS_TOKEN = "access_token";
app.use(middleware({access_token: LIGHTSTEP_ACCESS_TOKEN}));

app.get("/my-url", (req, res) => {
  sampleClient.getBookById("bookID", {span: req.span}, (err, book) => {
    if (book.isSpecial) {
    	// guarantee that this trace instance will be sampled
    	req.span.setTag("sampling.priority", 1);
    }
  });
});
```

## Tests
```
make test
```

## Swagger Spec

Currently, Wag doesn't implement the entire Swagger Spec. A couple things to keep in mind:
- All schemas should reference type definitions in /definitions. Any schemas defined in /paths will cause an error.
- Scheme, produces, and consumers can only be defined in the top-level swagger object, not individual operations. On the top level object the scheme must be 'http', produces must be 'application/json' and consumes must be 'application/json'

Below is a more comprehensive list of the features we don't yet support

### Unsupported Features
Mime Types

Multi-File Swagger Definitions

Schema:
- host
- tags
- scheme (must be http)
- consumes
- produces
- securityDefinitions
- security

Consumes:
- produces (must be application/json)
- consumes (must be application/json)
- schemes
- security

Form parameter type
Parameter:
- file parameter type
- collectionFormat
- global definitions
- possibly the json schema requirements? (uniqueItems, multipleOf, etc...)

Schema object (all these have to be defined in /definitions and are generated by go-swagger)

Discriminators

XML Modeling

Security Objects

Response:
  - Headers

## Development

The following directories and files are generated and should not be manually edited:
- samples/gen-&ast;/&ast;
- hardcoded/hardcoded.go

Once you've made your changes, run `make test` and check that the generated code looks ok, then check in those changes.
