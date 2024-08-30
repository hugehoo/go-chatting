package network

import (
	"bufio"
	"chat_controller/service"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
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
	s := &Server{engine: fiber.New(), service: service, port: port}
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
	t.server.engine.Get("/load-test", t.loadTest)
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

type message struct {
	Name    string    `json:"name"`
	Message string    `json:"message"`
	Room    string    `json:"room"`
	When    time.Time `json:"when"`
}

type LoadTestConfig struct {
	Connections int    `json:"connections"`
	Duration    int    `json:"duration"`
	ServerURL   string `json:"server_url"`
}

func (t *tower) loadTest(ctx *fiber.Ctx) error {
	config := new(LoadTestConfig)
	serverList := t.server.service.ResponseLiveServerList()
	log.Info("⭐️⭐️⭐️", serverList)
	if err := ctx.BodyParser(config); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if config.Connections <= 0 || config.Duration <= 0 || config.ServerURL == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid configuration"})
	}
	var wg sync.WaitGroup
	results := make(chan string, config.Connections)
	for i := 0; i < config.Connections; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := runWebSocketClient(config.ServerURL, config.Duration, id, results, ctx)
			if err != nil {
				results <- fmt.Sprintf("Client %d error: %v", id, err)
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var response []string
	for result := range results {
		response = append(response, result)
	}

	return ctx.JSON(fiber.Map{"results": response})
}

func runWebSocketClient(serverURL string, duration, id int, results chan<- string, ctx *fiber.Ctx) error {
	// Connect to the WebSocket server
	userId := fmt.Sprintf("user_%d", id)
	header := http.Header{
		"Cookie": []string{"auth=" + userId},
		//"room":   []string{fmt.Sprintf("Room%d", id%5)},
	}
	c, resp, err := websocket.DefaultDialer.Dial(serverURL, header)
	if err != nil {
		if resp != nil {
			log.Printf("Client %d received HTTP response: %d %s", id, resp.StatusCode, resp.Status)
			body, _ := ioutil.ReadAll(resp.Body)
			log.Printf("Response body: %s", string(body))
		}
		return fmt.Errorf("dial error: %v", err)
	}
	defer c.Close()

	messagesSent := 0
	messagesReceived := 0

	// Set up done channel for test duration
	done := make(chan bool)
	go func() {
		time.Sleep(time.Duration(duration) * time.Second)
		done <- true
	}()

	// Set up message sending goroutine
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				msg := message{
					Name:    userId,
					Message: getRandomMessage(),
					Room:    fmt.Sprintf("Room%d", id%5),
					When:    time.Now(),
				}
				msgBytes, err := json.Marshal(msg)
				if err != nil {
					log.Printf("Client %d error marshaling message: %v", id, err)
					continue
				}

				if err := c.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
					log.Printf("Client %d error writing message: %v", id, err)
					return
				}
				messagesSent++
				time.Sleep(time.Millisecond * 10) // Adjust this delay as needed
			}
		}
	}()

	// Main loop for receiving messages
	for {
		select {
		case <-done:
			results <- fmt.Sprintf("Client %d completed. Sent: %d, Received: %d", id, messagesSent, messagesReceived)
			return nil
		default:
			_, _, err := c.ReadMessage()
			if err != nil {
				log.Printf("Client %d read error: %v", id, err)
				results <- fmt.Sprintf("Client %d error: %v", id, err)
				return err
			}
			messagesReceived++
		}
	}
}

func getRandomMessage() string {
	file, err := os.Open("sample-chat-data.csv")
	if err != nil {
		return ""
	}
	defer file.Close()

	var messages []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		messages = append(messages, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return ""
	}

	if len(messages) == 0 {
		return ""
	}

	randomIndex := rand.Intn(len(messages))
	return messages[randomIndex]
}
