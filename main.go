package main

import (
    "log"
    "rest-service/routes"
    "rest-service/db"

    "github.com/gofiber/fiber/v2"
)

func main() {
    db.Connect()
    db.AutoMigrate()

    app := fiber.New()

	routes.Setup(app)

    // Start the server
    log.Println("Starting server on :4000")
    err := app.Listen(":4000")
    log.Fatal(err)
}
