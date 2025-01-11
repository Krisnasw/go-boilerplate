package routers

import (
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"go-grst-boilerplate/app/users-svc/handlers"
	"go-grst-boilerplate/app/users-svc/repositories"
	"go-grst-boilerplate/app/users-svc/services"
	userpb "go-grst-boilerplate/contracts"
)

func InitMicroservices(db *gorm.DB, grpcServer *grpc.Server) {
	// Initialize repositories, services, and handlers
	userRepository := repositories.New(db)
	userService := services.New(userRepository)
	userHandler := handlers.New(userService)

	// Register gRPC server
	userpb.RegisterUserServiceServer(grpcServer, userHandler)
}
