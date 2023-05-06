package http

import (
	"encoding/json"
	"fmt"

	"github.com/MohammadNE/PhoneBook/internal/repository"
	"github.com/MohammadNE/PhoneBook/pkg/token"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Server struct {
	config     *Config
	logger     *zap.Logger
	repository repository.Repository
	token      token.Token

	app *fiber.App
}

func New(cfg *Config, log *zap.Logger, repo repository.Repository, token token.Token) *Server {
	server := &Server{config: cfg, logger: log, repository: repo, token: token}

	server.app = fiber.New(fiber.Config{JSONEncoder: json.Marshal, JSONDecoder: json.Unmarshal})

	v1 := server.app.Group("api/v1")

	auth := v1.Group("auth")
	auth.Post("/register", server.register)
	auth.Post("/login", server.login)

	contacts := v1.Group("contacts", server.fetchUserId)
	contacts.Get("/", server.getContacts)
	contacts.Get("/:id", server.getContact)
	contacts.Put("/:id", server.updateContact)
	contacts.Delete("/:id", server.deleteContact)

	return server
}

func (server *Server) Serve() error {
	addr := fmt.Sprintf(":%d", server.config.ListenPort)
	if err := server.app.Listen(addr); err != nil {
		server.logger.Error("error resolving server", zap.Error(err))
		return err
	}
	return nil
}