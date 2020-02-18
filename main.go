package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nikvas0/dc-homework/handlers"
	"github.com/nikvas0/dc-homework/storage"
)

func main() {
	//err := storage.Init("sqlite3", "test.db")
	err := storage.Init("postgres", "host=localhost port=5432 user=myuser dbname=mydb password=mypass")
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer storage.Close()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/product", handlers.CreateProduct).Methods("POST")
	router.HandleFunc("/product/{id:[0-9]+}", handlers.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/products", handlers.GetProducts).Methods("GET").Queries("offset", "{offset:[0-9]+}", "limit", "{limit:[0-9]+}")
	router.HandleFunc("/product/{id:[0-9]+}", handlers.GetProduct).Methods("GET")
	router.HandleFunc("/product", handlers.UpdateProduct).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))
}
