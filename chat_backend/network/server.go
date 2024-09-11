package network

import (
	"chat_backend/types"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
)

type api struct {
	server      *Server
	RoomManager *RoomManager
}

func registerServer(server *Server) {
	a := &api{
		server: server,
		RoomManager: &RoomManager{
			rooms: make(map[string]*Room),
		},
	}
	room := &Room{
		ID:      "roomId",
		Clients: make(map[*Client]bool),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Forward: make(chan *message),
	}
	a.RoomManager.rooms["roomId"] = room

	service := server.service
	go room.Run(service) // 여기서 Run() 메서드를 고루틴으로 실행
	server.engine.Get("/room-list", a.getRooms)
	server.engine.Get("/room", a.room)
	server.engine.Get("/enter-room", a.enterRoom)
	//server.engine.Post("/make-room", a.makeRoom)

	//r := NewRoom(server.service)
	//go r.Run()
	//server.engine.Get("/room-chat", r.ServeHTTP)
	server.engine.Get("/room-chat/:id", websocket.New(a.handleRoomChat))
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
	//var req types.BodyRoomReq
	//if err := c.BodyParser(&req); err != nil {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	//		"error": "Invalid request body",
	//	})
	//}
	//
	//if err := a.server.service.MakeRoom(req.Name); err != nil {
	//	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	//		"error": "Failed to create room",
	//	})
	//}
	//
	roomID := generateUniqueRoomID() // 고유한 방 ID 생성
	rm := a.RoomManager
	log.Println(roomID)
	log.Println(rm)
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if _, exists := rm.rooms[roomID]; exists {
		return nil
	}

	room := &Room{
		ID:      roomID,
		Clients: make(map[*Client]bool),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Forward: make(chan *message),
	}
	rm.rooms[roomID] = room
	//if err := rm.server.service.MakeRoom(room.ID); err != nil {
	//	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	//		"error": "Failed to create room",
	//	})
	//}
	//go room.Run()
	if room == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to create room",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"roomID":  roomID,
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

func generateUniqueRoomID() string {
	return uuid.New().String()
}

func (rm *RoomManager) GetRoom(id string) *Room {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	return rm.rooms[id]
}

func (a *api) handleRoomChat(c *websocket.Conn) {
	roomID := c.Params("id")
	room := a.RoomManager.GetRoom(roomID)
	if room == nil {
		c.Close()
		return
	}
	client := &Client{
		Socket: c,
		Send:   make(chan *message, 256),
	}

	log.Println("before join")
	room.Join <- client
	log.Println("after join")

	go func() {
		defer func() {
			room.Leave <- client
			c.Close()
		}()

		for {
			var msg message
			err := c.ReadJSON(&msg)
			if err != nil {
				log.Println("break", &msg)
				break
			}
			msg.Room = roomID
			log.Println("in go", &msg)
			room.Forward <- &msg
		}
	}()

	for msg := range client.Send {
		log.Println("msg", msg)
		err := c.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
