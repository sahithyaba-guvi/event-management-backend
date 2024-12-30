package routes

import (
	adminpanel "em_backend/controllers/admin"
	"em_backend/library/middleware"

	"github.com/gofiber/fiber/v2"
)

func AdminPanel(app *fiber.App) {
	adminApi := app.Group("/admin", middleware.AuthenticationMiddlewareForAdmin)

	adminApi.Post("/addEvent", adminpanel.CreateEvent)
}
