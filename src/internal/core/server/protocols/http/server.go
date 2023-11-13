// server.go

package http

import (
	"errors"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger"
)

type Server struct {
	Address  string     // Address to listen on
	listener *fiber.App // TCP Listener object
	running  bool
}

func NewServer(address string) *Server {
	app := fiber.New(fiber.Config{
		Prefork:                  false,
		StrictRouting:            true,
		CaseSensitive:            true,
		DisableStartupMessage:    true,
		DisableHeaderNormalizing: true,
		EnablePrintRoutes:        true,
		RequestMethods: []string{
			fiber.MethodHead,
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodPatch,
			fiber.MethodDelete,
		},
	})

	app.Use(helmet.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	/*
		app.Use(func(c *fiber.Ctx) error {
			req := &transport.Request{}
			res := &transport.Response{}

			err := router.Resolve(req, res)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("An internal error occurred.")
			}

			b, err := proto.Marshal(res)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("An internal error occurred.")
			}

			return c.Send(b)
		})

		/*
		app.Use(etag.New())
		app.Use(cache.New(cache.Config{
			CacheControl: true,
			Expiration:   config.DEFAULT_CACHE * time.Minute,
			Methods: []string{
				fiber.MethodGet,
				fiber.MethodHead,
			},
		}))
	*/

	return &Server{
		Address:  address,
		listener: app,
	}
}

func (s *Server) Start() error {
	if s.running {
		return errors.New("server already started")
	}
	s.running = true
	logger.Info("server start on " + s.Address + " with pid:" + strconv.Itoa(os.Getpid()))
	return s.listener.Listen(s.Address)
}

func (s *Server) Stop() error {
	if !s.running {
		return errors.New("server already stoped")
	}

	s.running = false
	logger.Info("server stop on " + s.Address)
	return s.listener.Shutdown()
}
