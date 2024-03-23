package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      "SegoKuning Social App",
		ErrorHandler: NewErrorHandler(),
		Prefork:      true,
	})

	app.Use(logger.New(logger.Config{
		TimeFormat: "2 Jan 2006 15:04:05",
		TimeZone:   "Asia/Jakarta",
		Format:     "[${time}] ${status} - ${method} ${path} body: ${body} queryParams: ${queryParams}\n",
	}))

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return ctx.Status(code).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}
}
