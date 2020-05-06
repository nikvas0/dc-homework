package routes

import (
	"product_upload/handlers"

	"github.com/gorilla/mux"
)

func InitRoutesCommon(router *mux.Router) {
	router.HandleFunc("/products/upload", handlers.Upload).Methods("POST")
}

func InitRoutesV1(router *mux.Router) {
	InitRoutesCommon(router)
}

func InitRoutes(router *mux.Router) {
	InitRoutesV1(router.PathPrefix("/v1").Subrouter())
}
