package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/thisisrandom/fmfrest/database"
	"github.com/thisisrandom/fmfrest/internal"
)

type AuthController struct {
	db *database.Connection
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	LoginRequest
	Name string `json:"name"`
}

func (controller *AuthController) Login(c *fiber.Ctx) error {
	var user database.User
	var data LoginRequest

	if err := c.BodyParser(&data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Couldnt parse body")
	}

	tx := controller.
		db.
		Instance.
		Preload("Profile.Role").
		Preload("ContactInformation").
		First(&user, "email = ?", data.Email)

	if tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "db err")
	}

	if tx.RowsAffected < 1 {
		return fiber.NewError(fiber.StatusNotFound, "Use does not exist")
	}

	if _, err := internal.HashCompare(data.Password, *user.Password); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Username or password doesnt match")
	}

	claims := jwt.MapClaims{
		"id":   int(user.ID),
		"role": int(*user.Profile.RoleID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, _ := token.SignedString([]byte("secret"))

	return c.JSON(fiber.Map{
		"user":  user,
		"token": t,
	})
}

func (controller *AuthController) Register(c *fiber.Ctx) error {
	var data RegisterRequest

	if err := c.BodyParser(&data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Couldnt parse body")
	}

	hashedPw, err := internal.HashPassword(data.Password)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Hash")
	}

	m := database.User{
		Email:    &data.Email,
		Password: &hashedPw,
		Profile: database.Profile{
			Name: data.Name,
		},
	}

	if tx := controller.db.Instance.Create(&m); tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "cant save to db")
	}

	return c.JSON(m)
}

func NewAuthController(db *database.Connection) *AuthController {
	return &AuthController{
		db,
	}
}

func RegisterAuthController(router fiber.Router, db *database.Connection) {
	r := router.Group("auth")

	c := NewAuthController(
		db,
	)

	r.Post("/login", c.Login)
	r.Post("/register", c.Register)
}
