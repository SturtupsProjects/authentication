package repo

import (
	"database/sql"
	"fmt"
	"strings"

	pb "authentification/internal/generated/user"
	"authentification/internal/usecase"

	"github.com/jmoiron/sqlx"
)

type workersRepo struct {
	db *sqlx.DB
}

func NewWorkersRepo(db *sqlx.DB) usecase.WorkersRepo {
	return &workersRepo{db: db}
}

func (w *workersRepo) CreateSalary(in *pb.SalaryRequest) (*pb.SalaryResponse, error) {
	query := `INSERT INTO staff_salary (user_id, currency_code, salary_amount, salary_date, company_id)
	          VALUES ($1, $2, $3, $4, $5) RETURNING salary_id, created_at, updated_at`

	var resp pb.SalaryResponse
	err := w.db.QueryRow(query, in.UserId, in.CurrencyCode, in.Amount, in.SalaryDate, in.CompanyId).
		Scan(&resp.SalaryId, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		return nil, err
	}

	resp.UserId = in.UserId
	resp.CurrencyCode = in.CurrencyCode
	resp.Amount = in.Amount
	resp.SalaryDate = in.SalaryDate

	return &resp, nil
}

func (w *workersRepo) UpdateSalary(in *pb.SalaryUpdate) (*pb.SalaryResponse, error) {
	query := `UPDATE staff_salary 
	          SET currency_code = $1, salary_amount = $2, salary_date = $3, updated_at = NOW()
	          WHERE salary_id = $4 AND company_id = $5 RETURNING salary_id, updated_at`

	var resp pb.SalaryResponse
	err := w.db.QueryRow(query, in.CurrencyCode, in.Amount, in.SalaryDate, in.SalaryId, in.CompanyId).
		Scan(&resp.SalaryId, &resp.UpdatedAt)
	if err != nil {
		return nil, err
	}

	resp.CurrencyCode = in.CurrencyCode
	resp.Amount = in.Amount
	resp.SalaryDate = in.SalaryDate

	return &resp, nil
}

func (w *workersRepo) GetSalaryByID(in *pb.ID) (*pb.SalaryResponse, error) {
	query := `
		SELECT u.user_id, u.first_name, u.last_name, u.phone_number, u.company_id,
		       s.salary_id, s.currency_code, s.salary_amount, s.salary_date,
		       to_char(s.created_at, 'YYYY-MM-DD') as created_at,
		       to_char(s.updated_at, 'YYYY-MM-DD') as updated_at
		FROM staff_salary s
		JOIN users u ON u.user_id = s.user_id
		WHERE s.salary_id = $1 AND u.company_id = $2
	`

	var resp pb.SalaryResponse
	var salaryId, currencyCode, salaryDate, createdAt, updatedAt sql.NullString
	var amount sql.NullFloat64

	err := w.db.QueryRow(query, in.Id, in.CompanyId).
		Scan(
			&resp.UserId,
			&resp.FirstName,
			&resp.LastName,
			&resp.PhoneNumber,
			&resp.CompanyId,
			&salaryId,
			&currencyCode,
			&amount,
			&salaryDate,
			&createdAt,
			&updatedAt,
		)
	if err != nil {
		return nil, err
	}

	if salaryId.Valid {
		resp.SalaryId = salaryId.String
	}
	if currencyCode.Valid {
		resp.CurrencyCode = currencyCode.String
	}
	if amount.Valid {
		resp.Amount = amount.Float64
	}
	if salaryDate.Valid {
		resp.SalaryDate = salaryDate.String
	}
	if createdAt.Valid {
		resp.CreatedAt = createdAt.String
	}
	if updatedAt.Valid {
		resp.UpdatedAt = updatedAt.String
	}

	return &resp, nil
}

