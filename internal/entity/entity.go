package entity

type DBClient struct {
	Id       string `db:"id"`
	FullName string `db:"full_name"`
	Address  string `db:"address"`
	Phone    string `db:"phone"`
}
