package http

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/kaenova/kaenova-backend/service/live/repository"
)

var typeValidator = validator.New()

func (h *HttpService) helloWorld(c *fiber.Ctx) error {
	return c.SendString("hello from live service")
}

func (h *HttpService) getStatus(c *fiber.Ctx) error {
	return c.JSON(h.R.GetLiveState())
}

type handleOnline struct {
	Title    string `json:"title" xml:"title" form:"title" validate:"required"`
	Password string `json:"password" xml:"password" form:"password" validate:"required"`
}

func (h *HttpService) goLive(c *fiber.Ctx) error {
	var req handleOnline

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	log.Println(req)

	err := typeValidator.Struct(req)
	if err != nil {
		return err
	}

	if req.Password != h.Cfg.LiveKeyPassword {
		return c.SendStatus(400)
	}

	h.R.GoLive(repository.LiveOption{Title: req.Title})

	return c.SendStatus(200)
}

type handleOffline struct {
	Password string `json:"password" xml:"password" form:"password" validate:"required"`
}

func (h *HttpService) goOffline(c *fiber.Ctx) error {
	var req handleOffline

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	log.Println(req)

	err := typeValidator.Struct(req)
	if err != nil {
		return err
	}

	if req.Password != h.Cfg.LiveKeyPassword {
		return c.SendStatus(400)
	}

	h.R.GoOffline()

	return c.SendStatus(200)
}
