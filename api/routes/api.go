package routes

import (
	"combined-crawler/api/app/http/controllers"
	"combined-crawler/api/app/repositories"
	"combined-crawler/api/app/services"
	"combined-crawler/api/bootstrap"
	"github.com/gofiber/fiber/v2"
)

func SetUpApiRoutes(api fiber.Router) {
	// Initialize repositories and services
	repo := repositories.NewRepository(bootstrap.App().DB)
	siteCollectionService := services.NewSiteCollectionService(repo)
	collectionService := services.NewCollectionService(repo)
	urlCollectionService := services.NewUrlCollectionService(repo)
	secretCollectionService := services.NewSecretCollectionService(repo)

	// Initialize controllers
	siteCollectionController := controllers.NewSiteCollectionController(siteCollectionService)
	collectionController := controllers.NewCollectionController(collectionService)
	urlCollectionController := controllers.NewUrlCollectionController(urlCollectionService)
	secretCollectionController := controllers.NewSecretCollectionController(secretCollectionService)

	// Test controller
	testService := services.NewTestService(repo)
	testController := controllers.NewTestController(testService)

	// Define routes
	api.Get("/", testController.Test)
	api.Get("/test", testController.GetAllHandler)
	api.Get("/start-crawler/:SiteID/:zone", testController.StartCrawler)

	// SiteCollection routes
	site := api.Group("/site")
	site.Get("/", siteCollectionController.Index)
	site.Post("/", siteCollectionController.Create)
	site.Get("/:siteID", siteCollectionController.GetByID)
	site.Put("/:siteID", siteCollectionController.Update)
	site.Delete("/:siteID", siteCollectionController.Delete)

	// Collection routes
	collection := api.Group("/collection")
	collection.Get("/", collectionController.Index)
	collection.Post("/", collectionController.Create)
	collection.Get("/:collectionID", collectionController.GetByID)
	collection.Put("/:collectionID", collectionController.Update)
	collection.Delete("/:collectionID", collectionController.Delete)

	// UrlCollection routes
	url := api.Group("/urlcollections")
	url.Post("/", urlCollectionController.Create)
	url.Get("/:collectionID", urlCollectionController.GetByID)
	url.Put("/:collectionID", urlCollectionController.Update)
	url.Delete("/:collectionID", urlCollectionController.Delete)

	// SiteSecretCollection routes
	secret := api.Group("/site-secret")
	secret.Post("/", secretCollectionController.Create)
	secret.Get("/:siteID", secretCollectionController.GetByID)
	secret.Put("/:siteID", secretCollectionController.Update)
	secret.Delete("/:siteID", secretCollectionController.Delete)
}
