package controller

import (
	"github.com/gofiber/fiber/v2"
)

// Root
// @Summary Display welcome message.
// @Accept */*
// @Produce plain
// @Success 200 {string} string "Welcome to qa api"
// @Router / [get]
func RootController(c *fiber.Ctx) error {
	return c.SendString("Welcome to qa api")
}
