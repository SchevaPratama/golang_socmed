package model

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
	CredentialType  string `json:"credentialType" validate:"required,oneof=phone email"`
	CredentialValue string `json:"credentialValue" validate:"required,customCredential"`
	Name            string `json:"name" validate:"required,min=5,max=50"`
	Password        string `json:"password" validate:"required,min=5,max=15"`
}

type RegisterResponse struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

// Custom validation function for CredentialValue
func ValidateCredential(fl validator.FieldLevel) bool {
	credentialType := fl.Parent().FieldByName("CredentialType").String()
	credentialValue := fl.Field().String()

	switch credentialType {
	case "email":
		return isEmailValid(credentialValue)
	case "phone":
		return isPhoneValid(credentialValue)
	default:
		return false
	}
}

// isEmailValid checks if the email provided is valid by regex.
func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

// isPhoneValid checks if the phone number is valid.
func isPhoneValid(p string) bool {
	// Implement your phone number validation logic here
	// For example, check if it starts with "+" and has a length between 7 and 13
	return strings.HasPrefix(p, "+") && (len(p) >= 7 && len(p) <= 13)
}
