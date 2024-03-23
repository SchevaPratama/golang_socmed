package config

import (
	"golang_socmed/internal/handler"
	"golang_socmed/internal/repository"
	"golang_socmed/internal/router"
	"golang_socmed/internal/service"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	DB       *sqlx.DB
	App      *fiber.App
	Config   *viper.Viper
	Validate *validator.Validate
	Log      *logrus.Logger
	Aws      *s3.Client
}

func Bootstrap(config *BootstrapConfig) {
	// aws
	aws := NewAws(config.Config)

	// setup repositories
	userRepository := repository.NewUserRepository(config.DB)
	productRepository := repository.NewProductRepository(config.DB)
	bankAccountRepository := repository.NewBankAccountRepository(config.DB)
	postRepository := repository.NewPostRepository(config.DB)
	commentRepository := repository.NewCommentRepository(config.DB)

	// setup services
	imageService := service.NewImageService(aws, config.Validate, config.Log)
	productService := service.NewProductService(productRepository, imageService, config.Validate, config.Log)
	userService := service.NewUserService(userRepository, config.Validate, config.Log)
	bankAccountService := service.NewBankAccountService(bankAccountRepository, config.Validate, config.Log)
	postService := service.NewPostService(postRepository, config.Validate, config.Log, commentRepository)

	// setup handler
	productHandler := handler.NewProductHandler(productService, config.Log)
	userHandler := handler.NewUserHandler(userService, config.Log)
	bankAccountHandler := handler.NewBankAccountHandler(bankAccountService, config.Log)
	ImageHandler := handler.NewImageHandler(imageService, config.Log)
	postHandler := handler.NewPostHandler(postService, config.Log)

	// recover from panic
	config.App.Use(func(c *fiber.Ctx) error {

		defer func() {
			if r := recover(); r != nil {
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Internal Server Error",
				})
			}
		}()
		return c.Next()
	})

	// setup route
	routeConfig := router.RouteConfig{
		App:                config.App,
		ProductHandler:     productHandler,
		UserHandler:        userHandler,
		ImageHandler:       ImageHandler,
		BankAccountHandler: bankAccountHandler,
		PostHandler:        postHandler,
	}

	routeConfig.Setup()
}
