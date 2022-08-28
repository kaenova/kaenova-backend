package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kaenova/kaenova-backend/config"
	livechat "github.com/kaenova/kaenova-backend/service/live_chat"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	cfg := config.MakeConfig()

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://local.kaenova.my.id:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderUpgrade},
	}))

	liveChat := livechat.NewLiveChatSerice(cfg.HCaptchaSecret)
	liveChat.RegisterEchoRoute(e)

	e.Logger.Fatal(e.Start("0.0.0.0:1323"))
}
