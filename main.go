package main

import (
	"log"
	"net/http"

	routeHandler "./pkg/routeHandler"
)

const (
	Port = ":8080"
)

func main() {
	routeHandler.HandleHttpRoutes()

	log.Fatal(http.ListenAndServe(Port, nil))
}
