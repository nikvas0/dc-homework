package routes

import (
	"github.com/gorilla/mux"
	"github.com/nikvas0/dc-homework/auth/handlers"
)

func InitRoutesCommon(router *mux.Router) {
	router.HandleFunc("/signup", handlers.SignUp).Methods("POST")
	router.HandleFunc("/signin", handlers.SignIn).Methods("POST")
	router.HandleFunc("/validate", handlers.Validate).Methods("POST")
	router.HandleFunc("/refresh", handlers.Refresh).Methods("POST")
}

func InitRoutesV1(router *mux.Router) {
	InitRoutesCommon(router)
}

func InitRoutes(router *mux.Router) {
	InitRoutesV1(router.PathPrefix("/v1").Subrouter())
}
