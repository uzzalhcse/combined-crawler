package main

import (
	"combined-crawler/api/app/exceptions"
	"combined-crawler/api/bootstrap"
	"combined-crawler/api/routes"
	"combined-crawler/handlers/aqua"
	"combined-crawler/handlers/kojima"
	"combined-crawler/handlers/kyocera"
	"combined-crawler/handlers/midori_anzen"
	"combined-crawler/handlers/osg"
	"combined-crawler/handlers/sumitool"
	"combined-crawler/handlers/yamaya"
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"os"
)

func main() {
	//go startServer()
	//ninjacrawler.NewNinjaCrawler().RunAutoPilot()
	siteName := ""
	if len(os.Args) > 1 {
		firstArgument := os.Args[1]
		siteName = firstArgument
	}
	if siteName == "" {
		fmt.Println("Usage: go run start_vm.go SITE_NAME")
		os.Exit(1)
	}
	// register all sites configs to run the crawlers
	ninjacrawler.NewNinjaCrawler().
		AddSite(kyocera.Crawler()).
		AddSite(aqua.Crawler()).
		AddSite(yamaya.Crawler()).
		AddSite(midori_anzen.Crawler()).
		AddSite(osg.Crawler()).
		AddSite(sumitool.Crawler()).
		AddSite(kojima.Crawler()).
		StartOnly(siteName)
}

func startServer() {
	app := bootstrap.App()
	defer app.CloseDBConnection()
	app.ConnectDB()

	// Register routes
	routes.RegisterRoutes(app.App)

	// Launch the application in a goroutine
	go startApplication(app)

	// Graceful shutdown
	app.GracefulShutdown(func() {
		shutdownApplication(app)
	})
}
func startApplication(app *bootstrap.Application) {
	port := ":" + app.Config.App.Port
	if err := app.Run(port); err != nil {
		exceptions.PanicIfNeeded(err.Error())
	}
}

func shutdownApplication(app *bootstrap.Application) {
	if err := app.Shutdown(); err != nil {
		fmt.Println("Error during shutdown:", err.Error())
	}

	app.CloseDBConnection()
}
