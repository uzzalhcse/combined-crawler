package sandvik

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {

	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        categoryHandler,
			Engine: ninjacrawler.Engine{
				CookieConsent: &ninjacrawler.CookieAction{
					ButtonText:                  "Accept Cookies",
					MustHaveSelectorAfterAction: ".row.mb-6.col-md-10",
				},
			},
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Categories,
			Processor:        productHandler,
		},
	})
}
