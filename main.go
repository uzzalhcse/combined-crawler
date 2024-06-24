package main

import (
	"combined-crawler/handlers/markt"
	"combined-crawler/pkg/ninjacrawler"
)

func main() {

	// register all sites configs to run the crawlers
	ninjacrawler.NewNinjaCrawler().
		//AddSite(kyocera.Crawler()).
		//AddSite(aqua.Crawler()).
		AddSite(markt.Crawler()).
		//AddSite(sandvik.Crawler()).
		Start()
}
