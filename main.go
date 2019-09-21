package main

import (
	"log"
	"net/http"

	routeHandler "./pkg/routeHandler"
)

func main() {
	routeHandler.HandleHttpRoutes()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
