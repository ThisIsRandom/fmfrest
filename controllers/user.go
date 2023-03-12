package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/thisisrandom/fmfrest/database"
	"gorm.io/gorm"
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

func (controller *UserController) GetIdFromCtx(c *fiber.Ctx) uint {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["id"].(float64)
	return uint(userId)
}

func (controller *UserController) GetUser(c *fiber.Ctx) error {
	var r database.User

	userId := controller.GetIdFromCtx(c)

	tx := controller.
		db.
		Instance.
		Debug().
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

func (controller *UserController) Update(c *fiber.Ctx) error {

	var user database.User
	user.ID = controller.GetIdFromCtx(c)

	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Couldnt parse **update user**")
	}

	tx := controller.db.Instance.Debug().Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user)

	if tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Couldnt save to db **update user**")
	}

	return c.JSON(user)
}

func RegisterUserController(router fiber.Router, db *database.Connection) {
	c := NewUserController(
		db,
	)

	r := router.Group("user")

	r.Get("/", c.GetUser)
	r.Post("/", c.Update)
}
