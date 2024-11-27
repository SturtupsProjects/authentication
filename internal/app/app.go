package app

import (
	"authentification/config"
	"authentification/internal/usecase"
	"authentification/internal/usecase/repo"
	"authentification/internal/usecase/token"
	pb "authentification/pkg/generated/user"
	"authentification/pkg/logger"
	"authentification/pkg/postgres"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Run(cfg *config.Config) {

	logger1 := logger.NewLogger()

	db, err := postgres.Connection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = token.ConfigToken(cfg)
	if err != nil {
		log.Fatal(err)
	}
	userst := repo.NewUserRepo(db)

	authSr := usecase.NewAuthServiceServer(userst, logger1, cfg)

	listen, err := net.Listen("tcp", cfg.RUN_PORT)
	fmt.Println("listening on port " + cfg.RUN_PORT)
	if err != nil {
		logger1.Error("Error listening on port " + cfg.RUN_PORT)
		log.Fatal(err)
	}

	service := grpc.NewServer()
	pb.RegisterAuthServiceServer(service, authSr)

	if err := service.Serve(listen); err != nil {
		logger1.Error("Error starting server")
		log.Fatal(err)
	}
}
