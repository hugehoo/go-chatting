package network

import (
	"chat_controller/service"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
)

type Server struct {
	engine  *fiber.App
	service *service.Service
	port    string
}

type tower struct {
	server *Server
}

func NewNetwork(service *service.Service, port string) *Server {
	s := &Server{engine: fiber.New(), port: port}
	//s.engine.Use(fiber.Logger()) // fiber logger 찾아야함
	//s.engine.Config(fiber.Recovery) // fiber recovery 찾아야함

	t := &tower{server: s}
	t.server.engine.Get("/server-list", t.serverList) // 요 표현식 이해안됨, t.serverList() 도 아니고 그냥 t.serverList ?
	return s
}

func (t *tower) serverList(ctx *fiber.Ctx) error {
	return response(ctx, http.StatusOK, t.server.service.ResponseLiveServerList())
}

func (s *Server) Start() error {
	log.Println("Start Controller Server")
	app := fiber.New()
	err := app.Listen(s.port)
	if err != nil {
		panic(err)
	} else {
		return nil
	}
}
