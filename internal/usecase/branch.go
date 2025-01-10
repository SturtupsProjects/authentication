package usecase

import (
	"authentification/internal/generated/company"
	"context"
)

func (s *CompanyService) CreateBranch(ctx context.Context, req *company.CreateBranchRequest) (*company.BranchResponse, error) {
	s.log.Info("CreateBranch called", "request", req)
	result, err := s.repoB.CreateBranch(req)
	if err != nil {
		s.log.Error("Error creating branch", "error", err)
		return nil, err
	}
	s.log.Info("CreateBranch finished")
	return result, nil
}

func (s *CompanyService) GetBranch(ctx context.Context, req *company.GetBranchRequest) (*company.BranchResponse, error) {
	s.log.Info("GetBranch called", "branch_id", req.BranchId)
	result, err := s.repoB.GetBranch(req)
	if err != nil {
		s.log.Error("Error fetching branch", "error", err)
		return nil, err
	}
	s.log.Info("GetBranch finished")
	return result, nil
}

func (s *CompanyService) UpdateBranch(ctx context.Context, req *company.UpdateBranchRequest) (*company.BranchResponse, error) {
	s.log.Info("UpdateBranch called", "branch_id", req.BranchId)
	result, err := s.repoB.UpdateBranch(req)
	if err != nil {
		s.log.Error("Error updating branch", "error", err)
		return nil, err
	}
	s.log.Info("UpdateBranch finished")
	return result, nil
}

func (s *CompanyService) DeleteBranch(ctx context.Context, req *company.DeleteBranchRequest) (*company.Message, error) {
	s.log.Info("DeleteBranch called", "branch_id", req.BranchId)
	result, err := s.repoB.DeleteBranch(req)
	if err != nil {
		s.log.Error("Error deleting branch", "error", err)
		return nil, err
	}
	s.log.Info("DeleteBranch finished")
	return result, nil
}

func (s *CompanyService) ListBranches(ctx context.Context, req *company.ListBranchesRequest) (*company.ListBranchesResponse, error) {
	s.log.Info("ListBranches called", "pagination", req)
	result, err := s.repoB.ListBranches(req)
	if err != nil {
		s.log.Error("Error listing branches", "error", err)
		return nil, err
	}
	s.log.Info("ListBranches finished")
	return result, nil
}
