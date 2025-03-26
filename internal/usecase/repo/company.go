package repo

import (
	"authentification/internal/generated/company"
	"authentification/internal/usecase"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type CompanyRepo struct {
	db *sqlx.DB
}

func NewCompanyStorage(db *sqlx.DB) usecase.CompanyRepo {
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

	query = `insert into company_account_balance
    (company_id, monthly_fee, last_payment_date, next_due_date)
		values ($1, $2, $3, $4)`

	_, err = tx.Exec(query, result.CompanyId, 250000, time.Now(), time.Now().AddDate(0, 1, 0))
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

// ----------------------------------------- Company Balance ----------------------------------------------

func (r *CompanyRepo) ReplenishmentCompany(in *company.ReplenishmentRequest) (*company.ReplenishmentResponse, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var companyBalance struct {
		ID          uuid.UUID `db:"id"`
		Balance     float64   `db:"balance"`
		MonthlyFee  float64   `db:"monthly_fee"`
		Status      string    `db:"status"`
		NextDueDate time.Time `db:"next_due_date"`
	}

	// Получаем текущий баланс
	err = tx.Get(&companyBalance, `
		SELECT id, balance, monthly_fee, status, next_due_date 
		FROM company_account_balance 
		WHERE company_id = $1 FOR UPDATE`, in.CompanyId)

	if err != nil {
		return nil, errors.New("company not found or DB error")
	}

	// Пополняем баланс
	newBalance := companyBalance.Balance + in.Amount

	// Проверяем, хватит ли денег на оплату
	if companyBalance.Status == "blocked" && newBalance >= companyBalance.MonthlyFee {
		newBalance -= companyBalance.MonthlyFee
		nextDueDate := time.Now().AddDate(0, 1, 0) // +1 месяц

		_, err = tx.Exec(`
			UPDATE company_account_balance 
			SET balance = $1, status = 'active', next_due_date = $2, last_payment_date = NOW()
			WHERE id = $3`, newBalance, nextDueDate, companyBalance.ID)

		if err != nil {
			return nil, err
		}
	} else {
		// Просто обновляем баланс
		_, err = tx.Exec(`
			UPDATE company_account_balance 
			SET balance = $1 WHERE id = $2`, newBalance, companyBalance.ID)

		if err != nil {
			return nil, err
		}
	}

	// Записываем транзакцию
	_, err = tx.Exec(`
		INSERT INTO company_balance_transaction (id, company_id, transaction_date, category, amount)
		VALUES ($1, $2, NOW(), 'deposit', $3)`, uuid.New(), in.CompanyId, in.Amount)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &company.ReplenishmentResponse{
		Success:    true,
		Message:    "Balance replenished successfully",
		NewBalance: newBalance,
	}, nil
}

func (r *CompanyRepo) GetCompanyBalance(in *company.Id) (*company.CompanyBalance, error) {
	var balance float64
	var status string
	var nextDueDate sql.NullTime

	err := r.db.QueryRow(`
		SELECT balance, status, next_due_date 
		FROM company_account_balance 
		WHERE company_id = $1`, in.Id).
		Scan(&balance, &status, &nextDueDate)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("company not found")
		}
		return nil, err
	}

	return &company.CompanyBalance{
		Balance:     balance,
		Status:      status,
		NextDueDate: nextDueDate.Time.Format(time.RFC3339),
	}, nil
}

// 🔹 Получение истории транзакций компании
func (r *CompanyRepo) GetTransactionHistory(in *company.TransactionHistoryRequest) (*company.TransactionHistoryRes, error) {
	rows, err := r.db.Query(`
		SELECT id, transaction_date, category, amount
		FROM company_balance_transaction 
		WHERE company_id = $1 AND transaction_date BETWEEN $2 AND $3 
		ORDER BY transaction_date DESC`, in.CompanyId, in.StartDate, in.EndDate)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*company.Transaction
	for rows.Next() {
		var transaction company.Transaction
		var transactionDate time.Time

		err := rows.Scan(&transaction.TransactionId, &transactionDate, &transaction.Category, &transaction.Amount)
		if err != nil {
			return nil, err
		}
		transaction.TransactionDate = transactionDate.Format(time.RFC3339)
		transactions = append(transactions, &transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &company.TransactionHistoryRes{Transactions: transactions}, nil
}

func BalanceChecker(db *sqlx.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Получаем все компании, у которых наступила дата списания
	query := `
		SELECT id, balance, monthly_fee, next_due_date 
		FROM company_account_balance 
		WHERE next_due_date <= NOW();
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// 2. Обрабатываем каждую компанию
	for rows.Next() {
		var (
			companyID   string
			balance     float64
			monthlyFee  float64
			nextDueDate time.Time
		)

		if err := rows.Scan(&companyID, &balance, &monthlyFee, &nextDueDate); err != nil {
			return err
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if balance >= monthlyFee {
			// 3. Списываем деньги и обновляем дату следующего списания
			updateQuery := `
				UPDATE company_account_balance 
				SET balance = balance - $1, 
				    last_payment_date = NOW(), 
				    next_due_date = NOW() + INTERVAL '1 month',
				    status = 'active'
				WHERE id = $2;
			`
			_, err := tx.ExecContext(ctx, updateQuery, monthlyFee, companyID)
			if err != nil {
				tx.Rollback()
				log.Printf("Ошибка при списании для компании %s: %v", companyID, err)
				continue
			}

			// 4. Записываем транзакцию списания
			insertTransaction := `
				INSERT INTO company_balance_transaction (company_id, transaction_date, category, amount)
				VALUES ($1, NOW(), 'charge', $2);
			`
			_, err = tx.ExecContext(ctx, insertTransaction, companyID, monthlyFee)
			if err != nil {
				tx.Rollback()
				log.Printf("Ошибка при создании транзакции для компании %s: %v", companyID, err)
				continue
			}

			tx.Commit()
		} else {
			// 5. Недостаточно денег — блокируем аккаунт
			blockQuery := `
				UPDATE company_account_balance 
				SET status = 'blocked' 
				WHERE id = $1;
			`
			_, err := tx.ExecContext(ctx, blockQuery, companyID)
			if err != nil {
				tx.Rollback()
				log.Printf("Ошибка при блокировке компании %s: %v", companyID, err)
				continue
			}
			tx.Commit()
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