func (w *workersRepo) ListSalaries(in *pb.GetSalaryRequest) (*pb.GetSalaryList, error) {
	var salaries []*pb.SalaryResponse
	var args []interface{}
	argIndex := 1

	// Фильтрация по компании
	conditions := []string{"u.company_id = $1"}
	args = append(args, in.CompanyId)
	argIndex++

	allowedSortFields := map[string]string{
		"salary_date":   "s.salary_date",
		"created_at":    "created_at",
		"salary_amount": "s.salary_amount",
	}
	sortField, ok := allowedSortFields[strings.ToLower(in.SortField)]
	if !ok {
		sortField = "s.salary_date"
	}
	order := strings.ToUpper(in.Order)
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	baseQuery := fmt.Sprintf(`
		SELECT DISTINCT
			u.user_id, u.first_name, u.last_name, u.phone_number, u.company_id,
			s.salary_id, s.currency_code, s.salary_amount, s.salary_date,
			to_char(s.created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at,
			to_char(s.updated_at, 'YYYY-MM-DD HH24:MI:SS') as updated_at
		FROM users u
		LEFT JOIN staff_salary s ON u.user_id = s.user_id
		WHERE %s
		ORDER BY %s %s`,
		strings.Join(conditions, " AND "),
		sortField, order)

	if in.Limit > 0 && in.Page > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, in.Limit, (in.Page-1)*in.Limit)
		argIndex += 2
	}

	rows, err := w.db.Query(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("ListSalaries: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var salary pb.SalaryResponse
		var salaryId, currencyCode, salaryDate, createdAt, updatedAt sql.NullString
		var amount sql.NullFloat64

		if err := rows.Scan(
			&salary.UserId,
			&salary.FirstName,
			&salary.LastName,
			&salary.PhoneNumber,
			&salary.CompanyId,
			&salaryId,
			&currencyCode,
			&amount,
			&salaryDate,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("ListSalaries scan: %w", err)
		}

		if salaryId.Valid {
			salary.SalaryId = salaryId.String
		}
		if currencyCode.Valid {
			salary.CurrencyCode = currencyCode.String
		}
		if amount.Valid {
			salary.Amount = amount.Float64
		}
		if salaryDate.Valid {
			salary.SalaryDate = salaryDate.String
		}
		salary.CreatedAt = createdAt.String
		salary.UpdatedAt = updatedAt.String

		salaries = append(salaries, &salary)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListSalaries rows error: %w", err)
	}

	return &pb.GetSalaryList{
		Salaries: salaries,
	}, nil
}

func (w *workersRepo) CreateAdjustment(in *pb.AdjustmentRequest) (*pb.AdjustmentResponse, error) {
	query := `INSERT INTO salary_adjustments (user_id, adjustment_type, currency_code, adjustment_amount, adjustment_date, company_id)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING adjustment_id, created_at, updated_at`

	var resp pb.AdjustmentResponse
	err := w.db.QueryRow(query, in.UserId, in.AdjustmentType, in.CurrencyCode, in.Amount, in.AdjustmentDate, in.CompanyId).
		Scan(&resp.AdjustmentId, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		return nil, err
	}

	resp.UserId = in.UserId
	resp.AdjustmentType = in.AdjustmentType
	resp.CurrencyCode = in.CurrencyCode
	resp.Amount = in.Amount
	resp.AdjustmentDate = in.AdjustmentDate
	resp.CompanyId = in.CompanyId

	return &resp, nil
}

func (w *workersRepo) UpdateAdjustment(in *pb.AdjustmentUpdate) (*pb.AdjustmentResponse, error) {
	query := `
		UPDATE salary_adjustments
		SET adjustment_type = $1,
		    currency_code = $2,
		    adjustment_amount = $3,
		    adjustment_date = $4,
		    updated_at = NOW()
		WHERE adjustment_id = $5 AND company_id = $6
		RETURNING adjustment_id, user_id, adjustment_type, currency_code, adjustment_amount, 
		          adjustment_date, is_active, created_at, updated_at, company_id
	`

	var resp pb.AdjustmentResponse
	err := w.db.QueryRow(query,
		in.AdjustmentType, in.CurrencyCode, in.Amount, in.AdjustmentDate,
		in.AdjustmentId, in.CompanyId,
	).Scan(
		&resp.AdjustmentId, &resp.UserId, &resp.AdjustmentType, &resp.CurrencyCode,
		&resp.Amount, &resp.AdjustmentDate, &resp.IsActive, &resp.CreatedAt,
		&resp.UpdatedAt, &resp.CompanyId,
	)
	if err != nil {
		return nil, fmt.Errorf("UpdateAdjustment error: %w", err)
	}

	return &resp, nil
}

