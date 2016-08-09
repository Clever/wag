package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Clever/inter-service-api-testing/codegen-poc/generated/server"
	"github.com/gorilla/mux"
)

// TODO: Some initialization of config???

func main() {

	controller := server.ControllerImpl{}

	router := server.SetupServer(mux.NewRouter(), controller)
	server := &http.Server{
		// TODO: This should be configurable???
		Addr:    fmt.Sprintf(":8080"),
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
