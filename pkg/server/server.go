package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gjg-redis-go/pkg/config"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	RegisterRoutes(ctx *fiber.App)
}

type Server interface {
	GetFiberInstance() *fiber.App
	Start() error
	Shutdown() error
	RegisterRoutes()
}

type server struct {
	serverPort string
	fiber      *fiber.App
	handlers   []Handler
}

func NewServer(config *config.Config, handlers []Handler) Server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})
	serverPort := config.ServerPort

	return &server{
		serverPort: serverPort,
		fiber:      app,
		handlers:   handlers,
	}
}

func (server *server) Start() error {
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdownChannel
		_ = server.fiber.Shutdown()
	}()
	server.RegisterRoutes()

	serverAddress := fmt.Sprintf(":%s", server.serverPort)
	return server.fiber.Listen(serverAddress)
}

func (server *server) Shutdown() error {
	return server.fiber.Shutdown()
}

func (server *server) GetFiberInstance() *fiber.App {
	return server.fiber
}

func (server *server) RegisterRoutes() {
	for _, handler := range server.handlers {
		handler.RegisterRoutes(server.fiber)
	}
}
