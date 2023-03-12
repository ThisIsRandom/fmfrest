package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetIdFromCtx(c *fiber.Ctx) uint {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["id"].(float64)
	return uint(userId)
}
