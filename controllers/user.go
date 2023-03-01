package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/thisisrandom/fmfrest/database"
)

type UserController struct {
	db *database.Connection
}

func (controller *UserController) CreateRole(c *fiber.Ctx) error {
	var role database.Role

	if err := c.BodyParser(&role); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if tx := controller.db.Instance.Create(&role); tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}

	return c.JSON(role)
}

func (controller *UserController) GetUser(c *fiber.Ctx) error {
	var r database.User

	user := c.Locals("user").(*jwt.Token)

	claims := user.Claims.(jwt.MapClaims)
	userId := claims["id"].(float64)

	tx := controller.
		db.
		Instance.
		Preload("Profile.Role").
		Preload("ContactInformation").
		First(&r, "id = ?", userId)

	if tx.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, tx.Error.Error())
	}

	return c.JSON(r)
}

func NewUserController(db *database.Connection) *UserController {
	return &UserController{
		db,
	}
}

func RegisterUserController(router fiber.Router, db *database.Connection) {
	c := NewUserController(
		db,
	)

	r := router.Group("user")

	r.Get("/", c.GetUser)

}
