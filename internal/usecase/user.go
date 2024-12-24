package usecase

import (
	"authentification/config"
	"authentification/internal/entity"
	"authentification/internal/usecase/help"
	"authentification/internal/usecase/token"
	"context"
	"fmt"
	"log/slog"

	pb "authentification/internal/generated/user"
	"authentification/internal/usecase/repo"
)

type AuthServiceServer struct {
	repo *repo.UserRepo
	log  *slog.Logger
	conf *config.Config
	pb.UnimplementedAuthServiceServer
}

func NewAuthServiceServer(repo *repo.UserRepo, log *slog.Logger, conf *config.Config) *AuthServiceServer {
	return &AuthServiceServer{repo: repo, log: log, conf: conf}
}

func (s *AuthServiceServer) RegisterAdmin(ctx context.Context, req *pb.MessageResponse) (*pb.MessageResponse, error) {
	s.log.Info("RegisterAdmin called", "phone_number")

	pass, err := help.HashPassword(req.Message)
	if err != nil {
		s.log.Error("Failed to hash password", "error", err)
		return nil, fmt.Errorf("could not hash password: %w", err)
	}

	req.Message = pass
	resp, err := s.repo.AddAdmin(req)
	if err != nil {
		s.log.Error("Failed to register admin", "error", err)
		return nil, fmt.Errorf("could not register admin: %w", err)
	}

	s.log.Info("Admin registered successfully", "message", resp.Message)
	return resp, nil
}

// AddUser handles new user creation
func (s *AuthServiceServer) AddUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	s.log.Info("AddUser called", "email", req.Email)
	pass, err := help.HashPassword(req.Password)
	if err != nil {
		s.log.Error("Failed to hash password", "error", err)
		return nil, fmt.Errorf("could not hash password: %w", err)
	}
	req.Password = pass
	resp, err := s.repo.CreateUser(req)
	if err != nil {
		s.log.Error("Failed to create user", "error", err)
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	s.log.Info("User created successfully", "user_id", resp.UserId)
	return resp, nil
}

// GetUser retrieves a user by their ID
func (s *AuthServiceServer) GetUser(ctx context.Context, req *pb.UserIDRequest) (*pb.UserResponse, error) {
	s.log.Info("GetUser called", "user_id", req.Id)

	user, err := s.repo.GetUser(req)
	if err != nil {
		s.log.Error("Failed to retrieve user", "user_id", req.Id, "error", err)
		return nil, fmt.Errorf("could not retrieve user: %w", err)
	}

	s.log.Info("User retrieved successfully", "user_id", user.UserId)
	return user, nil
}

// GetUserList retrieves a list of users based on filters
func (s *AuthServiceServer) GetUserList(ctx context.Context, req *pb.FilterUserRequest) (*pb.UserListResponse, error) {
	s.log.Info("GetUserList called", "filters", req)

	users, err := s.repo.GetListUser(req)
	if err != nil {
		s.log.Error("Failed to retrieve user list", "error", err)
		return nil, fmt.Errorf("could not retrieve user list: %w", err)
	}

	s.log.Info("User list retrieved successfully", "count", len(users.Users))
	return users, nil
}

// UpdateUser updates a user's information
func (s *AuthServiceServer) UpdateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	s.log.Info("UpdateUser called", "user_id", req.UserId)
	if req.Password != "" {
		pass, err := help.HashPassword(req.Password)
		if err != nil {
			s.log.Error("Failed to hash password", "error", err)
			return nil, fmt.Errorf("could not hash password: %w", err)
		}
		req.Password = pass

	}
	updatedUser, err := s.repo.UpdateUser(req)
	if err != nil {
		s.log.Error("Failed to update user", "user_id", req.UserId, "error", err)
		return nil, fmt.Errorf("could not update user: %w", err)
	}

	s.log.Info("User updated successfully", "user_id", updatedUser.UserId)
	return updatedUser, nil
}

// DeleteUser removes a user from the system
func (s *AuthServiceServer) DeleteUser(ctx context.Context, req *pb.UserIDRequest) (*pb.MessageResponse, error) {
	s.log.Info("DeleteUser called", "user_id", req.Id)

	resp, err := s.repo.DeleteUser(req)
	if err != nil {
		s.log.Error("Failed to delete user", "user_id", req.Id, "error", err)
		return nil, fmt.Errorf("could not delete user: %w", err)
	}

	s.log.Info("User deleted successfully", "message", resp.Message)
	return resp, nil
}

// LogIn handles user login
func (s *AuthServiceServer) LogIn(ctx context.Context, req *pb.LogInRequest) (*pb.TokenResponse, error) {
	s.log.Info("LogIn called", "phone_number", req.PhoneNumber)

	loginResp, pass, CompanyId, err := s.repo.LogIn(req)
	if err != nil {
		s.log.Error("Login failed", "phone_number", req.PhoneNumber, "error", err)
		return nil, fmt.Errorf("login failed: %w", err)
	}
	ok := help.CheckPasswordHash(req.Password, pass)
	if ok == false {
		s.log.Error("Login failed", "phone_number", req.PhoneNumber, "error", err)
		return nil, fmt.Errorf("login failed: %s", "Invalid password")
	}
	// Simulate token generation (Replace with real token generation)
	accessToken, err := token.GenerateAccessToken(&entity.LogInToken{UserId: loginResp.UserId, Role: loginResp.Role, FirstName: loginResp.FirstName, PhoneNumber: loginResp.PhoneNumber, CompanyId: CompanyId})
	if err != nil {
		s.log.Error("Error in generating access token", "error", err)
		return nil, err
	}

	refreshToken, err := token.GenerateRefreshToken(&entity.LogInToken{UserId: loginResp.UserId, Role: loginResp.Role, FirstName: loginResp.FirstName, PhoneNumber: loginResp.PhoneNumber, CompanyId: CompanyId})
	if err != nil {
		s.log.Error("Error in generating refresh token", "error", err)
		return nil, err
	}

	expireAt := token.GetExpires()
	s.log.Info("Login successful", "user_id", loginResp.UserId)
	return &pb.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpireAt:     int32(expireAt),
	}, nil
}
