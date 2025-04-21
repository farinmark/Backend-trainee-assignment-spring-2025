package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/farinmark/Backend-trainee-assignment-spring-2025/internal/config"
	"github.com/farinmark/Backend-trainee-assignment-spring-2025/internal/db"
	"github.com/farinmark/Backend-trainee-assignment-spring-2025/internal/handlers"
	"github.com/farinmark/Backend-trainee-assignment-spring-2025/internal/middleware"
	"github.com/farinmark/Backend-trainee-assignment-spring-2025/internal/service"
	pb "github.com/farinmark/Backend-trainee-assignment-spring-2025/proto"
)

func main() {
	cfg := config.Load()

	dbConn := db.Connect(cfg)
	defer dbConn.Close()
	db.Migrate(dbConn)

	svc := service.New(dbConn)

	r := mux.NewRouter()
	r.Use(middleware.Logging)
	r.HandleFunc("/dummyLogin", handlers.DummyLogin).Methods("POST")
	r.HandleFunc("/register", handlers.Register(svc)).Methods("POST")
	r.HandleFunc("/login", handlers.Login(svc)).Methods("POST")

	auth := r.PathPrefix("").Subrouter()
	auth.Use(middleware.Authenticate(svc))
	auth.HandleFunc("/pvz", handlers.CreatePVZ(svc)).Methods("POST")
	auth.HandleFunc("/pvz", handlers.ListPVZ(svc)).Methods("GET")
	auth.HandleFunc("/pvz/{id:[0-9]+}/sessions", handlers.StartSession(svc)).Methods("POST")
	auth.HandleFunc("/pvz/{id:[0-9]+}/sessions/{sid:[0-9]+}/items", handlers.AddItem(svc)).Methods("POST")
	auth.HandleFunc("/pvz/{id:[0-9]+}/sessions/{sid:[0-9]+}/items", handlers.DeleteItem(svc)).Methods("DELETE")
	auth.HandleFunc("/pvz/{id:[0-9]+}/sessions/{sid:[0-9]+}/close", handlers.CloseSession(svc)).Methods("POST")

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":9000", nil))
	}()

	go func() {
		srv := &http.Server{
			Handler:      r,
			Addr:         ":8080",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		log.Println("HTTP server listening on :8080")
		log.Fatal(srv.ListenAndServe())
	}()

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("gRPC listen error: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPVZServiceServer(grpcServer, handlers.NewGRPCServer(svc))
	log.Println("gRPC server listening on :3000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}
