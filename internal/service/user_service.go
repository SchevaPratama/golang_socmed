package service

import (
	"context"
	"database/sql"
	"golang_socmed/internal/entity"
	helpers "golang_socmed/internal/helper"
	"golang_socmed/internal/model/converter"

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
	"github.com/google/uuid"
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
			ID:       uuid.New().String(),
			Email:    sql.NullString{String: request.CredentialValue, Valid: true},
			Name:     request.Name,
			Password: string(hashedPassword),
		}
	} else {
		user = &entity.User{
			ID:       uuid.New().String(),
			Phone:    sql.NullString{String: request.CredentialValue, Valid: true},
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
			Email:       user.Email.String,
			Name:        user.Name,
			AccessToken: t,
		}
	} else {
		resp = &model.RegisterResponse{
			Phone:       user.Phone.String,
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

	log.Println(user)

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
		"name": user.Name,
		"exp":  time.Now().Add(day * 1).Unix(),
	}

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	resp := &model.LoginRegisterResponse{
		Email:       user.Email.String,
		Phone:       user.Phone.String,
		Name:        user.Name,
		AccessToken: t,
	}

	return resp, nil
}

func (s *UserService) GetFriends(ctx context.Context, filter model.FriendFilter, userId string) ([]model.FriendResponse, error) {
	// if err := helpers.ValidationError(s.Validate, filter); err != nil {
	// 	s.Log.WithError(err).Error("failed to validate request query params")
	// 	return nil, err
	// }

	users, err := s.Repository.GetUsers(filter, userId)
	if err != nil {
		s.Log.WithError(err).Error("failed get product lists")
		return nil, err
	}

	newusers := make([]model.FriendResponse, len(users))
	for i, user := range users {
		newusers[i] = *converter.FriendConverter(&user)
	}

	return newusers, nil
}

func (s *UserService) AddFriend(ctx context.Context, userId string, request *model.FriendRequest) error {
	if err := helpers.ValidationError(s.Validate, request); err != nil {
		s.Log.WithError(err).Error("failed to validate request body")
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		}
	}

	_, errUuid := uuid.Parse(request.UserId)
	if errUuid != nil {
		return &fiber.Error{
			Code:    fiber.StatusNotFound,
			Message: errUuid.Error(),
		}
	}

	_, errs := s.Repository.GetById(request)
	if errs != nil {
		s.Log.WithError(errs).Error("User Not Found", errs.Error())
		return errs
	}

	err := s.Repository.AddFriend(request.UserId, userId)
	if err != nil {
		s.Log.WithError(err).Error("failed to update data")
		return err
	}

	return nil
}

func (s *UserService) DeleteFriend(ctx context.Context, userId string, request *model.FriendRequest) error {
	// if err := s.Validate.Struct(request); err != nil {
	if err := helpers.ValidationError(s.Validate, request); err != nil {
		s.Log.WithError(err).Error("failed to validate request body")
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		}
	}

	_, errs := s.Repository.GetById(request)
	if errs != nil {
		s.Log.WithError(errs).Error("User Not Found", errs.Error())
		return &fiber.Error{
			Code:    fiber.StatusNotFound,
			Message: "User Not Found",
		}
	}

	err := s.Repository.DeleteFriend(request.UserId, userId)
	if err != nil {
		s.Log.WithError(err).Error("failed to update data")
		return err
	}

	return nil
}

func (s *UserService) LinkPhoneEmail(ctx context.Context, userId string, types string, value string) error {
	log.Println(value)
	user, errs := s.Repository.GetById(&model.FriendRequest{UserId: userId})
	if errs != nil {
		s.Log.WithError(errs).Error("failed get user detail")
		return errs
	}

	if types == "email" && user.Email.String != "" {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Can't Change Email Of This Account",
		}
	}

	if types == "phone" && user.Phone.String != "" {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Can't Change Phone Of This Account",
		}
	}

	userData, _ := s.getEmailOrPhone(types, value)
	if userData != nil && types == "email" {
		return &fiber.Error{
			Code:    fiber.StatusConflict,
			Message: "Email already exists",
		}
	}

	if userData != nil && types == "phone" {
		return &fiber.Error{
			Code:    fiber.StatusConflict,
			Message: "Phone already exists",
		}
	}

	err := s.Repository.LinkPhoneEmail(types, value, userId)
	if err != nil {
		s.Log.WithError(err).Error("failed to update data")
		return err
	}

	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, userId string, request model.UpdateProfileRequest) error {

	user, errs := s.Repository.GetById(&model.FriendRequest{UserId: userId})
	if errs != nil {
		s.Log.WithError(errs).Error("failed get user detail")
		return errs
	}

	user.Name = request.Name
	user.ImageUrl = sql.NullString{String: request.ImageUrl, Valid: true}

	err := s.Repository.UpdateUser(user)
	if err != nil {
		s.Log.WithError(err).Error("failed to update data")
		return err
	}

	return nil
}

func (s *UserService) getEmailOrPhone(credentialType string, credentialValue string) (*entity.User, error) {
	var user entity.User
	if credentialType == "email" {
		user = entity.User{Email: sql.NullString{String: credentialValue, Valid: true}}
	}
	if credentialType == "phone" {
		user = entity.User{Phone: sql.NullString{String: credentialValue, Valid: true}}
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
