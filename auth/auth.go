package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nikvas0/dc-homework/auth/routes"
	"github.com/nikvas0/dc-homework/auth/storage"
)

func main() {
	err := storage.Init(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
			os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_SSLMODE")))

	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer storage.Close()

	router := mux.NewRouter().StrictSlash(true)
	routes.InitRoutes(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
