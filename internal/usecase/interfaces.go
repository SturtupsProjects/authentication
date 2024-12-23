package usecase

import (
	"authentification/internal/generated/company"
	pb "authentification/internal/generated/user"
)

type UsersRepo interface {
	AddAdmin(in *pb.MessageResponse) (*pb.MessageResponse, error)
	CreateUser(in *pb.UserRequest) (*pb.UserResponse, error)
	GetUser(in *pb.UserIDRequest) (*pb.UserResponse, error)
	GetListUser(in *pb.FilterUserRequest) (*pb.UserListResponse, error)
	DeleteUser(in *pb.UserIDRequest) (*pb.MessageResponse, error)
	UpdateUser(in *pb.UserRequest) (*pb.UserResponse, error)
	LogIn(in *pb.LogInRequest) (*pb.LogInResponse, string, string, error)
	CreateClient(in *pb.ClientRequest) (*pb.ClientResponse, error)
	GetClient(in *pb.UserIDRequest) (*pb.ClientResponse, error)
	GetListClient(in *pb.FilterClientRequest) (*pb.ClientListResponse, error)
	UpdateClient(in *pb.ClientUpdateRequest) (*pb.ClientResponse, error)
	DeleteClient(in *pb.UserIDRequest) (*pb.MessageResponse, error)
}
type CompanyRepo interface {
	CreateCompany(req *company.CreateCompanyRequest) (*company.CompanyResponse, error)
	GetCompany(req *company.GetCompanyRequest) (*company.CompanyResponse, error)
	UpdateCompany(req *company.UpdateCompanyRequest) (*company.CompanyResponse, error)
	DeleteCompany(req *company.DeleteCompanyRequest) (*company.Message, error)
	ListCompanies(req *company.ListCompaniesRequest) (*company.ListCompaniesResponse, error)
	ListCompanyUsers(req *company.ListCompanyUsersRequest) (*company.ListCompanyUsersResponse, error)
	CreateUserToCompany(req *company.CreateUserToCompanyRequest) (*company.Id, error)
}
