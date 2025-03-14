package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gnad103/go-ex/api-gateway/handlers"
	pb "github.com/gnad103/go-ex/proto"
)

func main() {
	// Connect to the Go user service
	userConn, err := grpc.Dial("go-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	defer userConn.Close()
	userClient := pb.NewUserServiceClient(userConn)

	// Connect to the Python product service
	productConn, err := grpc.Dial("python-service:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to product service: %v", err)
	}
	defer productConn.Close()
	productClient := pb.NewProductServiceClient(productConn)

	// Create API handler
	apiHandler := handlers.NewAPIHandler(userClient, productClient)

	// Create router
	router := mux.NewRouter()

	// User routes
	router.HandleFunc("/api/users", apiHandler.ListUsers).Methods("GET")
	router.HandleFunc("/api/users", apiHandler.CreateUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", apiHandler.GetUser).Methods("GET")

	// Product routes
	router.HandleFunc("/api/products", apiHandler.ListProducts).Methods("GET")
	router.HandleFunc("/api/products", apiHandler.CreateProduct).Methods("POST")
	router.HandleFunc("/api/products/{id}", apiHandler.GetProduct).Methods("GET")

	// Create HTTP server
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("API Gateway is running on port 8080...")
	log.Fatal(srv.ListenAndServe())
}
