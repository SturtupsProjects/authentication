package repo

import (
	"authentification/internal/generated/company"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CompanyRepo struct {
	db *sqlx.DB
}

func NewCompanyStorage(db *sqlx.DB) *CompanyRepo {
	return &CompanyRepo{db: db}
}

func (r *CompanyRepo) CreateCompany(in *company.CreateCompanyRequest) (*company.CompanyResponse, error) {
	var result company.CompanyResponse
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}

	query := `INSERT INTO company (name, website) 
              VALUES ($1, $2) 
              RETURNING company_id, name, website, COALESCE(logo, ''), created_at, updated_at`
	err = tx.QueryRow(query, in.Name, in.Website).Scan(
		&result.CompanyId,
		&result.Name,
		&result.Website,
		&result.Logo,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to insert company: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &result, nil
}

func (r *CompanyRepo) GetCompany(in *company.GetCompanyRequest) (*company.CompanyResponse, error) {
	var result company.CompanyResponse
	query := `SELECT company_id, name, COALESCE(website, ''), COALESCE(logo, ''), created_at, updated_at 
              FROM company WHERE company_id = $1`
	err := r.db.QueryRow(query, in.CompanyId).Scan(
		&result.CompanyId,
		&result.Name,
		&result.Website,
		&result.Logo,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("company not found")
		}
		return nil, err
	}
	return &result, nil
}

func (r *CompanyRepo) UpdateCompany(in *company.UpdateCompanyRequest) (*company.CompanyResponse, error) {
	var result company.CompanyResponse
	query := `UPDATE company 
              SET name = COALESCE($1, name), 
                  website = COALESCE($2, website), 
                  logo = COALESCE($3, logo),
                  updated_at = NOW()
              WHERE company_id = $4 
              RETURNING company_id, name, website, logo, created_at, updated_at`

	err := r.db.QueryRow(query, in.Name, in.Website, in.Logo, in.CompanyId).Scan(
		&result.CompanyId,
		&result.Name,
		&result.Website,
		&result.Logo,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *CompanyRepo) DeleteCompany(in *company.DeleteCompanyRequest) (*company.Message, error) {
	query := `DELETE FROM company WHERE company_id = $1`
	_, err := r.db.Exec(query, in.CompanyId)
	if err != nil {
		return nil, err
	}
	return &company.Message{Message: "Deleted company"}, nil
}

func (r *CompanyRepo) ListCompanies(in *company.ListCompaniesRequest) (*company.ListCompaniesResponse, error) {
	query := `SELECT company_id, name, COALESCE(website,''), COALESCE(logo, ''), created_at, updated_at 
              FROM company 
              ORDER BY created_at DESC 
              LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, in.Limit, in.Page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	companies := make([]*company.CompanyResponse, 0)
	for rows.Next() {
		var c company.CompanyResponse
		err := rows.Scan(
			&c.CompanyId,
			&c.Name,
			&c.Website,
			&c.Logo,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		companies = append(companies, &c)
	}

	return &company.ListCompaniesResponse{Companies: companies}, nil
}

func (r *CompanyRepo) ListCompanyUsers(in *company.ListCompanyUsersRequest) (*company.ListCompanyUsersResponse, error) {
	query := `SELECT user_id, CONCAT(first_name, ' ', last_name) AS name, role 
              FROM users 
              WHERE company_id = $1 
              LIMIT $2 OFFSET $3`

	offset := (in.Page - 1) * in.Limit
	rows, err := r.db.Query(query, in.CompanyId, in.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*company.UserResponse, 0)
	for rows.Next() {
		var u company.UserResponse
		err := rows.Scan(
			&u.UserId,
			&u.Name,
			&u.Role,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return &company.ListCompanyUsersResponse{Users: users}, nil
}
func (r *CompanyRepo) CreateUserToCompany(in *company.CreateUserToCompanyRequest) (*company.Id, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var companyID string
	err = tx.QueryRow(`SELECT company_id FROM company WHERE company_id = $1`, in.CompanyId).Scan(&companyID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if companyID == "" {
		err = tx.QueryRow(`INSERT INTO company (name, company_id) VALUES ($1, $2) RETURNING company_id`,
			in.FirstName, in.CompanyId).Scan(&companyID)
		if err != nil {
			return nil, err
		}
	}

	var userID string
	err = tx.QueryRow(`INSERT INTO users (first_name, last_name, email,phone_number, password, role, company_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING user_id`,
		in.FirstName, in.LastName, in.Email, in.PhoneNumber, in.Password, in.Role, companyID).Scan(&userID)
	if err != nil {
		return nil, err
	}

	return &company.Id{Id: userID}, nil
}
