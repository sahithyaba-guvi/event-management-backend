package routes

import (
	login "em_backend/controllers/login"

	"github.com/gofiber/fiber/v2"
)

func Login(app *fiber.App) {
	app.Post("/register", login.Register)
	app.Post("/login", login.Login)
	app.Post("/register-forgot-password-email", login.RegisterForgotPasswordEmail)
}
