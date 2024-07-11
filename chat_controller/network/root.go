package network

import (
	"chat_controller/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	// Logger middleware
	s.engine.Use(logger.New())

	// Recover middleware
	s.engine.Use(recover.New())

	// CORS middleware
	s.engine.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin,Content-Length,Content-Type,Access-Control-Allow-Headers,Access-Control-Allow-Origin,Authorization,X-Requested-With,Expires",
		ExposeHeaders:    "Origin,Content-Length,Content-Type,Access-Control-Allow-Headers,Access-Control-Allow-Origin,Authorization,X-Requested-With,Expires",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	t := &tower{server: s}
	t.server.engine.Get("/ping", t.healthCheck)
	t.server.engine.Get("/server-list", t.serverList) // 요 표현식 이해안됨, t.serverList() 도 아니고 그냥 t.serverList ?
	return s
}

func (t *tower) serverList(ctx *fiber.Ctx) error {
	return response(ctx, http.StatusOK, t.server.service.ResponseLiveServerList())
}

func (t *tower) healthCheck(ctx *fiber.Ctx) error {
	return ctx.SendString("헬스체크 준비 할 완료~")
}

func (s *Server) Start() error {
	log.Println("Start Controller Server on port", s.port)
	err := s.engine.Listen(s.port)
	if err != nil {
		panic(err)
	} else {
		return nil
	}
}
