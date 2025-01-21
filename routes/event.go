package routes

import (
	eventPanel "em_backend/controllers/events"
	"em_backend/library/middleware"

	"github.com/gofiber/fiber/v2"
)

func EventPanel(app *fiber.App) {
	eventApi := app.Group("/event", middleware.AuthenticationMiddleware)

	eventApi.Post("/getAllEvents", eventPanel.GetAllEvents)
	eventApi.Post("/getEventById", eventPanel.GetEventByID)
	eventApi.Post("/registerEvent", eventPanel.RegisterEvent)
	eventApi.Post("/registration-form", eventPanel.GetRegistrationForm)
	eventApi.Post("/getAllRegistrations", eventPanel.GetRegistrationDetails)
	eventApi.Post("/getQR-ticket", eventPanel.GetTicketQR)
	eventApi.Post("/verify-ticket", eventPanel.VerifyTicket)
}
