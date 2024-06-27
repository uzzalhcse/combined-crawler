package controllers

import (
	"combined-crawler/api/app/models"
	"combined-crawler/api/app/services"
	"github.com/gofiber/fiber/v2"
)

type CollectionController struct {
	Service *services.CollectionService
}

func NewCollectionController(service *services.CollectionService) *CollectionController {
	return &CollectionController{Service: service}
}

func (ctrl *CollectionController) Index(c *fiber.Ctx) error {
	collections, err := ctrl.Service.GetAllSiteCollections()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(collections)
}

func (ctrl *CollectionController) Create(c *fiber.Ctx) error {
	var collection models.Collection
	if err := c.BodyParser(&collection); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := ctrl.Service.Create(&collection)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(collection)
}

func (ctrl *CollectionController) GetByID(c *fiber.Ctx) error {
	collectionID := c.Params("collectionID")
	collection, err := ctrl.Service.GetByID(collectionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(collection)
}

func (ctrl *CollectionController) Update(c *fiber.Ctx) error {
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

func (ctrl *CollectionController) Delete(c *fiber.Ctx) error {
	collectionID := c.Params("collectionID")
	err := ctrl.Service.Delete(collectionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}
