package usecase

import (
	"authentification/internal/generated/company"
	"authentification/internal/usecase/help"
	"authentification/pkg/eskiz"
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

func (s *CompanyService) CreateCompanyBalance(ctx context.Context, req *company.CompanyBalanceRequest) (*company.CompanyBalanceResponse, error) {
	s.log.Info("CreateCompanyBalance called", "request", req)
	result, err := s.repoC.CreateBalance(req)
	if err != nil {
		s.log.Error("Error creating company balance", "error", err)
		return nil, err
	}
	s.log.Info("CreateCompanyBalance finished")
	return result, nil
}
func (s *CompanyService) GetCompanyBalance(ctx context.Context, req *company.Id) (*company.CompanyBalanceResponse, error) {
	s.log.Info("GetCompanyBalance called", "company_id", req.Id)
	result, err := s.repoC.GetBalance(req)
	if err != nil {
		s.log.Error("Error fetching company balance", "error", err)
		return nil, err
	}
	s.log.Info("GetCompanyBalance finished")
	return result, nil

}

func (s *CompanyService) UpdateCompanyBalance(ctx context.Context, req *company.CompanyBalanceRequest) (*company.CompanyBalanceResponse, error) {
	s.log.Info("UpdateCompanyBalance called", "company_id", req.CompanyId)
	result, err := s.repoC.UpdateBalance(req)
	if err != nil {
		s.log.Error("Error updating company balance", "error", err)
		return nil, err
	}
	s.log.Info("UpdateCompanyBalance finished")
	return result, nil
}

func (s *CompanyService) GetUsersBalanceList(ctx context.Context, req *company.FilterCompanyBalanceRequest) (*company.CompanyBalanceListResponse, error) {
	s.log.Info("GetUsersBalanceList called", "limit", req.Limit)
	result, err := s.repoC.ListBalances(req)
	if err != nil {
		s.log.Error("Error fetching company balance", "error", err)
		return nil, err
	}
	s.log.Info("GetUsersBalanceList finished")
	return result, nil
}

func (s *CompanyService) DeleteCompanyBalance(ctx context.Context, req *company.Id) (*company.Message, error) {
	s.log.Info("DeleteCompanyBalance called", "company_id", req.Id)
	result, err := s.repoC.DeleteBalance(req)
	if err != nil {
		s.log.Error("Error deleting company balance", "error", err)
		return nil, err
	}
	s.log.Info("DeleteCompanyBalance finished")
	return result, nil
}

func (s *CompanyService) SendSMS(ctx context.Context, req *company.SmsRequest) (*company.Message, error) {
	s.log.Info("sending SMS")
	res, err := s.repoC.GetBalance(&company.Id{Id: req.CompanyId})
	if err != nil {
		s.log.Error("failed to get balance")
		return nil, err
	}
	if res.Balance > 200 {
		token, err := eskiz.LoginToEskiz("ozodjon160605@gmail.com", "AUlv8VzGxAkhv2BbHmoqYLJrLpMPKxddPhUJD40g")
		if err != nil {
			s.log.Error("failed to login to Eskiz")
			return nil, err
		}
		err = eskiz.SendNotification(token, req.Phone, req.Message)
		if err != nil {
			s.log.Error("failed to send notification")
			return nil, err
		}
	} else {
		return &company.Message{Message: "Not enough balance"}, nil
	}

	return &company.Message{Message: "Message sent"}, nil
}
