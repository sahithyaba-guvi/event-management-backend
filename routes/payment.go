package routes

import (
	paymentPanel "em_backend/controllers/payment"
	"em_backend/library/middleware"

	"github.com/gofiber/fiber/v2"
)

func PaymentPanel(app *fiber.App) {
	paymentApi := app.Group("/payment", middleware.AuthenticationMiddleware)

	paymentApi.Post("/create-order", paymentPanel.CreateOrderHandler)
	paymentApi.Post("/verify-payment", paymentPanel.VerifyPaymentHandler)
}
