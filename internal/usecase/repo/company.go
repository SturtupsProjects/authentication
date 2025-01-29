package repo

import (
	"authentification/internal/generated/company"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

	var filters []string
	var args []interface{}

	filters = append(filters, "company_id = $1")
	args = append(args, in.CompanyId)

	if in.Name != "" {
		filters = append(filters, fmt.Sprintf("CONCAT(COALESCE(first_name, ''), ' ', COALESCE(last_name, '')) ILIKE $%d", len(args)+1))
		args = append(args, "%"+in.Name+"%")
	}

	baseQuery := fmt.Sprintf(`
		SELECT 
			user_id, 
			CONCAT(COALESCE(first_name, ''), ' ', COALESCE(last_name, '')) AS name, 
			role
		FROM users
		WHERE %s`, strings.Join(filters, " AND "))

	baseQuery += " ORDER BY name ASC"

	if in.Limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, in.Limit)
	}
	if in.Page > 0 && in.Limit > 0 {
		offset := (in.Page - 1) * in.Limit
		baseQuery += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, offset)
	}

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query company users: %w", err)
	}
	defer rows.Close()

	users := make([]*company.UserResponse, 0)
	for rows.Next() {
		var user company.UserResponse
		if err := rows.Scan(&user.UserId, &user.Name, &user.Role); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user rows: %w", err)
	}

	var totalCount int64
	totalQuery := fmt.Sprintf(`SELECT COUNT(*) FROM users WHERE %s`, strings.Join(filters, " AND "))
	if err := r.db.QueryRow(totalQuery, args[:len(args)-2]...).Scan(&totalCount); err != nil {
		return nil, fmt.Errorf("failed to count company users: %w", err)
	}

	return &company.ListCompanyUsersResponse{
		Users:      users,
		TotalCount: totalCount,
	}, nil
}

func (r *CompanyRepo) CreateUserToCompany(in *company.CreateUserToCompanyRequest) (*company.Id, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure proper rollback and commit logic with deferred function
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

	// Check if the company exists
	var companyID string
	err = tx.QueryRow(`SELECT company_id FROM company WHERE company_id = $1`, in.CompanyId).Scan(&companyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("company not found")
		}
		return nil, fmt.Errorf("failed to check company existence: %w", err)
	}

	// Insert user into the database
	var userID string
	err = tx.QueryRow(
		`INSERT INTO users (first_name, last_name, email, phone_number, password, role, company_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING user_id`,
		in.FirstName, in.LastName, in.Email, in.PhoneNumber, in.Password, in.Role, companyID,
	).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	// Return created user ID
	return &company.Id{Id: userID}, nil
}

func (r *CompanyRepo) CreateBalance(req *company.CompanyBalanceRequest) (*company.CompanyBalanceResponse, error) {
	var result company.CompanyBalanceResponse
	query := `INSERT INTO company_balance (company_id, amount) 
              VALUES ($1, $2) 
              RETURNING company_id, amount`
	err := r.db.QueryRow(query, req.CompanyId, req.Balance).Scan(&result.CompanyId, &result.Balance)
	if err != nil {
		return nil, fmt.Errorf("failed to create balance: %w", err)
	}
	return &result, nil
}

func (r *CompanyRepo) GetBalance(req *company.Id) (*company.CompanyBalanceResponse, error) {
	var result company.CompanyBalanceResponse
	query := `SELECT company_id, amount FROM company_balance WHERE company_id = $1`
	err := r.db.QueryRow(query, req.Id).Scan(&result.CompanyId, &result.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("balance not found")
		}
		return nil, err
	}
	return &result, nil
}

func (r *CompanyRepo) UpdateBalance(req *company.CompanyBalanceRequest) (*company.CompanyBalanceResponse, error) {
	var result company.CompanyBalanceResponse
	query := `UPDATE company_balance SET amount = $1 WHERE company_id = $2 RETURNING company_id, amount`
	err := r.db.QueryRow(query, req.Balance, req.CompanyId).Scan(&result.CompanyId, &result.Balance)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}
	return &result, nil
}

func (r *CompanyRepo) DeleteBalance(req *company.Id) (*company.Message, error) {
	query := `UPDATE company_balance SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE company_id = $1`
	result, err := r.db.Exec(query, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete balance: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("balance not found or already deleted")
	}

	return &company.Message{Message: "Balance deleted successfully"}, nil
}

func (r *CompanyRepo) ListBalances(req *company.FilterCompanyBalanceRequest) (*company.CompanyBalanceListResponse, error) {
	var balances []*company.CompanyBalanceResponse
	query := `SELECT company_id, amount FROM company_balance ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, req.Limit, (req.Page-1)*req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list balances: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var balance company.CompanyBalanceResponse
		if err := rows.Scan(&balance.CompanyId, &balance.Balance); err != nil {
			return nil, fmt.Errorf("failed to scan balance row: %w", err)
		}
		balances = append(balances, &balance)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over balance rows: %w", err)
	}

	return &company.CompanyBalanceListResponse{Users: balances}, nil
}
