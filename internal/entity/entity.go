package entity

import "time"

// ---------------------------------- Message ---------------------------------------------

type Message struct {
	Message string `json:"message"`
}

// -------- User structs for Repo -----------------------------------------

type User struct {
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Email       string `json:"email" db:"email"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Password    string `json:"password"`
	Role        string `json:"role" db:"role"`
}

type UserRequest struct {
	UserID      string    `json:"user_id" db:"user_id"` // Omitted for Create
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Email       string    `json:"email" db:"email"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Role        string    `json:"role" db:"role"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type UserUpdate struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
}

type UserID struct {
	ID string `json:"id"`
}

type FilterUser struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Role      string `json:"role,omitempty"`
}

type UserList struct {
	Users []UserRequest `json:"users"`
}

type AdminPass struct {
	Login    string `json:"phone_number" db:"phone_number"`
	Password string `json:"password" db:"password"`
}

type LogIn struct {
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Password    string `json:"password" db:"password"`
}

type Token struct {
	AccessToken  string `json:"access_token" db:"access_token"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
	ExpireAt     int    `json:"expire_at" db:"expire_at"`
}

type LogInReq struct {
	Id          string `json:"id" db:"user_id"`
	FirstName   string `json:"first_name" db:"first_name"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Role        string `json:"role" db:"role"`
}

type PhoneNumber struct {
	PhoneNumber string `json:"phone_number" db:"phone_number"`
}

type Error struct {
	Error error
}
