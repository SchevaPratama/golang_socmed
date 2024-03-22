package entity

import (
	"database/sql"
)

type User struct {
	ID        string
	Email     sql.NullString
	Phone     sql.NullString
	Name      string
	Password  string
	ImageUrl  sql.NullString `db:"image_url"`
	Friends   []string
	CreatedAt string
	UpdatedAt string
}

func (prod *User) TableName() string {
	return "users"
}
