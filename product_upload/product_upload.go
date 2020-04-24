package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nikvas0/dc-homework/product_upload/middleware"
	"github.com/nikvas0/dc-homework/product_upload/queues"
	"github.com/nikvas0/dc-homework/product_upload/routes"
)

func main() {
	err := queues.Init(os.Getenv("RABBITMQ"))
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer queues.Close()

	router := mux.NewRouter().StrictSlash(true)
	middleware.InitMiddleware(router)
	routes.InitRoutes(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
