package kyocera

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
		Selector:     "ul.product-list li.product-item,ul.heightLineParent.clearfix li",
		SingleResult: false,
		FindSelector: "a,div dl dt a",
		Attr:         "href",
	}
	crawler.Collection(constant.Categories).CrawlUrls(crawler.GetBaseCollection(), handleCategory("/prdct/tool/category/product"), ninjacrawler.Preference{MarkAsComplete: false})
	crawler.Collection(constant.Other).CrawlUrls(crawler.GetBaseCollection(), handleCategory("/prdct/tool/sgs/"), ninjacrawler.Preference{MarkAsComplete: false})
	crawler.Collection(constant.Products).CrawlUrls(crawler.GetBaseCollection(), handleCategory("/prdct/tool/product/"))

	crawler.Collection(constant.Products).CrawlUrls(constant.Other, productSelector)
	crawler.Collection(constant.Products).CrawlUrls(constant.Categories, productSelector)
}
