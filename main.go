package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	api := app.Group("api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": "ok"})
	})

	app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
