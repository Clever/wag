# Subrouters

This document describes the experimental implementation of Gorilla Mux subrouters on this branch.

1. [**Gorilla Mux Subrouters**](#introduction-to-subrouters)
2. [**Configuration**](#configuration)
3. [**Implementation: Root Level**](#root-level)
4. [**Implementation: Subrouter Level**](#subrouter-level)
5. [**Evaluation**](#evaluation)

## Introduction to Subrouters

Most API frameworks across languages have some concept of subrouters that can group route handlers by path prefix and apply common behavior across those grouped routes. This subrouting behavior can be applied for the following benefits, with widely varying degrees of functional improvement:

- Simple grouping of routes for readability and code organization
- Package-level separation of routes such that one route implementation can be installed in multiple places or in different packages
- Separating code ownership of routes among different teams within one server or repo (which may also be facilitated by or complemented by package-level separation)
- Applying middleware to a path prefix rather than to the whole server or to each route individually
- Marginally more efficient route matching for each request in the server at runtime since the API router is essentially matching routes against a tree of paths rather than a list

`gorilla/mux` provides the [`Route.PathPrefix()`](https://pkg.go.dev/github.com/gorilla/mux#Route.PathPrefix) and [`Router.PathPrefix()`](https://pkg.go.dev/github.com/gorilla/mux#Router.PathPrefix) functions for this purpose. The [Matching Routes](https://github.com/gorilla/mux#matching-routes) documentation shows how to use `PathPrefix()` in practice.

## Configuration

To install a subrouter in your wag service, follow the steps in this section.

First, create a `routers/` directory with a subdirectory for your router and a `swagger.yml` file in that subdirectory (replace `KEY` with a meaningful path segment for your router):

```sh
mkdir -p routers/KEY
touch routers/KEY/swagger.yml
```

Fill out your swagger.yml with a meaningful config. The main thing to call out is that `basePath` is required.

```yaml
swagger: '2.0'
info:
  # title is used in the same way as root swagger.yml, but usage may be changed in the future
  title: app-district-service
  description: A router that serves routes related to managing the off-Clever districts themselves, including composing calls to other foundational services like district-config-service and dac-service.

  # version and x-npm-package are functionally ignored in this case, but we have
  # not yet removed them from the validation or generation.
  version: 0.1.0
  x-npm-package: '@clever/app-district-service'
schemes:
  - http
produces:
  - application/json
responses:
  # same as root swagger.yml

# basePath is technically the same as in a traditional swagger.yml file, but it
# is REQUIRED in this case to generate the client. The value must match the path
# configured in the x-routers key in the root swagger.yml.
basePath: /v0/apps
paths:
  # same as root swagger.yml
definitions:
  # same as root swagger.yml
```

Next, add the `x-routers` extension key to the top level of the root `swagger.yml` file. This config enables the root router to discover its subrouters for usage in the server and client code.

```yaml
x-routers:
  - key: districts # matches the subdirectory name under routers/
    path: /v0/apps # matches basePath in routers/districts/swagger.yml
  # ... more subrouters
```

With config for both the root router and subrouter installed, generate both the root router and subrouters in a `go:generate` directive or `make` target, using the `-subrouter` argument for the subrouter invocations.

```go
//go:generate wag -output-path ./gen-go -js-path ./gen-js -file swagger.yml
//go:generate wag -subrouter -output-path ./routers/districts/gen-go -js-path ./routers/districts/gen-js -file ./routers/districts/swagger.yml
//go:generate wag -subrouter [... other subrouter args]
```

Finally, implement the subrouter controllers in the `routers/KEY/controller` packages and wire it up in your server's `main.go`.

```go
    s := server.New(
		    myController,
		    districtscontroller.Controller{},
		    sessionscontroller.Controller{},
		    *addr,
	  )
```

And that's it!

## Implementation: Subrouter Level

### main.go

- `wag` now accepts [a boolean `-subrouter` flag](https://github.com/Clever/wag/blob/subrouters/main.go#L68-L72) to indicate that this target spec is for a subrouter.
- It [loads the parent `swagger.yml` spec](https://github.com/Clever/wag/blob/subrouters/main.go#L90-L99) if it is a subrouter.
- It [uses the parent spec as part of validation](https://github.com/Clever/wag/blob/subrouters/main.go#L104-L106) and [validates that the subrouter `basePath` matches some configured `path` in `x-routers` of the parent spec](https://github.com/Clever/wag/blob/subrouters/validation/validation.go#L262-L283).
- It [passes the value of the `-subrouter` flag to `generateServer`](https://github.com/Clever/wag/blob/subrouters/main.go#L124), which [passes it to `server.Generate()`](https://github.com/Clever/wag/blob/subrouters/main.go#L174) and [skips middleware generation for subrouters](https://github.com/Clever/wag/blob/subrouters/main.go#L178-L184).

The `generateClient` call is unchanged because it uses the pre-existing notion of `basePath`.

### Server

- `server.Generate()` [passes its new `subrouter` boolean arg to `generateRouter()`](https://github.com/Clever/wag/blob/subrouters/server/genserver.go#L18).
- `generateRouter` [sets the new `routerTemplate.IsSubrouter` struct field to the value of the `subrouter` arg](https://github.com/Clever/wag/blob/subrouters/server/genserver.go#L53). Note that [`routerTemplate.IsSubrouter`](https://github.com/Clever/wag/blob/subrouters/server/genserver.go#L41) is different from `routerTemplate.Subrouters`, which is defined for routers that _have_ subrouters rather than a subrouter itself; see [Server](#server-1) under [Implementation: Root Level](#implementation-root-level).
- _DOES NOT_ [prepend the `basePath` from the spec to the path](https://github.com/Clever/wag/blob/subrouters/server/genserver.go#L59-L62) in generating operations for the router to handle.
- It [creates a limited slice of imports for the subrouter `router.go`](https://github.com/Clever/wag/blob/subrouters/server/genserver.go#L84-L92).
- It cuts out sections irrelevant to subrouters in the router template with an `{{if not .IsSubrouter}}` so that only the `handler` struct and `Register()` function remain.
  - [`Server` struct, `serverConfig` struct, `CompressionLevel` function, `Server.Serve` function](https://github.com/Clever/wag/blob/subrouters/server/router.go#L10-L84)
  - [`startLoggingProcessMetrics` function, `withMiddleware` function, `New` function, `NewRouter` creator function, `newRouter` function](https://github.com/Clever/wag/blob/subrouters/server/router.go#L89-L161)
  - [`NewWithMiddleware` function, `AttachMiddleware` function](https://github.com/Clever/wag/blob/subrouters/server/router.go#L191-L235)

The `buildSubrouters()` call in `generateRouter()` only applies to parent routers, not subrouters.

### Client

The client generation for subrouters is essentially unchanged from the default behavior. Since we've already validated the `basePath` against the matching `x-routers[*].path` value in the parent spec, we just rely on the existing behavior to prepend the `basePath` in the client. On the server side, we don't do that, because the `PathPrefix()` call handles the `basePath` from the side of `x-routers[*].path`. (Alternatively, we could maybe have the parent generation discover its subrouter paths from the subrouter specs and omit the `path` entirely; this is a future consideration.)

## Implementation: Root Level

### `main.go`

### Server

- Gets [`template.Subrouters` and a slice of subrouter `gen-go/server` package imports](https://github.com/Clever/wag/blob/subrouters/server/genserver.go#L78) with `buildSubrouters()`. `buildSubrouters()` [calls `swagger.ParseSubrouters()`](https://github.com/Clever/wag/blob/subrouters/server/genserver.go#L131) and [maps the subrouter keys to their server package imports](https://github.com/Clever/wag/blob/subrouters/server/genserver.go#L136-L147). This setup allows the `server` package to discover subrouters from the root `swagger.yml` file.
- It [uses the full list of imports to run the whole server](https://github.com/Clever/wag/blob/subrouters/server/genserver.go#L92-L117), including [importing `gen-go/server` packages for its subrouters](https://github.com/Clever/wag/blob/subrouters/server/genserver.go#L109). Mostly these functions are wrapping each other, which is why this pattern repeats for so many functions.
- Accepts additional controllers for subrouters in the `New`, `NewWithMiddleware`, `NewRouter`, `newRouter`, and `Register` functions in the `router.go` template.
  - [`New` with subrouter controllers](https://github.com/Clever/wag/blob/subrouters/server/router.go#L113-L127)
  - [`NewWithMiddleware` with subrouter controllers](https://github.com/Clever/wag/blob/subrouters/server/router.go#L192-L207)
  - [`NewRouter` with subrouter controllers](https://github.com/Clever/wag/blob/subrouters/server/router.go#L129-L142)
  - [`newRouter` with subrouter controllers](https://github.com/Clever/wag/blob/subrouters/server/router.go#L144-L159)
  - [`Register` with subrouter controllers](https://github.com/Clever/wag/blob/subrouters/server/router.go#L162-L189)
- Most importantly, the parent `Register()` [includes calls to subrouter `Register()` functions](https://github.com/Clever/wag/blob/subrouters/server/router.go#L182) with the actual Gorilla Mux subrouter created by `PathPrefix()`:

```go
KEYrouter.Register(router.PathPrefix("PATH").Subrouter(), subcontroller)
```

### Client

- All of `generateClient`, `generateInterface`, `generateClientInterface`, and `CreateModFile` use `swagger.ParseSubrouters()` to extract the `x-routers` extension config from the root `swagger.yml` file.
  - [`generateClient` call to `swagger.ParseSubrouters()`](https://github.com/Clever/wag/blob/subrouters/clients/go/gengo.go#L171)
  - [`generateInterface` call to `swagger.ParseSubrouters()`](https://github.com/Clever/wag/blob/subrouters/clients/go/gengo.go#L331)
  - [`generateClientInterface` call to `swagger.ParseSubrouters()`](https://github.com/Clever/wag/blob/subrouters/clients/go/gengo.go#L374)
  - [`CreateModFile` call to `swagger.ParseSubrouters()`](https://github.com/Clever/wag/blob/subrouters/clients/go/gengo.go#L247)
- `CreateModFile` [adds `replace` directives for subrouter `gen-go/client` and `gen-go/models` packages](https://github.com/Clever/wag/blob/subrouters/clients/go/gengo.go#L273-L294). It _should_, in the future, add this for the subrouter `gen-go/server` package as well.
- `generateInterface` [adds imports for the subrouter to `interface.go`](https://github.com/Clever/wag/blob/subrouters/clients/go/gengo.go#L337-L349).
- `generateClientInterface`, which is called by `generateInterface`, embeds the subrouter `Client` interfaces in the parent `Client` interface.
- `generateClient` [creates handler code for subrouter operations that wraps the subrouter client methods](https://github.com/Clever/wag/blob/subrouters/clients/go/gengo.go#L203-L224) and [instantiates subrouter clients with their `New()` functions within the parent `New()` function](https://github.com/Clever/wag/blob/subrouters/clients/go/gengo.go#L120-L122).

## Evaluation

Given that we (Maddy and the API team) have run this experiment to this point, what is its status? What is the value in this potential feature? How reliable and easy to work with are the design and implementation?

At this point, the feature is not and will not be supported by Infra. As such, this branch can be treated like a fork of an open source, where the API team and any other teams that use it have to pay [the maintenance costs](#maintenance) of updating it against the main branch.

### User Stories and Potential Benefits

First off, it's important to consider why product teams would actually want to consider adopting this subrouter fork, what potential value it provides.

Let's start by separating the benefits that can be achieved by other methods from the benefits which are less likely to be achieved without subrouters.

Examples of user stories that can be achieved without a subrouter implementation:

- When working on a service that is very large, with potentially tens of routes, I want to group routes within one service into multiple specs and controller packages so that the service can be split into more manageable, readable pieces.
- When working on a service that has routes with very divergent dependencies, I want to split up those route implementations across multiple controllers so that each route handler has access only to the dependencies it needs, or at least eliminates access to dependencies it doesn't need, to the extent that it's possible to enforce this grouping by path prefix.
- When working on a service that has multiple teams working in it (possibly because it is very large or has routes with divergent dependencies), I want to split up route specs and controller implementations by directory path so that I can use Github `CODEOWNERS` or other tools to define team ownership by directory path.
- When working with a service that has routes versioned by path prefix, I want to separate the versions into separate packages so I can more easily comprehend which version has which behavior and not generate models with `V2`, `V3`, etc suffixes.

Examples of user stories that can only be achieved with a subrouter implementation or are less likely to be achieved without one:

- When working on a service that is very large, with potentially tens of routes, I want to optimize route matching so that the Gorilla Mux router matches paths by walking a tree of prefixes until it reaches the leaves/terminal segments of the path rather than matching against a list of literal routes.
- When working on a service where routes have shared behavior by path prefix — for example, ensuring that a path segment of the type `/collection/{itemID}` has a valid object reference of `itemID` within `collection`, and perhaps that it meets other domain-specific constraints — I want to implement that behavior as middleware rather than handler by handler. This functionality is not supported by the current implementation, but it could easily be; however, the specific example given would require additional work to use spec extensions or support a non-compliant `basePath`, because path parameters are not supported in OpenAPI Spec's `basePath` regardless of the OAS version.

Broadly speaking, the balance of the benefits could be achieved through some other means. Those means could include

- Concatenation of multiple OpenAPI specs
- Usage of OAS v3 which allows referencing other specs (if I'm not mistaken)
- Embedding subcontrollers from subpackages in a root controller manually

In my opinion (Maddy), it would be ideal for the platform to provide explicit support for any such features as required. From the list above, only concatenation of multiple OpenAPI specs would require platform support on its own, as OpenAPI Spec v3 support is a feature request that is going to be prioritized in its own right and embedding subcontrollers does not require any platform support. The platform could also choose to support mapping middleware to specific routes in order to apply behavior across multiple routes but not the whole server through some means other than subrouters, but that isn't a trivial feature either, and I think it's likely preferable to use subrouters in order to get the route matching optimization as well.

All in all: you decide! If you happen to try this out, please register your feedback wherever appropriate (backend guild, the API team, the Infra team).

### Installation

To install the subrouters fork of `wag`, run the following command in your Go module root.

```sh
go get -u github.com/Clever/wag/v9@subrouters
```

You'll need to rerun this command any time that the `subrouters` branch is rebased off the `wag` main branch to get the latest changes in the platform

I recommend also configuring `wag` as a tool and using `go install tool` to update the version of `wag` in your path.

### Maintenance

In order to get upstream changes, we'll need to regularly rebase off the main branch for `wag`. Consumers of the subrouters experiment should reinstall via `go get -u`. We should schedule that maintenance to the degree possible.

### Implementation

It's worth considering: How much complexity does the subrouter implementation add to the `wag` implementation, and would that complexity be hard to maintain going forward?

My assessment is that yes, it's a little unnecessarily complex, but that that complexity is also a function of `wag`'s heavily procedural implementation. If we were to actually adopt this feature, I would suggest potentially suggest refactoring `wag` to be more object-oriented and interface-driven so that it can support multiple types of targets and implementations can be resolved and injected based on `wag` args and `swagger.yml` config. That kind of refactor could also support other types of features, like OAS v2 and v3 support simultaneously (although that's a little different in that it's probably mostly about mapping different inputs to the same output, the principle applies).
