package platformhttp

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

type Server struct {
	app    *fiber.App
	port   string
	listen func(string, ...fiber.ListenConfig) error
}

func NewServer(port string, app *fiber.App) *Server {
	return &Server{
		app:    app,
		port:   port,
		listen: app.Listen,
	}
}

func (s *Server) App() *fiber.App {
	return s.app
}

func (s *Server) Run() error {
	return s.listen(fmt.Sprintf(":%s", s.port))
}
