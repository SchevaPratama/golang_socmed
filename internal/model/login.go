package model

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=5,max=15"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

type LoginRegisterResponse struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}
