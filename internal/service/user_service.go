package service

import (
	"context"
	"golang_socmed/internal/entity"
	helpers "golang_socmed/internal/helper"
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

func (s *UserService) Register(ctx context.Context, request *model.RegisterRequest) (*model.LoginRegisterResponse, error) {

	// handle request
	err := helpers.ValidationError(s.Validate, request)
	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		}
	}

	user, err := s.getUsername(request.Username)
	if user != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusConflict,
			Message: "Username already exists",
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

	user = &entity.User{
		Name:     request.Name,
		Username: request.Username,
		Password: string(hashedPassword),
	}

	err = s.Repository.Create(user)
	if err != nil {
		log.Println("Gagal menyimpan user", err)
	}

	day := time.Hour * 24

	claims := jtoken.MapClaims{
		"ID":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(day * 1).Unix(),
	}

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

	//TODO :: config secret jwt
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	resp := &model.LoginRegisterResponse{
		Username:    user.Username,
		Name:        user.Name,
		AccessToken: t,
	}

	return resp, nil
}

func (s *UserService) Login(ctx context.Context, request *model.LoginRequest) (*model.LoginRegisterResponse, error) {

	// handle request
	err := s.Validate.Struct(request)
	if err != nil {
		return nil, &fiber.Error{
			Code:    400,
			Message: err.Error(),
		}
	}

	user, err := s.getUsername(request.Username)
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
		"ID":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(day * 1).Unix(),
	}

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	resp := &model.LoginRegisterResponse{
		Username:    user.Username,
		Name:        user.Name,
		AccessToken: t,
	}

	return resp, nil
}

func (s *UserService) getUsername(username string) (*entity.User, error) {
	user := entity.User{Username: username}
	err := s.Repository.GetByUsername(&user)
	if err != nil {
		s.Log.WithError(err).Error("Error Get User by Username", err.Error())
		return nil, &fiber.Error{
			Code:    fiber.StatusNotFound,
			Message: "User NotFound",
		}
	}
	return &user, nil
}
