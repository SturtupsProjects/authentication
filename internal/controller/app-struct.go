package controller

import (
	"authentification/internal/usecase"
	"authentification/internal/usecase/repo"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type Controller struct {
	Auth *usecase.UserUseCase
}

func NewController(db *sqlx.DB, log *slog.Logger) *Controller {

	authRepo := repo.NewUserRepo(db)

	ctr := &Controller{
		Auth: usecase.NewUserUseCase(authRepo, log),
	}

	return ctr
}
