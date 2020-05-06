package routes

import (
	"auth/handlers"

	"github.com/gorilla/mux"
)

func InitRoutesCommon(router *mux.Router) {
	router.HandleFunc("/signup", handlers.SignUp).Methods("POST")
	router.HandleFunc("/signin", handlers.SignIn).Methods("POST")
	router.HandleFunc("/validate", handlers.Validate).Methods("POST")
	router.HandleFunc("/refresh", handlers.Refresh).Methods("POST")
	router.HandleFunc("/confirm/{token}", handlers.Confirm).Methods("GET")

	router.HandleFunc("/role", handlers.UpdateRole).Methods("PUT")
}

func InitRoutesV1(router *mux.Router) {
	InitRoutesCommon(router)
}

func InitRoutes(router *mux.Router) {
	InitRoutesV1(router.PathPrefix("/v1").Subrouter())
}
