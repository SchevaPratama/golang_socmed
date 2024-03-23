package router

import (
	"golang_socmed/internal/handler"
	"golang_socmed/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                *fiber.App
	ProductHandler     *handler.ProductHandler
	UserHandler        *handler.UserHandler
	ImageHandler       *handler.ImageHandler
	BankAccountHandler *handler.BankAccountHandler
	PostHandler        *handler.PostHandler
}

func (c *RouteConfig) Setup() {

	authMiddleware := middleware.NewAuthMiddleware()

	c.App.Post("/v1/user/register", c.UserHandler.Register)
	c.App.Post("/v1/user/login", c.UserHandler.Login)

	//c.App.Patch("/v1/user", c.UserHandler.UpdateUser, authMiddleware)

	image := c.App.Group("/v1/image", authMiddleware)
	image.Post("/", c.ImageHandler.Upload)

	// product := c.App.Group("/v1/product", authMiddleware)
	// product.Get("", c.ProductHandler.List)
	// product.Post("", c.ProductHandler.Create)
	// product.Get("/:id", c.ProductHandler.Get)
	// product.Delete("/:id", c.ProductHandler.Delete)
	// product.Put("/:id", c.ProductHandler.Update)
	// product.Post("/:id/stock", c.ProductHandler.UpdateStock)
	// product.Post("/:id/buy", c.ProductHandler.Buy)

	friend := c.App.Group("/v1/friend", authMiddleware)
	friend.Get("", c.UserHandler.GetFriends)
	friend.Post("", c.UserHandler.AddFriend)
	friend.Delete("", c.UserHandler.DeleteFriend)

	user := c.App.Group("/v1/user", authMiddleware)
	user.Post("/link/email", c.UserHandler.LinkPhoneEmail)
	user.Post("/link/phone", c.UserHandler.LinkPhoneEmail)
	user.Patch("", c.UserHandler.UpdateUser)

	// c.App.Patch("/v1/bank/account", authMiddleware, func(c *fiber.Ctx) error {
	// 	return c.SendStatus(http.StatusNotFound)
	// })
	// bankAccount := c.App.Group("/v1/bank/account", authMiddleware)
	// bankAccount.Get("/", c.BankAccountHandler.List)
	// bankAccount.Get("/:id", c.BankAccountHandler.Get)
	// bankAccount.Patch("/:id", c.BankAccountHandler.Update)
	// bankAccount.Delete("/:id", c.BankAccountHandler.Delete)
	// bankAccount.Post("/", c.BankAccountHandler.Create)

	post := c.App.Group("/v1/post", authMiddleware)
	post.Get("/", c.PostHandler.List)
	post.Post("/", c.PostHandler.Create)

}
