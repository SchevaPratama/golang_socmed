package service

import (
	"context"

	"golang_socmed/internal/entity"
	helpers "golang_socmed/internal/helper"
	"golang_socmed/internal/model"
	"golang_socmed/internal/model/converter"
	"golang_socmed/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type BankAccountService struct {
	Repository *repository.BankAccountRepository
	Validate   *validator.Validate
	Log        *logrus.Logger
}

func NewBankAccountService(r *repository.BankAccountRepository, validate *validator.Validate, log *logrus.Logger) *BankAccountService {
	return &BankAccountService{Repository: r, Validate: validate, Log: log}
}

func (s *BankAccountService) List(ctx context.Context, userId string) ([]model.BankAccountResponse, error) {

	bankAccounts, err := s.Repository.List(userId)
	if err != nil {
		s.Log.Error("failed get product lists")
		return nil, err
	}

	newBankAccounts := make([]model.BankAccountResponse, len(bankAccounts))
	for i, bankAccount := range bankAccounts {
		newBankAccounts[i] = *converter.BankAccountConverter(&bankAccount)
	}

	return newBankAccounts, nil
}

func (s *BankAccountService) Get(ctx context.Context, id string, userId string) (*model.BankAccountResponse, error) {
	bankAccount := new(entity.BankAccount)

	err := s.Repository.Get(id, bankAccount, userId)
	if err != nil {
		s.Log.Error("failed get bank account detail")
		return nil, &fiber.Error{
			Code:    404,
			Message: "Data NotFound",
		}
	}

	return converter.BankAccountConverter(bankAccount), nil
}

func (s *BankAccountService) Create(ctx context.Context, request *model.BankAccountRequest, userId string) error {
	// handle request
	err := helpers.ValidationError(s.Validate, request)
	if err != nil {
		s.Log.Error("validation error: ", err.Error())
		return &fiber.Error{
			Code:    400,
			Message: err.Error(),
		}
	}

	payload := &entity.BankAccount{
		ID:       uuid.New().String(),
		Name:     request.BankAccountName,
		Number:   request.BankAccountNumber,
		BankName: request.BankName,
		UserId:   userId,
	}

	err = s.Repository.Create(payload)
	if err != nil {
		s.Log.Error("failed to insert new bank account")
		return &fiber.Error{
			Code:    500,
			Message: err.Error(),
		}
	}

	return nil
}

func (s *BankAccountService) Update(ctx context.Context, id string, request *model.BankAccountRequest, userId string) error {
	// handle request
	err := helpers.ValidationError(s.Validate, request)
	if err != nil {
		s.Log.Error("validation error: ", err.Error())
		return &fiber.Error{
			Code:    400,
			Message: err.Error(),
		}
	}

	bankAccount := new(entity.BankAccount)
	err = s.Repository.Get(id, bankAccount, userId)
	if err != nil {
		s.Log.WithError(err).Error("failed get bankAccount detail: ", err.Error())
		if err.Error() == "sql: no rows in result set" {
			return &fiber.Error{
				Code:    404,
				Message: "Data not found",
			}
		}
		return err
	}

	bankAccount.BankName = request.BankName
	bankAccount.Name = request.BankAccountName
	bankAccount.Number = request.BankAccountNumber

	err = s.Repository.Update(id, bankAccount)
	if err != nil {
		s.Log.Error("failed to update bank account")
		return &fiber.Error{
			Code:    500,
			Message: err.Error(),
		}
	}

	return nil
}

func (s *BankAccountService) Delete(ctx context.Context, id string, userId string) error {
	bankAccount := new(entity.BankAccount)
	err := s.Repository.Get(id, bankAccount, userId)
	if err != nil {
		s.Log.WithError(err).Error("failed get bankAccount detail: ", err.Error())
		if err.Error() == "sql: no rows in result set" {
			return &fiber.Error{
				Code:    404,
				Message: "Data not found",
			}
		}
		return err
	}

	err = s.Repository.Delete(id)
	if err != nil {
		s.Log.Error("failed to delete bank account")
		return &fiber.Error{
			Code:    500,
			Message: err.Error(),
		}
	}

	return nil
}
