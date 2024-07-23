package network

import (
	"chat_backend/service"
	"github.com/gorilla/websocket"
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

func (c *Client) Write() {
	defer c.Socket.Close()
	for msg := range c.Send {
		if err := c.Socket.WriteJSON(msg); err != nil {
			return
		}
	}
}
