package usecase

import (
	"authentification/internal/entity"
	"authentification/internal/usecase/help"
	"authentification/internal/usecase/token"
	"log/slog"
)

type UserUseCase struct {
	repo UsersRepo
	log  *slog.Logger
}

func NewUserUseCase(repo UsersRepo, log *slog.Logger) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		log:  log,
	}
}

func (u *UserUseCase) RegisterAdmin(in entity.AdminPass) (entity.Message, error) {
	hash, err := help.HashPassword(in.Password)
	if err != nil {
		u.log.Error("Error in help password", "error", err)
		return entity.Message{}, err
	}

	in.Password = hash

	res, err := u.repo.AddAdmin(in)
	if err != nil {
		u.log.Error("Error in adding admin", "error", err)
		return entity.Message{}, err
	}

	return res, nil
}

func (u *UserUseCase) AddUser(in entity.User) (entity.UserRequest, error) {
	hash, err := help.HashPassword(in.Password)
	if err != nil {
		u.log.Error("Error in help password", "error", err)
		return entity.UserRequest{}, err
	}

	in.Password = hash

	res, err := u.repo.CreateUser(in)
	if err != nil {
		u.log.Error("Error in adding user", "error", err)
		return entity.UserRequest{}, err
	}

	return res, nil
}

func (u *UserUseCase) GetUser(in entity.UserID) (entity.UserRequest, error) {
	res, err := u.repo.GetUser(in)
	if err != nil {
		u.log.Error("Error in getting user", "error", err)
		return entity.UserRequest{}, err
	}

	return res, nil
}

func (u *UserUseCase) UpdateUser(in entity.UserRequest) (entity.UserRequest, error) {
	res, err := u.repo.UpdateUser(in)
	if err != nil {
		u.log.Error("Error in updating user", "error", err)
		return entity.UserRequest{}, err
	}

	return res, nil
}

func (u *UserUseCase) GetUserList(in entity.FilterUser) (entity.UserList, error) {
	res, err := u.repo.GetListUser(in)
	if err != nil {
		u.log.Error("Error in getting user list", "error", err)
		return entity.UserList{}, err
	}

	return res, nil
}

func (u *UserUseCase) DeleteUser(in entity.UserID) (entity.Message, error) {
	res, err := u.repo.DeleteUser(in)
	if err != nil {
		u.log.Error("Error in deleting user", "error", err)
		return entity.Message{}, err
	}

	return res, nil
}

func (u *UserUseCase) LogIn(in entity.LogIn) (entity.Token, error) {
	phone := entity.PhoneNumber{in.PhoneNumber}

	res, err := u.repo.LogIn(phone)
	if err != nil {
		u.log.Error("Error in logging in", "error", err)
		return entity.Token{}, err
	}

	accessToken, err := token.GenerateAccessToken(res)
	if err != nil {
		u.log.Error("Error in generating access token", "error", err)
		return entity.Token{}, err
	}

	refreshToken, err := token.GenerateRefreshToken(res)
	if err != nil {
		u.log.Error("Error in generating refresh token", "error", err)
		return entity.Token{}, err
	}

	expireAt := token.GetExpires()

	return entity.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpireAt:     expireAt,
	}, nil
}
