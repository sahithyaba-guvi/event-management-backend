package main

import (
	"em_backend/routes"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	// Use CORS middleware
	app.Use(cors.New(cors.Config{
		AllowMethods: "GET,PUT,POST,DELETE",
		AllowOrigins: "*",
	}))

	// Use Logger middleware to print incoming requests
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip} - ${method} ${path} - ${status}\n",
	}))

	// Define routes
	routes.Login(app)
	routes.AdminPanel(app)
	routes.EventPanel(app)
	routes.PaymentPanel(app)

	// Start the server
	fmt.Println(app.Listen(":3001"))
}
