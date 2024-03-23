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

type FriendRequest struct {
	UserId string `json:"userId" validate:"required,min=5,max=50"`
}

type FriendResponse struct {
	UserId      string `json:"userId"`
	Name        string `json:"name"`
	ImageUrl    string `json:"imageUrl"`
	FriendCount int    `json:"friendCount"`
	CreatedAt   string `json:"createdAt"`
}

type FriendFilter struct {
	Limit      int    `json:"limit" validate:"min=0"`
	Offset     int    `json:"offset" validate:"min=0"`
	SortBy     string `json:"sortBy" validate:"oneof=createdAt friendCount"`
	OrderBy    string `json:"orderBy" validate:"oneof=asc desc"`
	OnlyFriend bool   `json:"onlyFriend" validate:"oneof=true false"`
	Search     string `json:"search"`
}

type EmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PhoneRequest struct {
	Phone string `json:"phone" validate:"required,phone"`
}

type UpdateProfileRequest struct {
	Name     string `json:"name" validate:"required,min=5,max=50"`
	ImageUrl string `json:"imageUrl" validate:"required,url,min=1"`
}

// ValidateCredential Custom validation function for CredentialValue
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

func PhoneValidation(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// Check if the phone number starts with "+" and has a length between 7 and 13
	return len(phone) >= 7 && len(phone) <= 13 && phone[0] == '+'
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
