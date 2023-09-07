// Package main Echo API documentation
package main

import (
	"fmt"
	"go-swagger-demo/handlers"
	"log"
	"net/http"
)

func main() {
	// Echo endpoint
	http.HandleFunc("/api/echo", handlers.EchoHandler)

	fmt.Println("Starting server on :3333...")
	log.Fatal(http.ListenAndServe(":3333", nil))
}
