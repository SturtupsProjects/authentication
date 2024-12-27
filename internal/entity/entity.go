package entity

type DBClient struct {
	Id        string `db:"id"`
	FullName  string `db:"full_name"`
	Address   string `db:"address"`
	Phone     string `db:"phone"`
	Type      string `db:"type"`
	CompanyId string `db:"company_id"`
}

type LogInToken struct {
	UserId      string `db:"user_id"`
	Role        string `db:"role"`
	FirstName   string `db:"first_name"`
	PhoneNumber string `db:"phone_number"`
	CompanyId   string `db:"company_id"`
}
