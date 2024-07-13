package api

import "github.com/gofiber/fiber/v2"

func HandleGetUsers(c *fiber.Ctx) error {
	return c.JSON("Jack")
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("Jack")
}