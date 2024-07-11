package osg

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func ProductHandler(crawler *ninjacrawler.Crawler) {

	crawler.CrawlPageDetail([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.ProductDetails,
			OriginCollection: constant.Categories,
			Processor:        handleProduct,
			Preference: ninjacrawler.Preference{
				ValidationRules:     []string{"PageTitle|required"},
				DoNotMarkAsComplete: true,
			},
		},
	})
}
func handleProduct(ctx ninjacrawler.CrawlerContext, fn func([]ninjacrawler.ProductDetailSelector, string)) error {

	items := []ninjacrawler.ProductDetailSelector{}
	ctx.Document.Find("#searchresults ul.productlist li.clearfix").Each(func(i int, s *goquery.Selection) {
		items = append(items, selectProduct(s))
	})

	nextPageUrl := ""
	lastPageUrl := ""
	currentUrl := strings.Split(ctx.UrlCollection.Url, "?")[0]
	ctx.Document.Find("nav.pager.clearfix:nth-child(1) ul li a").Each(func(i int, s *goquery.Selection) {
		rel, exists := s.Attr("rel")
		if exists {
			if rel == "next" {
				href, _ := s.Attr("href")
				nextPageUrl = currentUrl + href
			}
			if rel == "last" {
				href, _ := s.Attr("href")
				lastPageUrl = currentUrl + href
			}
		}
	})

	if ctx.UrlCollection.CurrentPageUrl == lastPageUrl {
		fn(items, "")
		return nil
	} else {
		fn(items, nextPageUrl)
	}
	return nil
}

func selectProduct(selection *goquery.Selection) ninjacrawler.ProductDetailSelector {
	productDetailSelector := ninjacrawler.ProductDetailSelector{
		Jan: "",
		PageTitle: &ninjacrawler.SingleSelector{
			Selector: "title",
		},
		Url: getUrlHandler,
		Images: func(ctx ninjacrawler.CrawlerContext) []string {
			fullUrl := ""
			el := selection.Find("p.thumb img").First()
			if url, ok := el.Attr("src"); ok {
				fullUrl = ctx.App.GetFullUrl(url)
			}
			return []string{fullUrl}
		},
		ProductCodes: func(ctx ninjacrawler.CrawlerContext) []string {
			toolNo := selection.Find("dl dd.toolNo").Text()
			return []string{toolNo}
		},
		Maker: func(ctx ninjacrawler.CrawlerContext) string {
			return "オーエスジー"

		},
		Brand:       "",
		ProductName: "",
		Category:    "",
		Description: "",
		Reviews: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		ItemTypes: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		ItemSizes: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		ItemWeights: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		SingleItemSize:   "",
		SingleItemWeight: "",
		NumOfItems:       "",
		ListPrice:        "",
		SellingPrice:     "",
		Attributes: func(ctx ninjacrawler.CrawlerContext) []ninjacrawler.AttributeItem {
			return getProductAttribute(ctx, selection)
		},
	}
	return productDetailSelector
}
