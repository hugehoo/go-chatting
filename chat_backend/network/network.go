package network

import (
	"chat_backend/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Server struct {
	engine  *fiber.App
	service *service.Service
	port    string
	ip      string
}

func NewServer(service *service.Service, port string) *Server {
	s := &Server{
		engine:  fiber.New(),
		service: service,
		port:    port,
	}

	defaultHeaderOptions := strings.Join([]string{
		fiber.HeaderOrigin,
		fiber.HeaderContentLength,
		fiber.HeaderContentType,
		fiber.HeaderAccessControlAllowHeaders,
		fiber.HeaderAccessControlAllowOrigin,
		fiber.HeaderAuthorization,
		fiber.HeaderXRequestedWith,
		fiber.HeaderExpires}, ",")

	s.engine.Use(cors.New(cors.Config{
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowHeaders:     defaultHeaderOptions,
		ExposeHeaders:    defaultHeaderOptions,
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	registerServer(s)
	return s
}

func (s *Server) setServerInfo() {
	if address, err := net.InterfaceAddrs(); err != nil {
		panic(err.Error())
	} else {
		var ip net.IP
		for _, addr := range address {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					ip = ipNet.IP
					break
				}
			}
		}
		if ip == nil {
			panic("no ip address found")
		} else {
			if err = s.service.ServerSet(ip.String()+s.port, true); err != nil {
				panic(err)
			} else {
				s.ip = ip.String()
			}
			s.service.PublishServerStatusEvent(s.ip+s.port, true)
		}
	}

}

func (s *Server) StartServer() error {

	s.setServerInfo()
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-channel

		address := s.ip + s.port
		if err := s.service.ServerSet(address, false); err != nil {
			log.Println("Failed To Set Server Info When Close", "err", err)
		}
		s.service.PublishServerStatusEvent(address, false)
		os.Exit(1)
	}()
	log.Println("Start Tx Server")
	return s.engine.Listen(s.port)
}
