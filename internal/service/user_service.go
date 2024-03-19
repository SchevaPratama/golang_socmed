package service

import (
	"context"
	"golang_socmed/internal/entity"

	// helpers "golang_socmed/internal/helper"
	"golang_socmed/internal/model"
	"golang_socmed/internal/repository"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repository *repository.UserRepository
	Validate   *validator.Validate
	Log        *logrus.Logger
}

func NewUserService(r *repository.UserRepository, validate *validator.Validate, log *logrus.Logger) *UserService {
	return &UserService{Repository: r, Validate: validate, Log: log}
}

func (s *UserService) Register(ctx context.Context, request *model.RegisterRequest) (*model.RegisterResponse, error) {
	// handle request
	// err := helpers.ValidationError(s.Validate, request)
	// if err != nil {
	// 	return nil, &fiber.Error{
	// 		Code:    fiber.StatusBadRequest,
	// 		Message: err.Error(),
	// 	}
	// }
	// log.Println(err)

	userData, _ := s.getEmailOrPhone(request.CredentialType, request.CredentialValue)
	if userData != nil && request.CredentialType == "email" {
		return nil, &fiber.Error{
			Code:    fiber.StatusConflict,
			Message: "Email already exists",
		}
	}

	if userData != nil && request.CredentialType == "phone" {
		return nil, &fiber.Error{
			Code:    fiber.StatusConflict,
			Message: "Phone already exists",
		}
	}

	bcryptSalt, err := strconv.Atoi(os.Getenv("BCRYPT_SALT"))
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcryptSalt)
	if err != nil {
		log.Println("Error hashedPassword")
	}

	var user *entity.User

	if request.CredentialType == "email" {
		user = &entity.User{
			Email:    request.CredentialValue,
			Name:     request.Name,
			Password: string(hashedPassword),
		}
	} else {
		user = &entity.User{
			Phone:    request.CredentialValue,
			Name:     request.Name,
			Password: string(hashedPassword),
		}
	}

	err = s.Repository.Create(request.CredentialType, request.CredentialValue, user)
	if err != nil {
		log.Println("Gagal menyimpan user", err)
	}

	day := time.Hour * 24

	claims := jtoken.MapClaims{
		"ID":   user.ID,
		"name": user.Name,
		"exp":  time.Now().Add(day * 1).Unix(),
	}

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

	//TODO :: config secret jwt
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}
	var resp *model.RegisterResponse

	if request.CredentialType == "email" {
		resp = &model.RegisterResponse{
			Email:       user.Email,
			Name:        user.Name,
			AccessToken: t,
		}
	} else {
		resp = &model.RegisterResponse{
			Phone:       user.Phone,
			Name:        user.Name,
			AccessToken: t,
		}
	}

	return resp, nil
}

func (s *UserService) Login(ctx context.Context, request *model.LoginRequest) (*model.LoginRegisterResponse, error) {

	// handle request
	// err := s.Validate.Struct(request)
	// if err != nil {
	// 	return nil, &fiber.Error{
	// 		Code:    400,
	// 		Message: err.Error(),
	// 	}
	// }

	user, err := s.getEmailOrPhone(request.CredentialType, request.CredentialValue)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		s.Log.WithError(err).Error("Error Password is Wrong", err.Error())
		return nil, &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Password is wrong",
		}
	}

	day := time.Hour * 24

	claims := jtoken.MapClaims{
		"ID":   user.ID,
		"Name": user.Name,
		"exp":  time.Now().Add(day * 1).Unix(),
	}

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	resp := &model.LoginRegisterResponse{
		Email:       user.Email,
		Phone:       user.Phone,
		Name:        user.Name,
		AccessToken: t,
	}

	return resp, nil
}

func (s *UserService) getEmailOrPhone(credentialType string, credentialValue string) (*entity.User, error) {
	var user entity.User
	if credentialType == "email" {
		user = entity.User{Email: credentialValue}
	}
	if credentialType == "phone" {
		user = entity.User{Phone: credentialValue}
	}
	err := s.Repository.GetByEmailOrPhone(credentialType, credentialValue, &user)
	if err != nil {
		s.Log.WithError(err).Error("Error Get User by Username", err.Error())
		return nil, &fiber.Error{
			Code:    fiber.StatusNotFound,
			Message: "User NotFound",
		}
	}
	return &user, nil
}
