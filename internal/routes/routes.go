package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mouradev1/buscacepsgolang/internal/controllers"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/cep/:cep", controllers.GetCep)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "API de Consulta de CEP está funcionando!",
			"version": "1.2.2",
			"status":  "OK",
			"router":  "/cep/:cep",
			"server":  "Golang with Fiber",
		})
	})
}