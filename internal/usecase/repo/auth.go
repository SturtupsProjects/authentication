package repo

import (
	"fmt"
	"strings"

	pb "authentification/internal/generated/user"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) AddAdmin(response *pb.MessageResponse) (*pb.MessageResponse, error) {

	var companyID string
	err := u.db.QueryRowx("INSERT INTO company (name) VALUES ($1) RETURNING company_id", "admin").Scan(&companyID)
	if err != nil {
		return nil, err
	}

	_, err = u.db.Exec(
		`INSERT INTO users(first_name, last_name, email, phone_number, password, role, company_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		"admin", "admin", "admin@admin.com", "admin", response.Message, "admin", companyID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &pb.MessageResponse{Message: "Admin added successfully"}, nil
}

func (u *UserRepo) CreateUser(in *pb.UserRequest) (*pb.UserResponse, error) {
	var user pb.UserResponse
	tx, err := u.db.Beginx() // Use Beginx() for sqlx
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
	err = tx.QueryRowx("INSERT INTO company (name) VALUES ($1) RETURNING company_id", in.FirstName).Scan(&companyID)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO users (first_name, last_name, email, phone_number, password, role, company_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING user_id, first_name, last_name, email, phone_number, role, created_at
	`
	err = tx.QueryRowx(query, in.FirstName, in.LastName, in.Email, in.PhoneNumber, in.Password, in.Role, companyID).
		Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &user, nil
}

func (u *UserRepo) GetUser(in *pb.UserIDRequest) (*pb.UserResponse, error) {
	var user User
	query := `SELECT user_id, first_name, last_name, email, phone_number, role, created_at, company_id FROM users WHERE user_id = $1`
	err := u.db.Get(&user, query, in.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &pb.UserResponse{UserId: user.UserID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		CompanyId:   user.CompanyID,
	}, nil
}

type User struct {
	UserID      string `db:"user_id"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	Email       string `db:"email"`
	PhoneNumber string `db:"phone_number"`
	Role        string `db:"role"`
	CreatedAt   string `db:"created_at"`
	CompanyID   string `db:"company_id"`
}

func (u *UserRepo) GetListUser(in *pb.FilterUserRequest) (*pb.UserListResponse, error) {
	var users []User
	var queryBuilder strings.Builder
	var args []interface{}
	argCounter := 1

	queryBuilder.WriteString(`SELECT user_id, first_name, last_name, email, phone_number, role, created_at, company_id FROM users WHERE company_id IS NOT NULL `)

	if in.FirstName != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND first_name ILIKE $%d", argCounter))
		args = append(args, "%"+in.FirstName+"%")
		argCounter++
	}

	if in.LastName != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND last_name ILIKE $%d", argCounter))
		args = append(args, "%"+in.LastName+"%")
		argCounter++
	}

	if in.Role != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND role = $%d", argCounter))
		args = append(args, in.Role)
		argCounter++
	}

	queryBuilder.WriteString(" ORDER BY created_at DESC")

	query := queryBuilder.String()
	err := u.db.Select(&users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Convert to pb.UserResponse
	var pbUsers []*pb.UserResponse
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.UserResponse{
			UserId:      user.UserID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt,
			CompanyId:   user.CompanyID,
		})
	}

	return &pb.UserListResponse{Users: pbUsers}, nil
}

func (u *UserRepo) DeleteUser(in *pb.UserIDRequest) (*pb.MessageResponse, error) {
	query := `DELETE FROM users WHERE user_id = $1`
	res, err := u.db.Exec(query, in.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}
	rows, _ := res.RowsAffected()
	return &pb.MessageResponse{Message: fmt.Sprintf("Deleted %d user(s)", rows)}, nil
}

func (u *UserRepo) UpdateUser(in *pb.UserRequest) (*pb.UserResponse, error) {
	if in.UserId == "" {
		return nil, fmt.Errorf("user ID is required for updating a user")
	}

	var (
		queryBuilder strings.Builder
		args         []interface{}
		argIndex     = 1
	)

	queryBuilder.WriteString("UPDATE users SET ")

	if in.FirstName != "" {
		queryBuilder.WriteString(fmt.Sprintf("first_name = $%d, ", argIndex))
		args = append(args, in.FirstName)
		argIndex++
	}

	if in.LastName != "" {
		queryBuilder.WriteString(fmt.Sprintf("last_name = $%d, ", argIndex))
		args = append(args, in.LastName)
		argIndex++
	}

	if in.Email != "" {
		queryBuilder.WriteString(fmt.Sprintf("email = $%d, ", argIndex))
		args = append(args, in.Email)
		argIndex++
	}

	if in.Password != "" {
		queryBuilder.WriteString(fmt.Sprintf("password = $%d, ", argIndex))
		args = append(args, in.Password)
		argIndex++
	}

	if in.PhoneNumber != "" {
		queryBuilder.WriteString(fmt.Sprintf("phone_number = $%d, ", argIndex))
		args = append(args, in.PhoneNumber)
		argIndex++
	}

	if in.Role != "" {
		queryBuilder.WriteString(fmt.Sprintf("role = $%d, ", argIndex))
		args = append(args, in.Role)
		argIndex++
	}

	if in.CompanyId != "" {
		queryBuilder.WriteString(fmt.Sprintf("company_id = $%d, ", argIndex))
		args = append(args, in.CompanyId)
		argIndex++
	}

	// Remove the trailing comma and space
	query := queryBuilder.String()
	query = strings.TrimSuffix(query, ", ")

	// Add WHERE clause
	query += fmt.Sprintf(" WHERE user_id = $%d RETURNING user_id, first_name, last_name, email, phone_number, role, created_at, company_id", argIndex)
	args = append(args, in.UserId)

	// Execute the query
	var user pb.UserResponse
	err := u.db.QueryRowx(query, args...).Scan(
		&user.UserId,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PhoneNumber,
		&user.Role,
		&user.CreatedAt,
		&user.CompanyId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}
func (u *UserRepo) LogIn(in *pb.LogInRequest) (*pb.LogInResponse, string, string, error) {
	var loginResp pb.LogInResponse
	var password string
	var companyID string
	query := `SELECT user_id, first_name, phone_number, role, password, company_id FROM users WHERE phone_number = $1`
	err := u.db.QueryRowx(query, in.PhoneNumber).Scan(
		&loginResp.UserId,
		&loginResp.FirstName,
		&loginResp.PhoneNumber,
		&loginResp.Role,
		&password,
		&companyID,
	)
	if err != nil {
		return nil, "", "", fmt.Errorf("login failed: %w", err)
	}
	return &loginResp, password, companyID, nil
}
