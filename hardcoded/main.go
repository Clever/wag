package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO: Some initialization of config???

func main() {

	controller = ControllerImpl{}

	router := withMiddleware(withRoutes(mux.NewRouter()))
	server := &http.Server{
		// TODO: This should be configurable???
		Addr:    fmt.Sprintf(":8080"),
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
