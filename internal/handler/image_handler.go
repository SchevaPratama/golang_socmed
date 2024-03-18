package handler

import (
	"golang_socmed/internal/service"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ImageHandler struct {
	Service *service.ImageService
	Log     *logrus.Logger
}

func NewImageHandler(s *service.ImageService, log *logrus.Logger) *ImageHandler {
	return &ImageHandler{
		Service: s,
		Log:     log,
	}
}

func (h *ImageHandler) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "File is required",
		})
	}

	// Validate file size
	if file.Size < 10*1024 || file.Size > 2*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "File size must be between 10KB and 2MB",
		})
	}

	// Validate file extension (.jpg or .jpeg)
	ext := path.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "File must be a .jpg or .jpeg",
		})
	}

	filename := uuid.NewString() + ext

	buffer, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	defer buffer.Close()

	url, err := h.Service.Upload(c.UserContext(), buffer, filename)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"imageUrl": url,
	})
}