func (w *workersRepo) CloseAdjustment(in *pb.ID) (*pb.AdjustmentResponse, error) {
	updateQuery := `
		UPDATE salary_adjustments
		SET is_active = FALSE,
		    updated_at = NOW()
		WHERE adjustment_id = $1 AND company_id = $2
		RETURNING adjustment_id
	`
	var adjustmentID string
	err := w.db.QueryRow(updateQuery, in.Id, in.CompanyId).Scan(&adjustmentID)
	if err != nil {
		return nil, fmt.Errorf("CloseAdjustment update: %w", err)
	}

	query := `
		SELECT u.user_id, u.first_name, u.last_name, u.phone_number, u.company_id,
		       s.adjustment_id, s.adjustment_type, s.currency_code, s.adjustment_amount, s.adjustment_date, s.is_active,
		       to_char(s.created_at, 'YYYY-MM-DD HH24:MI:SS'),
		       to_char(s.updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM users u
		LEFT JOIN salary_adjustments s ON u.user_id = s.user_id AND s.adjustment_id = $1
		WHERE u.company_id = $2
		LIMIT 1
	`

	var respAdj pb.AdjustmentResponse
	var adjID, adjType, currCode, adjDate, createdAt, updatedAt sql.NullString
	var adjAmount sql.NullFloat64
	var isActive sql.NullString

	err = w.db.QueryRow(query, in.Id, in.CompanyId).Scan(
		&respAdj.UserId,
		&respAdj.FirstName,
		&respAdj.LastName,
		&respAdj.PhoneNumber,
		&respAdj.CompanyId,
		&adjID,
		&adjType,
		&currCode,
		&adjAmount,
		&adjDate,
		&isActive,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("CloseAdjustment select: %w", err)
	}

	if adjID.Valid {
		respAdj.AdjustmentId = adjID.String
	}
	if adjType.Valid {
		respAdj.AdjustmentType = adjType.String
	}
	if currCode.Valid {
		respAdj.CurrencyCode = currCode.String
	}
	if adjAmount.Valid {
		respAdj.Amount = adjAmount.Float64
	}
	if adjDate.Valid {
		respAdj.AdjustmentDate = adjDate.String
	}
	if isActive.Valid {
		respAdj.IsActive = isActive.String
	}
	if createdAt.Valid {
		respAdj.CreatedAt = createdAt.String
	}
	if updatedAt.Valid {
		respAdj.UpdatedAt = updatedAt.String
	}

	return &respAdj, nil
}

func (w *workersRepo) GetAdjustmentByID(in *pb.ID) (*pb.AdjustmentResponse, error) {
	query := `
		SELECT u.user_id, u.first_name, u.last_name, u.phone_number, u.company_id,
		       a.adjustment_id, a.adjustment_type, a.currency_code, a.adjustment_amount, a.adjustment_date, a.is_active,
		       to_char(a.created_at, 'YYYY-MM-DD HH24:MI:SS'), 
		       to_char(a.updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM users u
		INNER JOIN salary_adjustments a ON u.user_id = a.user_id
		WHERE a.adjustment_id = $1 AND u.company_id = $2
		LIMIT 1
	`

	var resp pb.AdjustmentResponse
	var adjustmentId, adjustmentType, currencyCode, adjustmentDate, createdAt, updatedAt sql.NullString
	var adjustmentAmount sql.NullFloat64
	var isActive sql.NullString

	err := w.db.QueryRow(query, in.Id, in.CompanyId).
		Scan(
			&resp.UserId,
			&resp.FirstName,
			&resp.LastName,
			&resp.PhoneNumber,
			&resp.CompanyId,
			&adjustmentId,
			&adjustmentType,
			&currencyCode,
			&adjustmentAmount,
			&adjustmentDate,
			&isActive,
			&createdAt,
			&updatedAt,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("adjustment notfount")
		}
		return nil, fmt.Errorf("GetAdjustmentByID: %w", err)
	}

	if adjustmentId.Valid {
		resp.AdjustmentId = adjustmentId.String
	}
	if adjustmentType.Valid {
		resp.AdjustmentType = adjustmentType.String
	}
	if currencyCode.Valid {
		resp.CurrencyCode = currencyCode.String
	}
	if adjustmentAmount.Valid {
		resp.Amount = adjustmentAmount.Float64
	}
	if adjustmentDate.Valid {
		resp.AdjustmentDate = adjustmentDate.String
	}
	if isActive.Valid {
		resp.IsActive = isActive.String
	}
	if createdAt.Valid {
		resp.CreatedAt = createdAt.String
	}
	if updatedAt.Valid {
		resp.UpdatedAt = updatedAt.String
	}

	return &resp, nil
}

