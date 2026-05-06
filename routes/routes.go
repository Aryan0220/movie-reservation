package routes

import (
	"booking-system/handlers"
	"booking-system/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/register", handlers.Register)
	api.Post("/login", handlers.Login)

	movieApi := api.Group("/movie", middleware.Protected)
	// timetableApi := api.Group("/timetable", middleware.Protected)

	movieApi.Post("/add", handlers.AddMovie)
	movieApi.Patch("/update", handlers.UpdateMovie)
	movieApi.Delete("/delete", handlers.DeleteMovie)
	// movieApi.Get("/get", handlers.GetMovies)
	// movieApi.Post("/reserve", handlers.ReserveMovie)
	// movieApi.Delete("/cancel", handlers.CancelReservation)

	// timetableApi.Post("/add", handlers.AddShowTime)
	// timetableApi.Patch("/update", handlers.UpdateShowTime)
	// timetableApi.Get("/capacity", handlers.GetCapacity)
	// timetableApi.Get("/revenue", handlers.GetRevenue)
	// timetableApi.Get("/all/bookings", handlers.GetAllReservations)

	protected := api.Group("/user", middleware.Protected)

	protected.Patch("/promote", handlers.Promote)

	protected.Get("/me", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Protected route"})
	})
}
