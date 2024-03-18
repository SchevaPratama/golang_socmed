package model

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Username string `json:"username" validate:"required,min=5,max=15"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}
