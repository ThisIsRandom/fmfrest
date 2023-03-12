package controllers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thisisrandom/fmfrest/database"
	"github.com/thisisrandom/fmfrest/internal"
	"gorm.io/gorm"
)

type TaskController struct {
	db    *database.Connection
	store *internal.CloudinaryStorage
}

type MessageRequestBody struct {
	database.Message `json:"message"`
	TaskID           int `json:"taskId"`
}

func (controller *TaskController) FindAll(c *fiber.Ctx) error {
	var tasks []database.Task

	if roleId := internal.GetRoleFromCtx(c); roleId == 1 {
		userId := internal.GetIdFromCtx(c)

		if tx := controller.db.Instance.Find(&tasks, "user_id = ?", userId); tx.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
		}

		return c.JSON(tasks)
	}

	tx := controller.db.Instance.Find(&tasks)

	if tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}

	return c.JSON(tasks)
}

func (controller *TaskController) Create(c *fiber.Ctx) error {
	var task database.Task
	id := int(internal.GetIdFromCtx(c))
	task.UserID = &id

	if err := c.BodyParser(&task); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if tx := controller.db.Instance.Debug().Create(&task); tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}

	return c.JSON(task)
}

func (controller *TaskController) GetMessages(c *fiber.Ctx) error {
	var tasks []database.Task

	id := int(internal.GetIdFromCtx(c))

	if tx := controller.db.Instance.Debug().Preload("MessageStreams.Messages").Find(&tasks, "user_id = ?", id); tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return c.JSON(tasks)
		}

		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}

	return c.JSON(tasks)
}

func (controller *TaskController) Message(c *fiber.Ctx) error {
	var message database.Message

	id := int(internal.GetIdFromCtx(c))

	taskId, _ := strconv.Atoi(c.Params("id"))

	message.UserID = &id

	if err := c.BodyParser(&message); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "parser")
	}

	if message.MessageStreamID == nil {
		stream := database.MessageStream{
			TaskID: taskId,
			UserID: id,
			Messages: []database.Message{
				message,
			},
		}

		if tx := controller.db.Instance.Debug().Create(&stream); tx.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
		}

		return c.JSON(stream)
	}

	if tx := controller.db.Instance.Debug().Create(&message); tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}

	return c.JSON(message)
}

func (controller *TaskController) GetMessageStream(c *fiber.Ctx) error {
	var ms database.MessageStream

	streamId, _ := strconv.Atoi(c.Params("id"))

	if tx := controller.db.Instance.Debug().Preload("Messages.User").First(&ms, "ID = ?", streamId); tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}

	return c.JSON(ms)
}

func (controller *TaskController) GetMessageStreamsByUserId(c *fiber.Ctx) error {
	var messageStreams []database.MessageStream

	userId := internal.GetIdFromCtx(c)

	if tx := controller.db.Instance.Preload("Messages").Find(&messageStreams, "user_id = ?", userId); tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}

	return c.JSON(messageStreams)
}

func (controller *TaskController) Find(c *fiber.Ctx) error {
	var task database.Task

	taskId, _ := strconv.Atoi(c.Params("id"))

	if tx := controller.db.Instance.First(&task, "ID = ?", taskId); tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}

	return c.JSON(task)
}

/* func (controller *TaskController) Create(c *fiber.Ctx) error {
	var task database.Task
	var images []database.TaskImage

	if err := c.BodyParser(&task); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	form, err := c.MultipartForm()

	if err != nil {
		return err
	}

	if err != nil {
		fiber.NewError(fiber.StatusInternalServerError, "dadad")
	}

	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()

			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "fileopen error")
			}

			pattern := "fmf-*" + filepath.Ext(fileHeader.Filename)

			tmpFile, err := os.CreateTemp("", pattern)

			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "create temp fail")
			}

			_, err = io.Copy(tmpFile, file)

			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "copy to temp fail")
			}

			var d *internal.SaveImageResult
			d, err = controller.store.SaveImage(tmpFile.Name())

			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "save image err")
			}

			taskImage := database.TaskImage{
				Uri: d.Url,
			}

			images = append(images, taskImage)
		}
	}

	return c.JSON(images)
} */

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

	r.Get("/message/streams", c.GetMessageStreamsByUserId)

	r.Post("/message/:id?", c.Message)

	r.Get("/message", c.GetMessages)

	r.Get("/message/:id", c.GetMessageStream)

	r.Get("/:id", c.Find)
}
