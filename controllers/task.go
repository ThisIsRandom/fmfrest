package controllers

import (
	"io"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/thisisrandom/fmfrest/database"
	"github.com/thisisrandom/fmfrest/internal"
)

type TaskController struct {
	db    *database.Connection
	store *internal.CloudinaryStorage
}

func (controller *TaskController) FindAll(c *fiber.Ctx) error {
	var tasks []database.Task

	tx := controller.db.Instance.Find(&tasks)

	if tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}

	return c.JSON(tasks)
}

func (controller *TaskController) Create(c *fiber.Ctx) error {
	//var task database.Task
	var images []database.TaskImage

	/* if err := c.BodyParser(&task); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} */

	form, err := c.MultipartForm()

	if err != nil {
		return err
	}

	/* if err != nil {
		fiber.NewError(fiber.StatusInternalServerError, "dadad")
	} */

	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()

			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			pattern := "fmf-*" + filepath.Ext(fileHeader.Filename)

			tmpFile, err := os.CreateTemp("", pattern)

			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			_, err = io.Copy(tmpFile, file)

			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			var d *internal.SaveImageResult
			d, err = controller.store.SaveImage(tmpFile.Name())

			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			taskImage := database.TaskImage{
				Uri: d.Url,
			}

			images = append(images, taskImage)
		}

		return c.JSON(images)
	}

	return c.JSON(fiber.Map{"ok": "ok"})
}

func NewTaskController(db *database.Connection, store *internal.CloudinaryStorage) *TaskController {
	return &TaskController{
		db,
		store,
	}
}

func RegisterTaskController(router fiber.Router, db *database.Connection, store *internal.CloudinaryStorage) {
	r := router.Group("task")

	c := NewTaskController(db, store)

	r.Post("/", c.Create)

	r.Get("/", c.FindAll)
}
