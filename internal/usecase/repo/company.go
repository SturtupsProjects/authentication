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
	// Начальные параметры фильтрации
	var filters []string
	var args []interface{}

	// Фильтр по company_id (обязательный)
	filters = append(filters, "company_id = $1")
	args = append(args, in.CompanyId)

	// Фильтр по имени (необязательный)
	if in.Name != "" {
		filters = append(filters, fmt.Sprintf("COALESCE(first_name, '') || ' ' || COALESCE(last_name, '') ILIKE '%%' || $%d || '%%'", len(args)+1))
		args = append(args, in.Name)
	}

	// Построение основного запроса
	mainQuery := fmt.Sprintf(`
		SELECT 
			user_id, 
			COALESCE(first_name, '') || ' ' || COALESCE(last_name, '') AS name, 
			role
		FROM users
		WHERE %s
		ORDER BY name ASC`, strings.Join(filters, " AND "))

	// Добавляем лимит и смещение, если они заданы
	if in.Limit > 0 && in.Page > 0 {
		offset := (in.Page - 1) * in.Limit
		mainQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
		args = append(args, in.Limit, offset)
	}

	// Выполнение основного запроса
	rows, err := r.db.Query(mainQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query company users: %w", err)
	}
	defer rows.Close()

	// Формирование списка пользователей
	users := make([]*company.UserResponse, 0)
	for rows.Next() {
		var u company.UserResponse
		err := rows.Scan(
			&u.UserId,
			&u.Name,
			&u.Role,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, &u)
	}

	// Проверка на ошибки итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user rows: %w", err)
	}

	var totalCount int64
	totalQuery := `SELECT COUNT(*) FROM users WHERE company_id = $1`

	err = r.db.QueryRow(totalQuery, in.CompanyId).Scan(&totalCount)
	if err != nil {
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
