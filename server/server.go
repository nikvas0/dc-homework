package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"server/middleware"
	"server/routes"
	"server/storage"

	"github.com/gorilla/mux"
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
	middleware.InitMiddleware(router)
	routes.InitRoutes(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
