package middleware

import (
	"github.com/gorilla/mux"
	"github.com/nikvas0/dc-homework/lib/auth_middleware"
)

func InitMiddleware(router *mux.Router) {
	router.Use(auth_middleware.GetAuthMiddleware())
	return
}
