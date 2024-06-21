package main

import (
	"combined-crawler/handlers/markt"
	"combined-crawler/pkg/ninjacrawler"
)

func main() {
	ninjacrawler.NewNinjaCrawler().
		//AddSite(kyocera.Crawler()).
		//AddSite(aqua.Crawler()).
		AddSite(markt.Crawler()).
		//AddSite(sandvik.Crawler()).
		Start()
}
