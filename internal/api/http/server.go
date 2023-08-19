package http

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadne/phone-book/internal/repository"
	"github.com/mohammadne/phone-book/pkg/token"
	"go.uber.org/zap"
)

type Server struct {
	logger     *zap.Logger
	repository repository.Repository
	token      token.Token

	managmentApp *fiber.App
	clientApp    *fiber.App
}

func New(log *zap.Logger, repo repository.Repository, token token.Token) *Server {
	server := &Server{logger: log, repository: repo, token: token}

	// Managment Endpoints

	server.managmentApp = fiber.New(fiber.Config{JSONEncoder: json.Marshal, JSONDecoder: json.Unmarshal})

	server.managmentApp.Get("/healthz/liveness", server.liveness)
	server.managmentApp.Get("/healthz/readiness", server.readiness)

	// Client Endpoints

	server.clientApp = fiber.New(fiber.Config{JSONEncoder: json.Marshal, JSONDecoder: json.Unmarshal})

	v1 := server.clientApp.Group("api/v1")

	auth := v1.Group("auth")
	auth.Post("/register", server.register)
	auth.Post("/login", server.login)

	contacts := v1.Group("contacts", server.fetchUserId)
	contacts.Get("/", server.getContacts)
	contacts.Post("/", server.createContact)
	contacts.Get("/:id", server.getContact)
	contacts.Put("/:id", server.updateContact)
	contacts.Delete("/:id", server.deleteContact)

	return server
}

func (server *Server) Serve() {
	go func() {
		err := server.managmentApp.Listen(":8080")
		server.logger.Fatal("error resolving managment server", zap.Error(err))
	}()

	go func() {
		err := server.clientApp.Listen(":8081")
		server.logger.Fatal("error resolving client server", zap.Error(err))
	}()
}
