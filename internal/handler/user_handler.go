package handler

import (
	helpers "golang_socmed/internal/helper"
	"golang_socmed/internal/model"
	"golang_socmed/internal/service"
	"strconv"
	"strings"

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
	validate.RegisterValidation("phone", model.PhoneValidation)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	request := new(model.LoginRequest)

	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	err := helpers.LoginValidationError(validate, *request)
	if err != nil {
		h.Log.Error("validation error: ", err.Error())
		return &fiber.Error{
			Code:    400,
			Message: err.Error(),
		}
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

	err := helpers.RegisterValidationError(validate, *request)
	if err != nil {
		h.Log.Error("validation error: ", err.Error())
		return &fiber.Error{
			Code:    400,
			Message: err.Error(),
		}
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

func (h *UserHandler) GetFriends(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	keyword := c.Query("search")
	sortBy := c.Query("sortBy")
	orderBy := c.Query("orderBy")
	onlyFriend, _ := strconv.ParseBool(c.Query("onlyFriend"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	filter := &model.FriendFilter{
		Search:     &keyword,
		SortBy:     &sortBy,
		OrderBy:    &orderBy,
		OnlyFriend: &onlyFriend,
		Limit:      &limit,
		Offset:     &offset}

	if err := c.QueryParser(filter); err != nil {
		h.Log.WithError(err).Error("failed to process request")
		return fiber.ErrBadRequest
	}

	users, err := h.Service.GetFriends(c.UserContext(), filter, userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"data":    users,
		"meta": fiber.Map{
			"limit":  limit,
			"offset": offset,
			"total":  len(users),
		},
	})
}

func (h *UserHandler) AddFriend(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	request := new(model.FriendRequest)
	if err := c.BodyParser(request); err != nil {
		h.Log.WithError(err).Error("failed to process request")
		return fiber.ErrBadRequest
	}

	if userId == request.UserId {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "You Can't Add Friend Yourself"}
	}

	if err := h.Service.AddFriend(c.UserContext(), userId, request); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully add friend",
		"data":    request,
	})
}

func (h *UserHandler) DeleteFriend(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	request := new(model.FriendRequest)
	if err := c.BodyParser(request); err != nil {
		h.Log.WithError(err).Error("failed to process request")
		return fiber.ErrBadRequest
	}

	if userId == request.UserId {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "You Can't Delete Friend Yourself"}
	}

	if err := h.Service.DeleteFriend(c.UserContext(), userId, request); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully delete friend",
		"data":    request,
	})
}

func (h *UserHandler) LinkPhoneEmail(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	var value string
	types := strings.ReplaceAll(c.OriginalURL(), "/v1/user/link/", "")
	var message string

	if types == "email" {
		request := new(model.EmailRequest)
		if err := c.BodyParser(request); err != nil {
			h.Log.WithError(err).Error("failed to process request")
			return fiber.ErrBadRequest
		}
		err := helpers.ValidationError(validate, request)
		if err != nil {
			h.Log.Error("validation error: ", err.Error())
			return &fiber.Error{
				Code:    400,
				Message: err.Error(),
			}
		}
		value = request.Email
		message = "successfully link email to phone"
	}

	if types == "phone" {
		request := new(model.PhoneRequest)
		if err := c.BodyParser(request); err != nil {
			h.Log.WithError(err).Error("failed to process request")
			return fiber.ErrBadRequest
		}
		err := helpers.ValidationError(validate, request)
		if err != nil {
			h.Log.Error("validation error: ", err.Error())
			return &fiber.Error{
				Code:    400,
				Message: err.Error(),
			}
		}
		value = request.Phone
		message = "successfully link phone to email"
	}

	if err := h.Service.LinkPhoneEmail(c.UserContext(), userId, types, value); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}
