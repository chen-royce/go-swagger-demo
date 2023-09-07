## Intro
Many of us use Go for our backend web services. Today, we are going to be talking about [go-swagger](https://github.com/go-swagger/go-swagger), a tool for auto-generating OpenAPI-compatible YAML and JSON and Swagger docs from Go code.

### First, some context: What are OpenAPI and Swagger?
At a very high level, tools for documenting HTTP APIs in a structured way that's standardized and understandable by both humans and machines.

But let's get into the specifics:

#### [OpenAPI Specification](https://www.openapis.org/) (OAS)
* Specification offering a structured, defined format for describing REST APIs
* Can be written in either YAML or JSON formats
* Can describe things like available endpoints (and their supported parameters and HTTP methods), higher-level details like API licensing and terms, and more
* Designed to be both human- and machine-readable
* Can be used to generate code and client libraries, documentation, tests, ...

#### [Swagger](https://swagger.io/)
* A suite of open-source tools built around OAS to help with API design, documentation, and consumption
* Swagger Editor ([old UI](https://editor.swagger.io/), [new UI](https://editor-next.swagger.io/)) is a browser-based editor to assist you in writing these JSON/YAML files and visualizing their outputs in real-time
* [Swagger UI](https://github.com/swagger-api/swagger-ui) turns OAS definitions into interactive documentation
* [Swagger Codegen](https://github.com/swagger-api/swagger-codegen) generates code (API clients, server stubs) from OAS definitions

### Swagger's pain points
go-swagger attempts to solve for a couple of pain points pertaining to Swagger documentation:
1. Writing YAML/JSON by hand is really hard.
2. Keeping documentation up-to-date is really hard, and the further it lives from your day-to-day workflow (your code, your VSCode, etc.), the harder it is.

go-swagger proposes solutions to both of the above issues:
1. Doesn't require any YAML or JSON to be written by hand
2. Keeps documentation living right inside your IDE by generating YAML or JSON from your application code, so you don't need to download or visit any external tools like Confluence or the Swagger Editor

## Using go-swagger
We're going to specifically focus on using go-swagger to generate Swagger documentation, although it's capable of [much more](https://github.com/go-swagger/go-swagger#features) (e.g. generating fake servers and clients, CLI tools, and more).

The [complete go-swagger docs are here](https://goswagger.io/).

### Installation
#### :beer: Homebrew
```bash
brew tap go-swagger/go-swagger
brew install go-swagger
```

#### :wrench: Go install
:warning: Prerequisites: A properly configured `$PATH` for running Go binaries globally. For more help, see [this Stack Overflow thread](https://stackoverflow.com/questions/70832925/how-to-properly-install-go-with-paths-and-all).
```bash
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```

## Starting to document your API
The following is a minimalistic spec designed just to showcase how go-swagger works. It's not an exhaustive walkthrough of everything that go-swagger can do and describes a toy problem API, but it is enough to get a basic API doc up and running and to show the core features.

For a deeper dive into what go-swagger can do, again refer to the [official go-swagger docs](https://goswagger.io/).

It's recommended that you go through some of the [initial setup for generating and deploying your documentation](#generation-and-deployment) as you experiment here in order to see your changes in real-time, even if they're local changes.

### API metadata
The [swagger:meta](https://goswagger.io/use/spec/meta.html) annotation provides metadata about our API.

**Example:**
```go
// Package main Echo API documentation
//
// The Echo API repeats and formats via the query parameters
// passed into it
//
// Terms Of Service:
// there are no TOS at this moment, use at your own risk
//
// Schemes: http
// BasePath: /api
// Version: 0.0.1
// Contact: Royce Chen<rchen@rvohealth.com>
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package classification
```

### Documenting routes
The [swagger:route](https://goswagger.io/use/spec/route.html) annotation allows you to document your API routes.

#### Syntax
```
swagger:route [method] [path pattern] [?tag1 tag2 tag3] [operation id]
```

##### Example
```
swagger:route GET /echo echoText
```

#### Properties
- Consumes/Produces (e.g. application/json)
- Schemes (e.g. HTTP, HTTPS)
- Deprecated
- Security (e.g. API keys)
- Parameters (supports custom structs)
- Responses (supports custom structs)
 
##### Example
```go
// EchoHandler echoes and formats text content based on provided query parameters
// swagger:route GET /echo echoText
//
// consumes:
//
// produces:
// - application/json
//
// schemes: http
//
// deprecated: false
//
// security:
// - api_key:
//
// parameters:
//
// - + name: data
// in: query
// description: the string to echo
// required: true
// type: string
//
// - + name: repetitions
// in: query
// description: the number of times to echo the string
// required: false
// type: integer
//
// responses:
// - 200: echoHandlerResponse
// - 500: errorResponse500
func EchoHandler(w http.ResponseWriter, r *http.Request) {
// code omitted for brevity
}
```

### Documenting data
#### Response data
Full features and documentation [here](https://goswagger.io/use/spec/response.html).
```go
// A 500 error returned in the response
// swagger:response errorResponse500
type errorResponse500 struct {
	// The error message
	// example: bad input
	Error string `json:"error"`
	// The status code
	// example: 500
	Status int `json:"status"`
}
```
#### Parameters
For more robust documentation of parameters than can be achieved in-line within a `swagger:route` comment. Supports more advanced functionality (e.g. enum types).

Full features and documentation [here](https://goswagger.io/use/spec/params.html).
```go
// swagger:enum echoCaseType
type echoCaseType string

const (
	upperCase echoCaseType = "upper"
	lowerCase echoCaseType = "lower"
	spongeBob echoCaseType = "spongebob"
)

// swagger:parameters echoText
type responseStringCase struct {
	// Description: Capitalization for response string
	// in: url
	// required: false
	// example: "upper"
	Case echoCaseType
}
```

#### Swagger models
For a more flexible implementation that can be used as both an input/parameter and a response type, use a [swagger:model](https://goswagger.io/use/spec/model.html). 
```go
// EchoHandlerResponse is the response from the echo handler endpoint
// swagger:model echoHandlerResponse
type EchoHandlerResponse struct {
	// The response string
	//
	// required: true
	// example: ECHO ECHO ECHO
	ResponseString string
}
```

<a id="generation-and-deployment"></a>
## Generating and deploying your docs
### Locally
A Makefile can be helpful for helping other teammates spin up and view docs.
```Makefile
install-swagger: go install github.com/go-swagger/go-swagger/cmd/swagger@latest
swagger: install-swagger
	swagger generate spec -o ./swagger.yaml --scan-models
serve-swagger: swagger
	swagger serve -F=swagger swagger.yaml
```

### Go HTTP server

We can easily deploy our Swagger YAML or JSON file as a documentation endpoint on our existing Go API using the [github.com/go-openapi/runtime/middleware](github.com/go-openapi/runtime/middleware) package.

Import the package:
```
go get github.com/go-openapi/runtime
```

Then, add endpoints for the raw Swagger YAML/JSON and the formatted docs.
```go
// Displays Swagger YAML/JSON only
http.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

// Displays formatted docs
opts := middleware.RedocOpts{SpecURL: "/swagger.yaml", Path: "/api/docs"}
sh := middleware.Redoc(opts, nil)
http.Handle("/api/docs", sh)

// ... Other code omitted for brevity

log.Fatal(http.ListenAndServe(":3333", nil))
```

## Final thoughts
### Pros
- Works via comments (which we should be doing anyways), thus documentation lives close to the code and close to current workflow
- No need to hand-write YAML, JSON, etc.
- Resultant specs can be used for code generation, documentation, myriad of purposes

### Cons
- Potential to add bloat to codebase; should integrate into code as much as possible vs. declaring custom types
- Plaintext comments don't always scale well
- Some IDEs will try to autoformat your comments, which can occasionally result in misformatted docs and breakages when go-swagger parses the comments -- need to be attentive to Git diffs when committing changes