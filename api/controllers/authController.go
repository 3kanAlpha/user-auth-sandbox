package controllers

import (
	"auth-api/models"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var user models.User

	user.FirstName = "Self"
	user.LastName = "Note"
	user.Email = "selfnote@yahoo.co.jp"
	user.Password = "pass"
	return c.JSON(user)
}
