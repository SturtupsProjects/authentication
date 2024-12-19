package repo

import (
	"fmt"
	"strings"

	pb "authentification/pkg/generated/user"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) AddAdmin(response *pb.MessageResponse) (*pb.MessageResponse, error) {
	_, err := u.db.Exec(
		`INSERT INTO users(first_name, last_name, email, phone_number, password, role)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		"admin", "admin", "admin@admin.com", "admin", response.Message, "admin",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add admin: %w", err)
	}
	return &pb.MessageResponse{Message: "Admin added successfully"}, nil
}

func (u *UserRepo) CreateUser(in *pb.UserRequest) (*pb.UserResponse, error) {
	var user pb.UserResponse
	query := `
		INSERT INTO users (first_name, last_name, email, phone_number, password, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING user_id, first_name, last_name, email, phone_number, role, created_at
	`
	err := u.db.QueryRowx(query, in.FirstName, in.LastName, in.Email, in.PhoneNumber, in.Password, in.Role).
		Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &user, nil
}

func (u *UserRepo) GetUser(in *pb.UserIDRequest) (*pb.UserResponse, error) {
	var user User
	query := `SELECT user_id, first_name, last_name, email, phone_number, role, created_at FROM users WHERE user_id = $1`
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
}

func (u *UserRepo) GetListUser(in *pb.FilterUserRequest) (*pb.UserListResponse, error) {
	var users []User
	var queryBuilder strings.Builder
	var args []interface{}
	argCounter := 1

	queryBuilder.WriteString(`SELECT user_id, first_name, last_name, email, phone_number, role, created_at FROM users WHERE 1=1`)

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

	// Remove the trailing comma and space
	query := queryBuilder.String()
	query = strings.TrimSuffix(query, ", ")

	// Add WHERE clause
	query += fmt.Sprintf(" WHERE user_id = $%d RETURNING user_id, first_name, last_name, email, phone_number, role, created_at", argIndex)
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
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}
func (u *UserRepo) LogIn(in *pb.LogInRequest) (*pb.LogInResponse, string, error) {
	var loginResp pb.LogInResponse
	var password string
	query := `SELECT user_id, first_name, phone_number, role, password FROM users WHERE phone_number = $1`
	err := u.db.QueryRowx(query, in.PhoneNumber).Scan(
		&loginResp.UserId,
		&loginResp.FirstName,
		&loginResp.PhoneNumber,
		&loginResp.Role,
		&password,
	)
	if err != nil {
		return nil, "", fmt.Errorf("login failed: %w", err)
	}
	return &loginResp, password, nil
}
