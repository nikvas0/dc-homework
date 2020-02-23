package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nikvas0/dc-homework/server/routes"
	"github.com/nikvas0/dc-homework/server/storage"
)

func main() {
	//err := storage.Init("sqlite3", "test.db")
	err := storage.Init("postgres", "host=localhost port=5432 user=myuser dbname=mydb password=mypass")
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer storage.Close()

	router := mux.NewRouter().StrictSlash(true)
	routes.InitRoutes(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
