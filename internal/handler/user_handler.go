package handler

import (
	"golang_socmed/internal/model"
	"golang_socmed/internal/service"

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

func (h *UserHandler) Login(c *fiber.Ctx) error {
	request := new(model.LoginRequest)

	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
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

	resp, err := h.Service.Register(c.UserContext(), request)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data":    resp,
	})
}
