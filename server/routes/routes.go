package routes

import (
	"github.com/gorilla/mux"
	"github.com/nikvas0/dc-homework/server/handlers"
)

func InitRoutesCommon(router *mux.Router) {
	router.HandleFunc("/product", handlers.CreateProduct).Methods("POST")
	router.HandleFunc("/product/{id:[0-9]+}", handlers.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/product/{id:[0-9]+}", handlers.GetProduct).Methods("GET")
	router.HandleFunc("/product", handlers.UpdateProduct).Methods("PUT")
}

func InitRoutesV1(router *mux.Router) {
	InitRoutesCommon(router)

	router.HandleFunc("/products", handlers.GetAllProducts).Methods("GET")
}

func InitRoutesV2(router *mux.Router) {
	InitRoutesCommon(router)
	// Пример немного искусственный (можно было просто детектить параметры в первой версии), но надо было что-то сделать.
	router.HandleFunc("/products", handlers.GetProductsPage).Methods("GET").Queries("offset", "{offset:[0-9]+}", "limit", "{limit:[0-9]+}")
}

func InitRoutes(router *mux.Router) {
	InitRoutesV1(router.PathPrefix("/v1").Subrouter())
	InitRoutesV2(router.PathPrefix("/v2").Subrouter())
}
