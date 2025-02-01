package main

import "github.com/gofiber/fiber/v2"

func main() {
	handler := &Handler{}
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Get("/user", handler.GetAllUsers)

	app.Listen(":3000")
}
