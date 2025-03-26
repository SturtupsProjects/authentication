package app

import (
	"authentification/config"
	pc "authentification/internal/generated/company"
	pb "authentification/internal/generated/user"
	"authentification/internal/usecase"
	"authentification/internal/usecase/repo"
	"authentification/internal/usecase/token"
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

	userSt := repo.NewUserRepo(db)
	companySt := repo.NewCompanyStorage(db)
	branchSt := repo.NewBranchRepo(db)
	workerRepo := repo.NewWorkersRepo(db)

	authSr := usecase.NewAuthServiceServer(userSt, workerRepo, logger1, cfg)
	companySr := usecase.NewCompanyService(companySt, branchSt, logger1)

	// ------------ starting cron -------------------

	//{
	//	c := cron.New()
	//
	//	_, err = c.AddFunc("0 */3 * * *", func() {
	//		err := repo.BalanceChecker(db)
	//		if err != nil {
	//			log.Printf("Ошибка выполнения BalanceChecker: %v", err)
	//		} else {
	//			log.Println("BalanceChecker успешно выполнен")
	//		}
	//	})
	//
	//	if err != nil {
	//		log.Fatalf("Ошибка при добавлении задачи: %v", err)
	//	}
	//
	//	c.Start()
	//}

	log.Println("Cron dan otdi")

	// ----------------------- starting services --------------------
	{
		listen, err := net.Listen("tcp", cfg.RUN_PORT)
		fmt.Println("listening on port " + cfg.RUN_PORT)
		if err != nil {
			logger1.Error("Error listening on port " + cfg.RUN_PORT)
			log.Fatal(err)
		}

		service := grpc.NewServer()
		pb.RegisterAuthServiceServer(service, authSr)
		pc.RegisterCompanyServiceServer(service, companySr)
		if err := service.Serve(listen); err != nil {
			logger1.Error("Error starting server")
			log.Fatal(err)
		}
	}
}
