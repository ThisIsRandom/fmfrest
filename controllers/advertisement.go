package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thisisrandom/fmfrest/database"
)

type AdvertisementController struct {
	db *database.Connection
}

type DeleteRequest struct {
	Id int `json:"id"`
}

func (controller *AdvertisementController) Create(c *fiber.Ctx) error {
	var advertisement database.Advertisement

	if err := c.BodyParser(&advertisement); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "couldnt parse advertisement")
	}

	if tx := controller.db.Instance.Create(&advertisement); tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}

	_, err := c.FormFile("file")

	if err != nil {
		fiber.NewError(404, err.Error())
	}

	return c.JSON(advertisement)
}

func (controller *AdvertisementController) Delete(c *fiber.Ctx) error {
	var advertisement database.Advertisement

	id, err := strconv.Atoi(c.Params("id"))

	advertisement.ID = uint(id)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if tx := controller.db.Instance.Delete(&advertisement); tx.Error != nil {
		return fiber.NewError(fiber.StatusBadRequest, tx.Error.Error())
	}

	return c.SendStatus(200)
}

func (controller *AdvertisementController) FindAll(c *fiber.Ctx) error {
	var response []database.Advertisement

	if tx := controller.db.Instance.Find(&response); tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}

	return c.JSON(response)
}

func NewAdvertisementController(db *database.Connection) *AdvertisementController {
	return &AdvertisementController{
		db,
	}
}

func RegisterAdvertisementController(router fiber.Router, db *database.Connection) {
	r := router.Group("advertisement")

	c := NewAdvertisementController(db)

	r.Get("/", c.FindAll)
	r.Post("/", c.Create)
	r.Delete("/:id<int>", c.Delete)
}
