package controllers

import (
	"combined-crawler/api/app/models"
	"combined-crawler/api/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type SecretCollectionController struct {
	Service *services.SecretCollectionService
}

func NewSecretCollectionController(service *services.SecretCollectionService) *SecretCollectionController {
	return &SecretCollectionController{Service: service}
}

func (ctrl *SecretCollectionController) Index(c *fiber.Ctx) error {
	siteCollections, err := ctrl.Service.GetAllGlobalSecret()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(siteCollections)
}
func (ctrl *SecretCollectionController) Create(c *fiber.Ctx) error {
	var secretCollection models.SiteSecret
	if err := c.BodyParser(&secretCollection); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := ctrl.Service.Create(&secretCollection)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(secretCollection)
}

func (ctrl *SecretCollectionController) GetByID(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	siteCollection, err := ctrl.Service.GetByID(siteID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	//return c.JSON(siteCollection.Secrets)
	// Prepare the data for .env file format
	var envData string
	for key, value := range siteCollection.Secrets {
		envData += fmt.Sprintf("%s=%v\n", key, value)
	}

	// You can then return the data or store it in a file
	return c.SendString(envData)
}

func (ctrl *SecretCollectionController) Update(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	var update map[string]interface{}
	if err := c.BodyParser(&update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := ctrl.Service.Update(siteID, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func (ctrl *SecretCollectionController) Delete(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	err := ctrl.Service.Delete(siteID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}
