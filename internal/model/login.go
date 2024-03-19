package model

type LoginRequest struct {
	CredentialType  string `json:"credentialType" validate:"required,oneof=phone email"`
	CredentialValue string `json:"credentialValue" validate:"required,customCredential"`
	Password        string `json:"password" validate:"required,min=5,max=15"`
}

type LoginRegisterResponse struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}
