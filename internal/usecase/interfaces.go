package usecase

import "authentification/internal/entity"

type UsersRepo interface {
	AddAdmin(in entity.AdminPass) (entity.Message, error)
	CreateUser(in entity.User) (entity.UserRequest, error)
	GetUser(in entity.UserID) (entity.UserRequest, error)
	GetListUser(in entity.FilterUser) (entity.UserList, error)
	DeleteUser(in entity.UserID) (entity.Message, error)
	UpdateUser(in entity.UserRequest) (entity.UserRequest, error)
	LogIn(in entity.PhoneNumber) (entity.LogInReq, error)
}
