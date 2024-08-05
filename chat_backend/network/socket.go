package network

import (
	"chat_backend/service"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

type Room struct {
	Forward chan *message
	Join    chan *Client
	Leave   chan *Client
	Clients map[*Client]bool
	service *service.Service
}

type Client struct {
	Socket *websocket.Conn
	Send   chan *message
	Room   *Room
	Name   string `json:"name"`
}

type message struct {
	Name    string    `json:"name"`
	Message string    `json:"message"`
	Room    string    `json:"room"`
	When    time.Time `json:"when"`
}

const (
	SocketBufferSize  = 1024
	messageBufferSize = 256
)

func NewRoom(service *service.Service) *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
		service: service,
	}
}

func (c *Client) Read() {
	defer c.Socket.Close()
	for {
		var msg *message
		if err := c.Socket.ReadJSON(&msg); err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.Name
		c.Room.Forward <- msg // what is Forward for?
	}
}

// Write : 채팅방 모든 클라이언트 에게 메시지를 전송(write)
func (c *Client) Write() {
	defer c.Socket.Close()
	for msg := range c.Send {
		if err := c.Socket.WriteJSON(msg); err != nil {
			return
		}
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			delete(r.Clients, client)
			close(client.Send)
		case msg := <-r.Forward:
			go r.service.InsertChatting(msg.Name, msg.Message, msg.Room)
			for client := range r.Clients {
				client.Send <- msg
			}
		}
	}
}

func (r *Room) ServeHTTP(c *fiber.Ctx) error {
	return websocket.New(func(conn *websocket.Conn) {
		authCookie := c.Cookies("auth")
		if authCookie == "" {
			log.Println("auth cookie is failed")
			return
		}

		// Create client and join the room
		client := &Client{
			Socket: conn,
			Send:   make(chan *message, messageBufferSize),
			Room:   r,
			Name:   authCookie,
		}
		r.Join <- client

		// Defer leaving the room
		// 또한 defer 를 통해서 client 가 끝날 떄를 대비하여 퇴장하는 작업을 연기한다.
		defer func() { r.Leave <- client }()

		// 이 후 고루틴을 통해서 write 를 실행 시킨다.
		go client.Write()

		// Read messages (this will block until the connection is closed)
		// 이 후 메인 루틴에서 read를 실행함으로써 해당 요청을 닫는것을 차단 -> 연결을 활성화 시키는 것이다. 채널을 활용하여
		client.Read()
	})(c)
}
