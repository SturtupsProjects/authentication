package repo

import (
	"authentification/internal/generated/company"
	"database/sql"
	"errors"
	"fmt"

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
	query := `INSERT INTO branches (company_id, name, address, phone) 
              VALUES ($1, $2, $3, $4) 
              RETURNING branch_id, company_id, name, address, phone, created_at, updated_at`
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
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}
	return &result, nil
}

func (r *BranchRepo) GetBranch(in *company.GetBranchRequest) (*company.BranchResponse, error) {
	var result company.BranchResponse
	query := `SELECT branch_id, company_id, name, address, phone, created_at, updated_at 
              FROM branches 
              WHERE branch_id = $1 AND company_id = $2 and deleted_at = 0`
	err := r.db.QueryRow(query, in.BranchId, in.CompanyId).Scan(
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
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}
	return &result, nil
}

func (r *BranchRepo) UpdateBranch(in *company.UpdateBranchRequest) (*company.BranchResponse, error) {
	var result company.BranchResponse
	query := `UPDATE branches 
              SET name = COALESCE($1, name), 
                  address = COALESCE($2, address), 
                  phone = COALESCE($3, phone),
                  updated_at = NOW()
              WHERE branch_id = $4 AND company_id = $5 AND deleted_at = 0
              RETURNING branch_id, company_id, name, address, phone, created_at, updated_at`
	err := r.db.QueryRow(query, in.Name, in.Address, in.PhoneNumber, in.BranchId, in.CompanyId).Scan(
		&result.BranchId,
		&result.CompanyId,
		&result.Name,
		&result.Address,
		&result.PhoneNumber,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update branch: %w", err)
	}
	return &result, nil
}

func (r *BranchRepo) DeleteBranch(in *company.DeleteBranchRequest) (*company.Message, error) {
	query := `UPDATE branches 
              SET deleted_at = EXTRACT(EPOCH FROM NOW()) 
              WHERE branch_id = $1 AND company_id = $2`
	result, err := r.db.Exec(query, in.BranchId, in.CompanyId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete branch: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("branch not found or already deleted")
	}

	return &company.Message{Message: "Branch deleted successfully"}, nil
}

func (r *BranchRepo) ListBranches(in *company.ListBranchesRequest) (*company.ListBranchesResponse, error) {
	query := `SELECT branch_id, company_id, name, address, phone, created_at, updated_at 
              FROM branches 
              WHERE company_id = $1 AND deleted_at = 0
              ORDER BY created_at DESC 
              LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, in.CompanyId, in.Limit, (in.Page-1)*in.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var branches []*company.BranchResponse
	for rows.Next() {
		var branch company.BranchResponse
		err := rows.Scan(
			&branch.BranchId,
			&branch.CompanyId,
			&branch.Name,
			&branch.Address,
			&branch.PhoneNumber,
			&branch.CreatedAt,
			&branch.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		branches = append(branches, &branch)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &company.ListBranchesResponse{Branches: branches}, nil
}
