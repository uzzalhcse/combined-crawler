package kyocera

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {

	handleCategory := func(keyword string) ninjacrawler.UrlSelector {
		return ninjacrawler.UrlSelector{
			Selector:     "div.index.clearfix ul.clearfix li",
			FindSelector: "a",
			Attr:         "href",
			Handler: func(urlCollection ninjacrawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {
				if strings.Contains(fullUrl, keyword) {
					return fullUrl, nil
				}
				return "", nil
			},
		}
	}
	productSelector := ninjacrawler.UrlSelector{
		Selector:     "ul.heightLineParent.clearfix li",
		FindSelector: "dl dt a",
		Attr:         "href",
	}
	productOtherSelector := ninjacrawler.UrlSelector{
		Selector:     "ul.product-list li.product-item",
		FindSelector: "a",
		Attr:         "href",
	}
	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        handleCategory("/prdct/tool/category/product/"),
			Preference:       ninjacrawler.Preference{DoNotMarkAsComplete: false},
		},
		{
			Entity:           constant.Other,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        handleCategory("/prdct/tool/sgs/"),
			Preference:       ninjacrawler.Preference{DoNotMarkAsComplete: false},
		},
		{
			Entity:           constant.Products,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        handleCategory("/prdct/tool/product/"),
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Categories,
			Processor:        productSelector,
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Other,
			Processor:        productOtherSelector,
		},
	})
}
