// Package main Echo API documentation
//
// The Echo API echoes and formats content via query parameters
// passed into GET requests
//
// Terms Of Service:
// there are no TOS at this moment, use at your own risk
//
//	Schemes: http
//	BasePath: /api
//	Version: 0.0.1
//	Contact: Royce Chen<rchen@rvohealth.com>
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import (
	"fmt"
	"go-swagger-demo/handlers"
	"log"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

func main() {
	// Echo endpoint
	http.HandleFunc("/api/echo", handlers.EchoHandler)

	// Displays Swagger YAML/JSON only
	http.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// Displays formatted docs
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml", Path: "/api/docs"}
	sh := middleware.Redoc(opts, nil)
	http.Handle("/api/docs", sh)

	fmt.Println("Starting server on :3333...")
	log.Fatal(http.ListenAndServe(":3333", nil))
}
