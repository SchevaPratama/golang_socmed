package handler

import (
	helpers "golang_socmed/internal/helper"
	"golang_socmed/internal/model"
	"golang_socmed/internal/service"
	"log"
	"net/url"
	"path/filepath"
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

	filter := model.FriendFilter{
		Search:     "",
		SortBy:     "createdAt",
		OrderBy:    "desc",
		OnlyFriend: false,
		Limit:      5,
		Offset:     0,
	}

	maps := c.Queries()

	filter.Search = c.Query("search")
	if val, ok := maps["sortBy"]; ok {
		if val != "" {
			if val == "createdAt" || val == "friendCount" {
				filter.SortBy = val
			} else {
				return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "sortBy value not defined"}
			}
		} else {
			return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "sortBy can't be empty"}
		}
	} else {
		filter.SortBy = "createdAt"
	}

	if val, ok := maps["orderBy"]; ok {
		if val != "" {
			if val == "asc" || val == "desc" {
				filter.OrderBy = val
			} else {
				return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "orderBy value not defined"}
			}
		} else {
			return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "orderBy can't be empty"}
		}
	} else {
		filter.OrderBy = "desc"
	}

	if val, ok := maps["onlyFriend"]; ok {
		if val != "" {
			onlyFriendParsed, _ := strconv.ParseBool(c.Query("onlyFriend"))
			if onlyFriendParsed == true || onlyFriendParsed == false {
				filter.OnlyFriend = onlyFriendParsed
			} else {
				return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "onlyFriend Param is not boolean"}
			}
		} else {
			return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "Limit can't be empty"}
		}
	} else {
		filter.OnlyFriend = false
	}

	if c.Query("onlyFriend") != "" {
		onlyFriend, err := strconv.ParseBool(c.Query("onlyFriend"))
		if err != nil {
			return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "onlyFriend Param is not boolean"}
		} else {
			filter.OnlyFriend = onlyFriend
		}
	} else {
		filter.OnlyFriend = false
	}

	if val, ok := maps["limit"]; ok {
		if val != "" {
			limitParsed, _ := strconv.Atoi(val)
			if limitParsed < 0 {
				return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "Limit minimal 0"}
			}
			filter.Limit = limitParsed
		} else {
			return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "Limit can't be empty"}
		}
	} else {
		filter.Limit = 5
	}

	// if c.Query("limit") == "" {
	// 	return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "Limit can't be empty"}
	// } else {
	// 	filter.Limit = 5
	// }

	// limitParsed, _ := strconv.Atoi(c.Query("limit"))
	// if limitParsed < 0 {
	// 	return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "Limit minimal 0"}
	// }

	// filter.Limit = limitParsed

	if val, ok := maps["offset"]; ok {
		if val != "" {
			offsetParsed, _ := strconv.Atoi(val)
			if offsetParsed < 0 {
				return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "Offset minimal 0"}
			}
			filter.Offset = offsetParsed
		} else {
			return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "Offset can't be empty"}
		}
	} else {
		filter.Offset = 0
	}

	log.Println(c.Queries())

	// if c.Query("offset") == "" {
	// 	return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "Offset can't be empty"}
	// } else {
	// 	filter.Offset = 0
	// }

	// offsetParsed, _ := strconv.Atoi(c.Query("offset"))
	// if offsetParsed < 0 {
	// 	return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: "Offet minimal 0"}
	// }

	// filter.Offset = offsetParsed

	errParse := c.QueryParser(&filter)
	if errParse != nil {
		h.Log.WithError(errParse).Error("failed to process request")
		return fiber.ErrBadRequest
	}

	log.Println(filter)

	users, err := h.Service.GetFriends(c.UserContext(), filter, userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"data":    users,
		"meta": fiber.Map{
			"limit":  filter.Limit,
			"offset": filter.Offset,
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
	var types string
	// types := strings.ReplaceAll(c.OriginalURL(), "/v1/user/link/", "")
	if strings.ReplaceAll(c.OriginalURL(), "/v1/user/link", "") == "" {
		types = "email"
	}
	if strings.ReplaceAll(c.OriginalURL(), "/v1/user/link", "") == "/phone" {
		types = "phone"
	}
	log.Println(types)
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

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	var request model.UpdateProfileRequest

	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	err := validate.Struct(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	isImage, err := isImageUrl(request.ImageUrl)
	if err != nil || isImage == false {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ImageUrl must be a valid image url",
		})
	}

	err = h.Service.UpdateUser(c.UserContext(), userId, request)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated",
	})
}

func isImageUrl(urlStr string) (bool, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false, err
	}

	ext := filepath.Ext(u.Path)
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp":
		return true, nil
	default:
		return false, nil
	}
}
