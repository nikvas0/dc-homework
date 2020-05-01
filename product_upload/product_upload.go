package main

import (
	"log"
	"net/http"
	"os"

	"product_upload/middleware"
	"product_upload/queues"
	"product_upload/routes"

	"github.com/gorilla/mux"
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
