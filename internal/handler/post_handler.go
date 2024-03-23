package handler

import (
	"golang_socmed/internal/model"
	"golang_socmed/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type PostHandler struct {
	Service *service.PostService
	Log     *logrus.Logger
}

func NewPostHandler(s *service.PostService, log *logrus.Logger) *PostHandler {
	return &PostHandler{
		Service: s,
		Log:     log,
	}
}

func (b *PostHandler) List(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	// searchTag := c.Query("searchTag", "hehe", "[]string")

	// search := c.Query("search")

	// limit := c.QueryInt("limit", 5)
	// offset := c.QueryInt("offset", 0)

	// filter := &model.PostFilter{
	// 	Search: &search,
	// 	Limit:  &limit,
	// 	Offset: &offset,
	// }

	filter := new(model.PostFilter)
	filter.Limit = c.QueryInt("limit", 5)
	filter.Offset = c.QueryInt("offset", 0)

	if err := c.QueryParser(filter); err != nil {
		b.Log.WithError(err).Error("failed to process request")
		return fiber.ErrBadRequest
	}

	products, err := b.Service.List(c.UserContext(), filter, userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"data":    products,
		"meta": fiber.Map{
			"limit":  1,
			"offset": 1,
			"total":  len(products),
		},
	})
}

func (b *PostHandler) Create(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	request := new(model.PostRequest)

	if err := c.BodyParser(request); err != nil {
		b.Log.WithError(err).Error("failed to process request")
		return fiber.ErrBadRequest
		// return &fiber.Error{Message: "Opppss", Code: 400}
	}

	err := b.Service.Create(c.UserContext(), request, userId)
	if err != nil {
		// return fiber.ErrBadRequest
		return &fiber.Error{Message: err.Error(), Code: 400}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    1,
		"message": "success insert new posts",
		"data":    request,
	})
}
