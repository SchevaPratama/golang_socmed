package entity

import "database/sql"

type User struct {
	ID        string
	Email     sql.NullString
	Phone     sql.NullString
	Name      string
	Password  string
	Friends   []string
	CreatedAt string
}

func (prod *User) TableName() string {
	return "users"
}
