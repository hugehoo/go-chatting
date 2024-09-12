package types

import (
	"github.com/gofiber/contrib/websocket"
	"log"
	"time"
)

type Client struct {
	Socket *websocket.Conn
	Send   chan *Message
	Room   *RoomSocket
	Name   string `json:"name"`
}

func (c *Client) Read() {
	defer func(Socket *websocket.Conn) {
		err := Socket.Close()
		if err != nil {
			log.Println("Error occurred while closing Socket connection")
		}
	}(c.Socket)

	for {
		var msg *Message
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
