package repo

import (
	"authentification/internal/generated/company"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type BranchRepo struct {
	db *sqlx.DB
}

func NewBranchRepo(db *sqlx.DB) *BranchRepo {
	return &BranchRepo{db: db}
}

func (r *BranchRepo) CreateBranch(in *company.CreateBranchRequest) (*company.BranchResponse, error) {
	var result company.BranchResponse
	query := `INSERT INTO branches (company_id, name, address, phone_number) 
              VALUES ($1, $2, $3, $4) 
              RETURNING branch_id, company_id, name, address, phone_number, created_at, updated_at`
	err := r.db.QueryRow(query, in.CompanyId, in.Name, in.Address, in.PhoneNumber).Scan(
		&result.BranchId,
		&result.CompanyId,
		&result.Name,
		&result.Address,
		&result.PhoneNumber,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create branch: %v", err)
	}
	return &result, nil
}

func (r *BranchRepo) GetBranch(in *company.GetBranchRequest) (*company.BranchResponse, error) {
	var result company.BranchResponse
	query := `SELECT branch_id, company_id, name, address, phone_number, created_at, updated_at 
              FROM branches 
              WHERE branch_id = $1`
	err := r.db.QueryRow(query, in.BranchId).Scan(
		&result.BranchId,
		&result.CompanyId,
		&result.Name,
		&result.Address,
		&result.PhoneNumber,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("branch not found")
		}
		return nil, err
	}
	return &result, nil
}

func (r *BranchRepo) UpdateBranch(in *company.UpdateBranchRequest) (*company.BranchResponse, error) {
	var result company.BranchResponse
	query := `UPDATE branches 
              SET name = COALESCE($1, name), 
                  address = COALESCE($2, address), 
                  phone_number = COALESCE($3, phone_number),
                  updated_at = NOW()
              WHERE branch_id = $4 
              RETURNING branch_id, company_id, name, address, phone_number, created_at, updated_at`
	err := r.db.QueryRow(query, in.Name, in.Address, in.PhoneNumber, in.BranchId).Scan(
		&result.BranchId,
		&result.CompanyId,
		&result.Name,
		&result.Address,
		&result.PhoneNumber,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update branch: %v", err)
	}
	return &result, nil
}

func (r *BranchRepo) DeleteBranch(in *company.DeleteBranchRequest) (*company.Message, error) {
	query := `UPDATE branches SET deleted_at = $2 WHERE branch_id = $1`
	_, err := r.db.Exec(query, in.BranchId, time.Now().Unix())
	if err != nil {
		return nil, fmt.Errorf("failed to delete branch: %v", err)
	}
	return &company.Message{Message: "Branch deleted successfully"}, nil
}

func (r *BranchRepo) ListBranches(in *company.ListBranchesRequest) (*company.ListBranchesResponse, error) {
	query := `SELECT branch_id, company_id, name, address, phone_number, created_at, updated_at 
              FROM branches 
              WHERE company_id = $1
              ORDER BY created_at DESC 
              LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, in.CompanyId, in.Limit, (in.Page-1)*in.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %v", err)
	}
	defer rows.Close()

	branches := make([]*company.BranchResponse, 0)
	for rows.Next() {
		var b company.BranchResponse
		err := rows.Scan(
			&b.BranchId,
			&b.CompanyId,
			&b.Name,
			&b.Address,
			&b.PhoneNumber,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		branches = append(branches, &b)
	}
	return &company.ListBranchesResponse{Branches: branches}, nil
}
