package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	live "github.com/kaenova/kaenova-backend/service/live"
	liveCfg "github.com/kaenova/kaenova-backend/service/live/config"
	livechat "github.com/kaenova/kaenova-backend/service/live_chat"
	livechatCfg "github.com/kaenova/kaenova-backend/service/live_chat/config"
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

	f.Use(logger.New())

	// Live Chat Serivce
	liveChat := livechat.NewLiveChatSerice(livechatCfg.Config{
		HCaptchaSecret: utils.EnvOrDefault("HCAPTCHA_SECRET", "0x0000000000000000000000000000000000000000"),
	})
	liveChat.RegisterHttpRoute(f)
	liveChat.RegisterWebsocketRoute(f)

	// Live Service
	live := live.NewLiveChatSerice(liveCfg.Config{
		LiveKeyPassword: utils.EnvOrDefault("LIVE_PASSWORD", "changeme"),
	})
	live.RegisterHttpRoute(f)

	f.Listen(":1323")
}
