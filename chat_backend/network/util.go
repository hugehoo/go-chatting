package network

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/fiber/v2/middleware/basicauth"
)

func response(c *fiber.Ctx, status int, res interface{}, data ...string) {
	c.Status(status).JSON(data)
}
