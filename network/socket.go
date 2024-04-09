package network

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	. "websocket/types"
)

// Http 통신을 Websocket 통신으로 업그레이드 해준다.
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  SocketBufferSize,
	WriteBufferSize: MessageBufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type message struct {
	Name    string
	Message string
	Time    int64
}

type client struct {
	Send   chan *message
	Name   string
	Room   *Room
	Socket *websocket.Conn
}

// Room | chatting 방에 관한 정보
type Room struct {
	Forward chan *message    // 수신되는 메시지 보관하는 값 - 들어오는 메시지를 다른 클라에게 전송
	Join    chan *client     // Socket 이 연결되는 경우에 작동
	Leave   chan *client     // Socket 이 끊어지는 경ㅇ우에 대해 작동
	Clients map[*client]bool // 현재 방에 있는 Client 정보를 저장
}

// room 초기화 method
func NewRoom() *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *client),
		Leave:   make(chan *client),
		Clients: make(map[*client]bool),
	}
}

func (c *client) Read() {
	// client 가 들어오는 메시지를 읽는 함수
	defer c.Socket.Close()
	for {
		var msg *message
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				break
			} else {
				panic(err)
			}
		} else {
			msg.Time = time.Now().Unix()
			msg.Name = c.Name
			c.Room.Forward <- msg
		}
	}
}

func (c *client) Write() {
	// client 가 메시지를 전송하는 함수
	defer c.Socket.Close()
	for msg := range c.Send {
		err := c.Socket.WriteJSON(msg)
		if err != nil {
			panic(err)
		}
	}
}

func (r *Room) RunInit() {
	// Room 에 있는 모든 채널값을 받는 역할
	for { // 무한정으로 반복한다.
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			r.Clients[client] = false
			delete(r.Clients, client)
			close(client.Send)
		case msg := <-r.Forward:
			// 모든 클라이언트들에게 전파를 해줘야한다.
			for client := range r.Clients {
				client.Send <- msg
			}
		}
	}
}

func (r *Room) SocketServe(c *gin.Context) {
	// socket 생성
	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	userCookie, err := c.Request.Cookie("auth")
	if err != nil {
		panic(err)
	}

	// client 가 만들어짐
	client := &client{
		Socket: socket,
		Send:   make(chan *message, MessageBufferSize),
		Room:   r,
		Name:   userCookie.Value,
	}

	r.Join <- client

	// SocketServe 메서드 블럭을 다 돌고 해당 메서드를 벗어날 때 실행된다 : defer
	defer func() { r.Leave <- client }()

	go client.Write()
	client.Read()

}
