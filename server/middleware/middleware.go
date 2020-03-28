package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nikvas0/dc-homework/server/objects"
)

func GetAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				next.ServeHTTP(w, r)
				return
			}

			authString := r.Header.Get("Authorization")
			splittedAuth := strings.Split(authString, " ")
			if len(splittedAuth) != 2 || splittedAuth[0] != "Bearer" {
				log.Println("Auth middleware error: bad token.")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			body, err := json.Marshal(map[string]interface{}{
				"token": splittedAuth[1],
			})
			if err != nil {
				log.Println("Auth middleware error: failed to marshal token.")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			req, err := http.NewRequest("POST", "http://auth:8080/v1/validate", bytes.NewBuffer(body))
			if err != nil {
				log.Println("Auth middleware error: failed to create request.")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			req.Header.Add("Content-Type", "application/json")

			client := http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Println("Auth middleware error: failed to make request: %w.", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if resp.StatusCode != http.StatusOK {
				log.Println("Auth middleware error: failed to authorize.")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			defer resp.Body.Close()
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Auth middleware error: failed to read response.")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userData := objects.UserData{}
			if err := json.Unmarshal(respBody, &userData); err != nil {
				log.Println("Auth middleware error: failed to unmarshal response.")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			log.Printf("Auth success: %d %s", userData.ID, userData.Email)

			ctx := context.WithValue(r.Context(), "user_id", userData.ID)
			ctx = context.WithValue(ctx, "email", userData.Email)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
			return
		})
	}
}

func InitMiddleware(router *mux.Router) {
	router.Use(GetAuthMiddleware())
	return
}
