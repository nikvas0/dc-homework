package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"auth/queues"
	"auth/routes"
	"auth/rpc_server"
	"auth/storage"
	pb "lib/proto/auth"

	"auth/middleware"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
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

	err = queues.Init(os.Getenv("RABBITMQ"))
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer queues.Close()

	router := mux.NewRouter().StrictSlash(true)
	middleware.InitMiddleware(router)
	routes.InitRoutes(router)

	lis, err := net.Listen("tcp", os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &rpc_server.AuthServer{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", router))
}
