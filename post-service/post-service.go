package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"social-network/common/config"
	postpb "social-network/common/proto/post"
	"social-network/post-service/handlers"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	serviceConfig := config.LoadConfig("common/config/config.json")
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	db = connectToDB(connStr)

	initDB()

	grpcServer := grpc.NewServer()

	postServer := handlers.NewPostServer(db, &serviceConfig)
	postpb.RegisterPostServiceServer(grpcServer, *postServer)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serviceConfig.PostService.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Post-service gRPC server running on port %d\n", serviceConfig.PostService.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
