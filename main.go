package main

import (
	"combined-crawler/handlers/aqua"
	"combined-crawler/pkg/ninjacrawler"
)

func main() {
	//ninjacrawler.NewNinjaCrawler().RunAutoPilot()
	siteRegistry()
	//PluginGenerate()
	//generic_crawler.NewGenericCrawler().RunAutoPilot()
}
func siteRegistry() {
	// register all sites configs to run the crawlers
	ninjacrawler.NewNinjaCrawler().
		//AddSite(kyocera.Crawler()).
		AddSite(aqua.Crawler()).
		//AddSite(yamaya.Crawler()).
		//AddSite(midori_anzen.Crawler()).
		//AddSite(osg.Crawler()).
		//AddSite(sumitool.Crawler()).
		//AddSite(kojima.Crawler()).
		//AddSite(panasonic.Crawler()).
		//AddSite(suntory.Crawler()).
		//AddSite(markt.Crawler()).
		//AddSite(as1.Crawler()).
		//AddSite(sony.Crawler()).
		//AddSite(kitamura.Crawler()).
		Start()
}
