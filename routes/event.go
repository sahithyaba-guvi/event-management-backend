package routes

import (
	eventPanel "em_backend/controllers/events"
	"em_backend/library/middleware"

	"github.com/gofiber/fiber/v2"
)

func EventPanel(app *fiber.App) {
	eventApi := app.Group("/event", middleware.AuthenticationMiddleware)

	eventApi.Post("/getAllEvents", eventPanel.GetAllEvents)
}
