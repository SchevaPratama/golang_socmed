package handler

import (
	"golang_socmed/internal/model"
	"golang_socmed/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	Service *service.ProductService
	Log     *logrus.Logger
}

func NewProductHandler(s *service.ProductService, log *logrus.Logger) *ProductHandler {
	return &ProductHandler{
		Service: s,
		Log:     log,
	}
}

func (b *ProductHandler) List(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}
	keyword := c.Query("search")
	condition := c.Query("condition")
	sortBy := c.Query("sortBy")
	orderBy := c.Query("orderBy")
	maxPrice, _ := strconv.Atoi(c.Query("maxPrice"))
	minPrice, _ := strconv.Atoi(c.Query("minPrice"))
	userOnly, _ := strconv.ParseBool(c.Query("userOnly"))
	showEmptyStock, _ := strconv.ParseBool(c.Query("showEmptyStock"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	filter := &model.ProductFilter{
		Condition:      &condition,
		Keyword:        &keyword,
		SortBy:         &sortBy,
		OrderBy:        &orderBy,
		MaxPrice:       &maxPrice,
		MinPrice:       &minPrice,
		UserOnly:       &userOnly,
		ShowEmptyStock: &showEmptyStock,
		Limit:          &limit,
		Offset:         &offset}

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
			"limit":  limit,
			"offset": offset,
			"total":  len(products),
		},
	})
}

func (b *ProductHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	product, err := b.Service.Get(c.UserContext(), id.String())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    1,
		"message": "detail of product",
		"data":    product,
	})
}

func (b *ProductHandler) Create(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	request := new(model.ProductRequest)

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
		"message": "success insert new product",
		"data":    request,
	})
}

func (b *ProductHandler) Delete(c *fiber.Ctx) error {
	id, errUuid := uuid.Parse(c.Params("id"))
	if errUuid != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errUuid.Error(),
		})
	}

	_, err := b.Service.Get(c.UserContext(), id.String())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := b.Service.Delete(c.UserContext(), id.String()); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    1,
		"message": "success delete a product",
	})
}

func (b *ProductHandler) Update(c *fiber.Ctx) error {
	id, errUUID := uuid.Parse(c.Params("id"))
	if errUUID != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errUUID.Error(),
		})
	}

	_, err := b.Service.Get(c.UserContext(), id.String())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	request := new(model.ProductRequest)
	if err := c.BodyParser(request); err != nil {
		b.Log.WithError(err).Error("failed to process request")
		return fiber.ErrBadRequest
	}

	if err := b.Service.Update(c.UserContext(), id.String(), request); err != nil {
		return &fiber.Error{Message: err.Error(), Code: 400}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    1,
		"message": "success Update a product",
		"data":    request,
	})
}

func (b *ProductHandler) UpdateStock(c *fiber.Ctx) error {
	id, errUUID := uuid.Parse(c.Params("id"))
	if errUUID != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errUUID.Error(),
		})
	}

	product, err := b.Service.Get(c.UserContext(), id.String())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	request := new(model.StockRequest)
	if err := c.BodyParser(request); err != nil {
		b.Log.WithError(err).Error("failed to process request")
		return fiber.ErrBadRequest
	}
	product.Stock = request.Stock

	if err := b.Service.UpdateStock(c.UserContext(), id.String(), request); err != nil {
		return &fiber.Error{Message: err.Error(), Code: 400}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    1,
		"message": "success Update a product",
		"data":    product,
	})
}

func (b *ProductHandler) Buy(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	id, errUUID := uuid.Parse(c.Params("id"))
	if errUUID != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errUUID.Error(),
		})
	}

	request := new(model.BuyRequest)

	if err := c.BodyParser(request); err != nil {
		b.Log.WithError(err).Error("failed to process request")
		return fiber.ErrBadRequest
	}

	request.ProductId = id.String()

	err := b.Service.Buy(c.UserContext(), request, userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "payment processed successfully",
		"data":    request,
	})
}
