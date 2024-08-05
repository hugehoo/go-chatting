package network

import "github.com/gofiber/fiber/v2"

func response(c *fiber.Ctx, status int, res interface{}, data ...string) {
	c.Status(status).JSON(data)
}
