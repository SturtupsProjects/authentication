package usecase

import pb "authentification/pkg/generated/user"

type UsersRepo interface {
	AddAdmin(in *pb.MessageResponse) (*pb.MessageResponse, error)
	CreateUser(in *pb.UserRequest) (*pb.UserResponse, error)
	GetUser(in *pb.UserIDRequest) (*pb.UserResponse, error)
	GetListUser(in *pb.FilterUserRequest) (*pb.UserListResponse, error)
	DeleteUser(in *pb.UserIDRequest) (*pb.MessageResponse, error)
	UpdateUser(in *pb.UserRequest) (*pb.UserResponse, error)
	LogIn(in *pb.LogInRequest) (*pb.LogInResponse, string, error)
}
