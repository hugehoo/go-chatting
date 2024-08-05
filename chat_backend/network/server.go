package network

import (
	"chat_backend/types"
	"github.com/gofiber/fiber/v2"
)

type api struct {
	server *Server
}

func registerServer(server *Server) {
	a := &api{server: server}
	server.engine.Get("/room-list", a.getRooms)
	server.engine.Get("/room", a.room)
	server.engine.Get("/enter-room", a.enterRoom)
	server.engine.Post("/make-room", a.makeRoom)

	r := NewRoom(server.service)
	go r.Run()

	server.engine.Get("/room-chat", r.ServeHTTP)
}

func (a *api) getRooms(c *fiber.Ctx) error {
	res, err := a.server.service.RoomList()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve room list",
		})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (a *api) makeRoom(c *fiber.Ctx) error {
	var req types.BodyRoomReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := a.server.service.MakeRoom(req.Name); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create room",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success",
	})
}

func (a *api) room(c *fiber.Ctx) error {
	var req types.FormRoomReq

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"error": "Invalid request Form",
			})
	}
	if err := c.BodyParser(req.Name); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"error": "Failed to get Room",
			})
	}
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "Success"})
}

func (a *api) enterRoom(c *fiber.Ctx) error {
	var req types.FormRoomReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid request Form"})
	}
	if err := c.BodyParser(req.Name); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "Failed to enter Room"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "enter room",
	})
}
