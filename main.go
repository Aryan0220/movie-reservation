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
	config.ConnectDB()

	app := fiber.New()

	routes.Setup(app)
	log.Println(config.GetEnv("PORT"))
	app.Listen(":" + config.GetEnv("PORT"))
}
