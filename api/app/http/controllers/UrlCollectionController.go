package controllers

import (
	"combined-crawler/api/app/models"
	"combined-crawler/api/app/services"
	"github.com/gofiber/fiber/v2"
)

type UrlCollectionController struct {
	Service *services.UrlCollectionService
}

func NewUrlCollectionController(service *services.UrlCollectionService) *UrlCollectionController {
	return &UrlCollectionController{Service: service}
}

func (ctrl *UrlCollectionController) Create(c *fiber.Ctx) error {
	var urlCollection models.UrlCollection
	if err := c.BodyParser(&urlCollection); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := ctrl.Service.Create(&urlCollection)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(urlCollection)
}

func (ctrl *UrlCollectionController) GetByID(c *fiber.Ctx) error {
	collectionID := c.Params("collectionID")
	urlCollection, err := ctrl.Service.GetByID(collectionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(urlCollection)
}

func (ctrl *UrlCollectionController) Update(c *fiber.Ctx) error {
	collectionID := c.Params("collectionID")
	var update map[string]interface{}
	if err := c.BodyParser(&update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := ctrl.Service.Update(collectionID, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func (ctrl *UrlCollectionController) Delete(c *fiber.Ctx) error {
	collectionID := c.Params("collectionID")
	err := ctrl.Service.Delete(collectionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}
