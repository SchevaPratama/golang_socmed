package entity

type User struct {
	ID       string
	Email    string
	Phone    string
	Name     string
	Password string
}

func (prod *User) TableName() string {
	return "users"
}
