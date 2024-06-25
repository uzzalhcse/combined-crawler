package main

import "combined-crawler/pkg/ninjacrawler"

func main() {
	ninjacrawler.NewNinjaCrawler().RunAutoPilot()

	// register all sites configs to run the crawlers
	//ninjacrawler.NewNinjaCrawler().
	//	//AddSite(kyocera.Crawler()).
	//	AddSite(aqua.Crawler()).
	//	Start()
}
