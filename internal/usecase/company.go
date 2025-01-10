package usecase

import (
	"authentification/internal/generated/company"
	"authentification/internal/usecase/help"
	"context"
	"log/slog"
)

type CompanyService struct {
	repoC CompanyRepo
	repoB BranchRepo
	log   *slog.Logger
	company.UnimplementedCompanyServiceServer
}

func NewCompanyService(repoC CompanyRepo, repoB BranchRepo, log *slog.Logger) *CompanyService {
	return &CompanyService{
		repoC: repoC,
		repoB: repoB,
		log:   log,
	}
}

func (s *CompanyService) CreateCompany(ctx context.Context, req *company.CreateCompanyRequest) (*company.CompanyResponse, error) {
	s.log.Info("CreateCompany called", "request", req)
	result, err := s.repoC.CreateCompany(req)
	if err != nil {
		s.log.Error("Error creating company", "error", err)
		return nil, err
	}
	s.log.Info("CreateCompany finished")
	return result, nil
}

func (s *CompanyService) GetCompany(ctx context.Context, req *company.GetCompanyRequest) (*company.CompanyResponse, error) {
	s.log.Info("GetCompany called", "company_id", req.CompanyId)
	result, err := s.repoC.GetCompany(req)
	if err != nil {
		s.log.Error("Error fetching company", "error", err)
		return nil, err
	}
	s.log.Info("GetCompany finished")
	return result, nil
}

func (s *CompanyService) UpdateCompany(ctx context.Context, req *company.UpdateCompanyRequest) (*company.CompanyResponse, error) {
	s.log.Info("UpdateCompany called", "company_id", req.CompanyId)
	result, err := s.repoC.UpdateCompany(req)
	if err != nil {
		s.log.Error("Error updating company", "error", err)
		return nil, err
	}
	s.log.Info("UpdateCompany finished")
	return result, nil
}

func (s *CompanyService) DeleteCompany(ctx context.Context, req *company.DeleteCompanyRequest) (*company.Message, error) {
	s.log.Info("DeleteCompany called", "company_id", req.CompanyId)
	result, err := s.repoC.DeleteCompany(req)
	if err != nil {
		s.log.Error("Error deleting company", "error", err)
		return nil, err
	}
	s.log.Info("DeleteCompany finished")
	return result, nil
}

func (s *CompanyService) ListCompanies(ctx context.Context, req *company.ListCompaniesRequest) (*company.ListCompaniesResponse, error) {
	s.log.Info("ListCompanies called", "pagination", req)
	result, err := s.repoC.ListCompanies(req)
	if err != nil {
		s.log.Error("Error listing companies", "error", err)
		return nil, err
	}
	s.log.Info("ListCompanies finished")
	return result, nil
}

func (s *CompanyService) ListCompanyUsers(ctx context.Context, req *company.ListCompanyUsersRequest) (*company.ListCompanyUsersResponse, error) {
	s.log.Info("ListCompanyUsers called", "company_id", req.CompanyId)
	result, err := s.repoC.ListCompanyUsers(req)
	if err != nil {
		s.log.Error("Error listing company users", "error", err)
		return nil, err
	}
	s.log.Info("ListCompanyUsers finished")
	return result, nil
}

func (s *CompanyService) CreateUserToCompany(ctx context.Context, req *company.CreateUserToCompanyRequest) (*company.Id, error) {
	s.log.Info("CreateUserToCompany called", "request", req)
	hashedPassword, err := help.HashPassword(req.Password)
	if err != nil {
		s.log.Error("Error in hashing password", "error", err)
		return nil, err
	}

	req.Password = hashedPassword
	result, err := s.repoC.CreateUserToCompany(req)
	if err != nil {
		s.log.Error("Error creating user to company", "error", err)
		return nil, err
	}
	s.log.Info("CreateUserToCompany finished")
	return result, nil
}