func (w *workersRepo) ListAdjustments(in *pb.GetAdjustmentRequest) (*pb.AdjustmentList, error) {
	conditions := []string{"u.company_id = $1"}
	args := []interface{}{in.CompanyId}
	argIdx := 2

	if in.AdjustmentType != "" {
		conditions = append(conditions, fmt.Sprintf("a.adjustment_type = $%d", argIdx))
		args = append(args, in.AdjustmentType)
		argIdx++
	}

	if in.IsActive != "" {
		conditions = append(conditions, fmt.Sprintf("a.is_active = $%d", argIdx))
		args = append(args, in.IsActive)
		argIdx++
	}

	allowedSortFields := map[string]string{
		"adjustment_date": "a.adjustment_date",
		"created_at":      "a.created_at",
		"amount":          "a.adjustment_amount",
	}

	orderByClause := ""
	if sortField, sortExists := allowedSortFields[strings.ToLower(in.SortField)]; sortExists {
		order := strings.ToUpper(in.Order)
		if order != "ASC" && order != "DESC" {
			order = "ASC"
		}
		orderByClause = fmt.Sprintf("ORDER BY %s %s", sortField, order)
	}

	limitOffsetClause := ""
	if in.Limit > 0 && in.Page > 0 {
		offset := (in.Page - 1) * in.Limit
		limitOffsetClause = fmt.Sprintf("LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
		args = append(args, in.Limit, offset)
		argIdx += 2
	}

	query := fmt.Sprintf(`
		SELECT u.user_id, u.first_name, u.last_name, u.phone_number, u.company_id,
		       a.adjustment_id, a.adjustment_type, a.currency_code, a.adjustment_amount AS amount, 
		       a.adjustment_date, a.is_active,
		       to_char(a.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at,
		       to_char(a.updated_at, 'YYYY-MM-DD HH24:MI:SS') AS updated_at
		FROM users u
		LEFT JOIN salary_adjustments a ON u.user_id = a.user_id
		WHERE %s
		%s
		%s
	`, strings.Join(conditions, " AND "), orderByClause, limitOffsetClause)

	rows, err := w.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ListAdjustments: %w", err)
	}
	defer rows.Close()

	var adjustments []*pb.AdjustmentResponse
	for rows.Next() {
		var adjustment pb.AdjustmentResponse

		var adjustmentId, adjustmentType, currencyCode, adjustmentDate, createdAt, updatedAt sql.NullString
		var adjustmentAmount sql.NullFloat64
		var isActive sql.NullString

		if err := rows.Scan(
			&adjustment.UserId,
			&adjustment.FirstName,
			&adjustment.LastName,
			&adjustment.PhoneNumber,
			&adjustment.CompanyId,
			&adjustmentId,
			&adjustmentType,
			&currencyCode,
			&adjustmentAmount,
			&adjustmentDate,
			&isActive,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("ListAdjustments scan: %w", err)
		}

		if adjustmentId.Valid {
			adjustment.AdjustmentId = adjustmentId.String
		}
		if adjustmentType.Valid {
			adjustment.AdjustmentType = adjustmentType.String
		}
		if currencyCode.Valid {
			adjustment.CurrencyCode = currencyCode.String
		}
		if adjustmentAmount.Valid {
			adjustment.Amount = adjustmentAmount.Float64
		}
		if adjustmentDate.Valid {
			adjustment.AdjustmentDate = adjustmentDate.String
		}
		if isActive.Valid {
			adjustment.IsActive = isActive.String
		}
		if createdAt.Valid {
			adjustment.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			adjustment.UpdatedAt = updatedAt.String
		}

		adjustments = append(adjustments, &adjustment)
	}

	return &pb.AdjustmentList{Adjustments: adjustments}, nil
}

