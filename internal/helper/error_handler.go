package helpers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func NewErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var e *fiber.Error
		if errors.As(err, &e) {
			return c.Status(e.Code).JSON(fiber.Map{
				"code":    e.Code,
				"message": e.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}
}
