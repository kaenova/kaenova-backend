package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	livechat "github.com/kaenova/kaenova-backend/service/live_chat"
	"github.com/kaenova/kaenova-backend/service/live_chat/config"
	"github.com/kaenova/kaenova-backend/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	f := fiber.New()

	f.Use(cors.New(cors.Config{
		AllowOrigins: "http://local.kaenova.my.id:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	liveChat := livechat.NewLiveChatSerice(config.Config{
		HCaptchaSecret: utils.EnvOrDefault("HCAPTCHA_SECRET", "0x0000000000000000000000000000000000000000"),
	})
	liveChat.Http.RegisterRoute(f)
	liveChat.Websocket.RegisterWebsocketRoute(f)

	f.Listen(":1323")
}
