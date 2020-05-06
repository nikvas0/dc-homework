package middleware

import (
	"lib/auth_middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func InitMiddleware(router *mux.Router) {
	router.Use(auth_middleware.GetAuthMiddleware(
		func(r *http.Request) bool {
			return r.Method == "PUT"
		}))
	return
}
