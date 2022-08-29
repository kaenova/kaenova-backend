package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kaenova/kaenova-backend/service/live_chat/config"
	"github.com/kaenova/kaenova-backend/service/live_chat/repository"
)

type HttpService struct {
	R   *repository.Repository
	Cfg *config.Config
}

func NewHttpService(r *repository.Repository, c *config.Config) HttpService {
	return HttpService{
		R:   r,
		Cfg: c,
	}
}

func (h *HttpService) RegisterRoute(e *fiber.App) {
	g := e.Group("/livechat")
	g.Get("/hello", h.helloWorld)
	g.Post("/register", h.registerUser)
	g.Get("/chat", h.getAllChat)
}
