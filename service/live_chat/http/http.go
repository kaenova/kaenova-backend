package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kaenova/kaenova-backend/service/live_chat/config"
	"github.com/kaenova/kaenova-backend/service/live_chat/repository"
)

type HttpService struct {
	R   repository.RepositoryI
	Cfg *config.Config
}

type HttpServiceI interface {
	RegisterHttpRoute(e *fiber.App)
}

func NewHttpService(r repository.RepositoryI, c *config.Config) HttpServiceI {
	return &HttpService{
		R:   r,
		Cfg: c,
	}
}

func (h *HttpService) RegisterHttpRoute(e *fiber.App) {
	g := e.Group("/livechat")
	g.Get("/hello", h.helloWorld)
	g.Post("/register", h.registerUser)
	g.Get("/chat", h.getAllChat)
}
