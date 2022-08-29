package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kaenova/kaenova-backend/service/live/config"
	"github.com/kaenova/kaenova-backend/service/live/repository"
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
	g := e.Group("/live")
	g.Get("/hello", h.helloWorld)
	g.Get("/status", h.getStatus)
	g.Post("/golive", h.goLive)
	g.Post("/goofline", h.goOffline)
}
