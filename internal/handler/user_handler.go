package handler

import (
	"golang_socmed/internal/model"
	"golang_socmed/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	Service *service.UserService
	Log     *logrus.Logger
}

func NewUserHandler(s *service.UserService, log *logrus.Logger) *UserHandler {
	return &UserHandler{
		Service: s,
		Log:     log,
	}
}

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("customCredential", model.ValidateCredential)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	request := new(model.LoginRequest)

	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validate the request
	if err := validate.Struct(request); err != nil {
		// Handle validation errors
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	resp, err := h.Service.Login(c.UserContext(), request)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User logged successfully",
		"data":    resp,
	})
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	request := new(model.RegisterRequest)

	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validate the request
	if err := validate.Struct(request); err != nil {
		// Handle validation errors
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	resp, err := h.Service.Register(c.UserContext(), request)
	if err != nil {
		return err
	}

	if request.CredentialType == "email" {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User registered successfully",
			"data": fiber.Map{
				"email":       resp.Email,
				"name":        resp.Name,
				"accessToken": resp.AccessToken,
			},
		})
	} else {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User registered successfully",
			"data": fiber.Map{
				"phone":       resp.Phone,
				"name":        resp.Name,
				"accessToken": resp.AccessToken,
			},
		})
	}

}
