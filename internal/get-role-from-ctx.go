package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetRoleFromCtx(c *fiber.Ctx) uint {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["role"].(float64)
	return uint(userId)
}
