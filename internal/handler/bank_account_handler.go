package handler

import (
	"golang_socmed/internal/model"
	"golang_socmed/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type BankAccountHandler struct {
	Service *service.BankAccountService
	Log     *logrus.Logger
}

func NewBankAccountHandler(s *service.BankAccountService, log *logrus.Logger) *BankAccountHandler {
	return &BankAccountHandler{
		Service: s,
		Log:     log,
	}
}

func (b *BankAccountHandler) List(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	products, err := b.Service.List(c.UserContext(), userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    products,
	})
}

func (b *BankAccountHandler) Get(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return err
	}

	bankAccount, err := b.Service.Get(c.UserContext(), id.String(), userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    bankAccount,
	})
}

func (b *BankAccountHandler) Create(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	request := new(model.BankAccountRequest)

	if err := c.BodyParser(request); err != nil {
		b.Log.WithError(err).Error("failed to process request")
		return fiber.ErrBadRequest
	}

	err := b.Service.Create(c.UserContext(), request, userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "account added successfully",
	})
}

func (b *BankAccountHandler) Update(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	request := new(model.BankAccountRequest)
	if err := c.BodyParser(request); err != nil {
		b.Log.WithError(err).Error("failed to process request")
		return fiber.ErrBadRequest
	}

	err = b.Service.Update(c.UserContext(), id.String(), request, userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "account updated successfully",
	})
}

func (b *BankAccountHandler) Delete(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return &fiber.Error{
			Code:    400,
			Message: err.Error(),
		}
	}

	err = b.Service.Delete(c.UserContext(), id.String(), userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "account deleted successfully",
	})
}
