package controllers

import (
	"combined-crawler/api/app/helper"
	"combined-crawler/api/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

// TestController defines a controller for handling test-related requests
type TestController struct {
	*BaseController
	TestService *services.TestService
}

// NewTestController creates a new instance of the test controller
func NewTestController(testService *services.TestService) *TestController {
	that := NewBaseController()
	return &TestController{
		BaseController: that,
		TestService:    testService,
	}
}

func (that *TestController) Test(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{
		"message": "Hello World",
		"status":  "Success",
	})
}

// GetAllHandler handles the route to get all records
func (that *TestController) GetAllHandler(c *fiber.Ctx) error {
	models, err := that.TestService.GetAllSiteCollections()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(models)
}

func (that *TestController) StartCrawler(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	err := helper.GenerateBinaryBuild(siteID)
	if err != nil {
		return err
	}
	zone := "asia-northeast1-a"
	resp, err := helper.CreateVM(siteID, zone)
	if err != nil {
		return err
	}
	fmt.Println("VM created successfully", resp)

	// TODO: Get the instance name,id and store into db
	return c.JSON(resp)
}