func (w *workersRepo) GetWorkerAllInfo(in *pb.ID) (*pb.WorkerAllInfo, error) {
	// Получаем данные пользователя и зарплату одним запросом.
	query := `
		SELECT u.user_id, u.first_name, u.last_name, u.phone_number, u.email, u.role, u.company_id,
		       s.salary_id, s.currency_code, s.salary_amount, s.salary_date,
		       to_char(s.created_at, 'YYYY-MM-DD'), 
		       to_char(s.updated_at, 'YYYY-MM-DD')
		FROM users u
		LEFT JOIN staff_salary s ON u.user_id = s.user_id
		WHERE u.user_id = $1 AND u.company_id = $2
		LIMIT 1
	`

	var resp pb.WorkerAllInfo
	var salaryId, currencyCode, salaryDate, sCreatedAt, sUpdatedAt sql.NullString
	var amount sql.NullFloat64

	err := w.db.QueryRow(query, in.Id, in.CompanyId).
		Scan(
			&resp.UserId,
			&resp.FirstName,
			&resp.LastName,
			&resp.PhoneNumber,
			&resp.Email,
			&resp.Role,
			&resp.CompanyId,
			&salaryId,
			&currencyCode,
			&amount,
			&salaryDate,
			&sCreatedAt,
			&sUpdatedAt,
		)
	if err != nil {
		return nil, fmt.Errorf("GetWorkerAllInfo user/salary: %w", err)
	}

	// Заполняем данные зарплаты, если они есть.
	if salaryId.Valid {
		resp.Salary.SalaryId = salaryId.String
	}
	if currencyCode.Valid {
		resp.Salary.CurrencyCode = currencyCode.String
	}
	if amount.Valid {
		resp.Salary.Amount = amount.Float64
	}
	if salaryDate.Valid {
		resp.Salary.SalaryDate = salaryDate.String
	}
	if sCreatedAt.Valid {
		resp.CreatedAt = sCreatedAt.String
	}
	if sUpdatedAt.Valid {
		resp.UpdatedAt = sUpdatedAt.String
	}

	// Получаем корректировки для пользователя.
	adjQuery := `
		SELECT adjustment_id, adjustment_type, currency_code, adjustment_amount, adjustment_date,
		       is_active, to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'), to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM salary_adjustments
		WHERE user_id = $1
		ORDER BY adjustment_date DESC
	`
	rows, err := w.db.Query(adjQuery, in.Id)
	if err != nil {
		return nil, fmt.Errorf("GetWorkerAllInfo adjustments: %w", err)
	}
	defer rows.Close()

	var adjustments []*pb.Adjustment
	for rows.Next() {
		var adjustment pb.Adjustment
		var adjustmentId, adjustmentType, currencyCode, adjustmentDate, createdAt, updatedAt sql.NullString
		var adjustmentAmount sql.NullFloat64
		var isActive sql.NullString

		if err = rows.Scan(
			&adjustmentId,
			&adjustmentType,
			&currencyCode,
			&adjustmentAmount,
			&adjustmentDate,
			&isActive,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("GetWorkerAllInfo adjustments scan: %w", err)
		}

		if adjustmentId.Valid {
			adjustment.AdjustmentId = adjustmentId.String
		}
		if adjustmentType.Valid {
			adjustment.AdjustmentType = adjustmentType.String
		}
		if currencyCode.Valid {
			adjustment.CurrencyCode = currencyCode.String
		}
		if adjustmentAmount.Valid {
			adjustment.Amount = adjustmentAmount.Float64
		}
		if adjustmentDate.Valid {
			adjustment.AdjustmentDate = adjustmentDate.String
		}
		if isActive.Valid {
			adjustment.IsActive = isActive.String
		}
		if createdAt.Valid {
			adjustment.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			adjustment.UpdatedAt = updatedAt.String
		}

		adjustments = append(adjustments, &adjustment)
	}
	resp.Adjustments = adjustments

	return &resp, nil
}
