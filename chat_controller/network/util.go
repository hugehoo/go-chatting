package network

import (
	"github.com/gofiber/fiber/v2"
)

func response(c *fiber.Ctx, status int, res interface{}, data ...string) error {
	return c.Status(status).JSON(data)
}
