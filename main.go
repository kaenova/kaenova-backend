package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/kaenova/kaenova-backend/config"
	livechat "github.com/kaenova/kaenova-backend/service/live_chat"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	cfg := config.MakeConfig()

	// e := echo.New()

	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"http://local.kaenova.my.id:3000"},
	// 	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderUpgrade},
	// }))
	// e.Use(middleware.Logger())
	// // e.Use(middleware.Recover())

	f := fiber.New()

	f.Use(cors.New(cors.Config{
		AllowOrigins: "http://local.kaenova.my.id:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	liveChat := livechat.NewLiveChatSerice(cfg.HCaptchaSecret)
	liveChat.RegisterRoute(f)

	f.Listen(":1323")

	// e.Logger.Fatal(e.Start(":1323"))
}
