package controllers

import (
	"combined-crawler/api/app/models"
	"combined-crawler/api/app/services"
	"github.com/gofiber/fiber/v2"
)

type SiteCollectionController struct {
	Service *services.SiteCollectionService
}

func NewSiteCollectionController(service *services.SiteCollectionService) *SiteCollectionController {
	return &SiteCollectionController{Service: service}
}

func (ctrl *SiteCollectionController) Index(c *fiber.Ctx) error {
	siteCollections, err := ctrl.Service.GetAllSiteCollections()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(siteCollections)
}
func (ctrl *SiteCollectionController) Create(c *fiber.Ctx) error {
	var siteCollection models.SiteCollection
	if err := c.BodyParser(&siteCollection); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := ctrl.Service.Create(&siteCollection)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(siteCollection)
}

func (ctrl *SiteCollectionController) GetByID(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	siteCollection, err := ctrl.Service.GetByID(siteID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(siteCollection)
}

func (ctrl *SiteCollectionController) Update(c *fiber.Ctx) error {
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

func (ctrl *SiteCollectionController) Delete(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	err := ctrl.Service.Delete(siteID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}
