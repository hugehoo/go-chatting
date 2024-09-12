package types

import (
	service "chat_backend/service"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
)

type RoomSocket struct {
	ID      string
	Clients map[*Client]bool
	Join    chan *Client
	Leave   chan *Client
	Forward chan *Message
	Service *service.Service
}

const (
	messageBufferSize = 256
)

func NewRoomSocket(roomId string) *RoomSocket {
	return &RoomSocket{
		ID:      roomId,
		Clients: make(map[*Client]bool),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Forward: make(chan *Message),
	}
}

func (r *RoomSocket) Run(service *service.Service) {
	for {
		select {
		case client := <-r.Join:
			log.Println("JOIN")
			r.Clients[client] = true
		case client := <-r.Leave:
			log.Println("LEAVE")
			delete(r.Clients, client)
			close(client.Send)
		case msg := <-r.Forward:
			log.Println("FORWARD", msg)
			log.Println("r.service", service)
			go service.InsertChatting(msg.Name, msg.Message, msg.Room)
			for client := range r.Clients {
				client.Send <- msg
			}
		}
	}
}

func (r *RoomSocket) ServeHTTP(c *fiber.Ctx) error {
	if r == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "RoomSocket is not initialized")
	}

	authCookie := c.Cookies("auth")
	if authCookie == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Auth cookie is missing")
	}

	return websocket.New(func(conn *websocket.Conn) {
		if conn == nil {
			log.Println("WebSocket connection failed")
			return
		}

		client := &Client{
			Socket: conn,
			Send:   make(chan *Message, messageBufferSize),
			Room:   r,
			Name:   authCookie,
		}

		if client == nil {
			log.Println("Failed to create client")
			return
		}

		r.Join <- client

		// Defer leaving the room
		// 또한 defer 를 통해서 client 가 끝날 떄를 대비하여 퇴장하는 작업을 연기한다.
		defer func() {
			if r != nil && client != nil {
				r.Leave <- client
			}
		}()

		// 이 후 고루틴을 통해서 write 를 실행 시킨다.
		go client.Write()

		// Read messages (this will block until the connection is closed)
		// 이 후 메인 루틴에서 read 를 실행함으로써 해당 요청을 닫는것을 차단 -> 연결을 활성화 시키는 것이다. 채널을 활용하여
		client.Read()
	})(c)
}
