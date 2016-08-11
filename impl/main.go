package main

import (
	"log"

	"github.com/Clever/wag/generated/server"
)

// TODO: Some initialization of config???

func main() {

	controller := server.ControllerImpl{}
	s := server.New(controller, 8080)
	log.Fatal(s.Serve())
}
