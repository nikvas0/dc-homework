package auth_middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	pb "lib/proto/auth"

	"google.golang.org/grpc"
)

func GetAuthMiddleware(needAuth func(r *http.Request) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				next.ServeHTTP(w, r)
				return
			}
			if !needAuth(r) {
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

			conn, err := grpc.Dial("auth:50051", grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Printf("did not connect: %v", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			defer conn.Close()
			c := pb.NewAuthServiceClient(conn)

			gctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			reply, err := c.Validate(gctx, &pb.ValidateRequest{Token: splittedAuth[1]})
			if err != nil {
				log.Printf("Could not authorize: %v", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			log.Printf("Auth success: %d %s", reply.GetUser(), reply.GetEmail)

			ctx := context.WithValue(r.Context(), "user_id", reply.GetUser())
			ctx = context.WithValue(ctx, "email", reply.GetEmail())
			ctx = context.WithValue(ctx, "role", reply.GetRole())
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
			return
		})
	}
}
