package main

import (
	"log"
	"net"

	"github.com/liju-github/EcommerceUserService/configs"
	"github.com/liju-github/EcommerceUserService/db"
	"github.com/liju-github/EcommerceUserService/proto/user"
	"github.com/liju-github/EcommerceUserService/repository"
	"github.com/liju-github/EcommerceUserService/service"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()
	dbConn, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := repository.NewUserRepository(dbConn)
	userService := service.NewUserService(userRepo)

	listener, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}

	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, userService)

	log.Println("Server is running")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
