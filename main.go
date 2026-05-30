package main

import (
	"booking-system/config"
	"booking-system/routes"
	"log"
	"github.com/gofiber/fiber/v2"
)

func main() {
	val := config.LoadEnv()
	log.Println(val)
	if err := config.ConnectDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	app := fiber.New()

	routes.Setup(app)
	log.Println(config.GetEnv("PORT"))
	if err := app.Listen(":" + config.GetEnv("PORT")); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
