package routes

import (
    "github.com/gofiber/fiber/v2"
    "rest-service/handlers"
)

func Setup(app *fiber.App) {

	app.Get("people", handlers.People)
	app.Post("people/add", handlers.Add)
	app.Put("people/modify", handlers.Modify)
    app.Delete("people/delete", handlers.DeleteId)
}
