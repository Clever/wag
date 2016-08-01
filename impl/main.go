package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Clever/inter-service-api-testing/codegen-poc/generated"
	"github.com/gorilla/mux"
)

// TODO: Some initialization of config???

func main() {

	controller := generated.ControllerImpl{}

	router := generated.SetupServer(mux.NewRouter(), controller)
	server := &http.Server{
		// TODO: This should be configurable???
		Addr:    fmt.Sprintf(":8080"),
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
