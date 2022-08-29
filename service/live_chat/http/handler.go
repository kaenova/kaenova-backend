package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/kaenova/kaenova-backend/service/live_chat/model"
	"github.com/kaenova/kaenova-backend/utils"
)

var typeValidator = validator.New()

func (h *HttpService) helloWorld(c *fiber.Ctx) error {
	return c.SendString("hello from livechat service")
}

func (h *HttpService) getAllChat(c *fiber.Ctx) error {
	return c.JSON(h.R.Messages)

}

type handleRegisterUser struct {
	Name         string `json:"name" validate:"required"`
	HCaptchaCode string `json:"hcaptcha" validate:"required"`
}

func (h *HttpService) registerUser(c *fiber.Ctx) error {
	var req handleRegisterUser

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	err := typeValidator.Struct(req)
	if err != nil {
		return err
	}

	// Check google captcha
	valid, err := utils.VerifyHcaptcha(utils.HCaptchaReqeust{
		Secret:   h.Cfg.HCaptchaSecret,
		Response: req.HCaptchaCode,
	})

	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if !valid {
		return c.Status(http.StatusBadRequest).SendString("not a valid human request")
	}

	user := model.CreateUser(req.Name)
	h.R.AddAuthenticatedUser(user)

	return c.SendString(user.ID)
}
