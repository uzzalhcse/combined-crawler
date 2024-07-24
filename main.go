package main

import (
	"combined-crawler/handlers/aqua"
	"combined-crawler/handlers/kojima"
	"combined-crawler/handlers/kyocera"
	"combined-crawler/handlers/midori_anzen"
	"combined-crawler/handlers/osg"
	"combined-crawler/handlers/sumitool"
	"combined-crawler/handlers/yamaya"
	"combined-crawler/pkg/generic_crawler"
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"os"
)

func main() {
	//ninjacrawler.NewNinjaCrawler().RunAutoPilot()
	//siteRegistry()
	//PluginGenerate()
	generic_crawler.NewGenericCrawler().RunAutoPilot()
}
func siteRegistry() {
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
