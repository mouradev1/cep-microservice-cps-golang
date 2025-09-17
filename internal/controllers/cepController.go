package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mouradev1/buscacepsgolang/internal/services"
)

func GetCep(c *fiber.Ctx) error {
    cep := c.Params("cep")
	log.Println("Received CEP:", cep)
    result, status, err := services.GetCepDataService(c, cep)
    if err != nil {
        return c.Status(status).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(result)
}