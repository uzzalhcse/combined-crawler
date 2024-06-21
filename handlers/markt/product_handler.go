package markt

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"strings"
)

func ProductHandler(crawler *ninjacrawler.Crawler) {

	crawler.ProductDetailSelector = ninjacrawler.ProductDetailSelector{
		Jan: func(ctx ninjacrawler.CrawlerContext) string {
			res := ctx.Document.Find("div.p-product-detail > div.p-product-detail__code > dl:nth-child(1) > span").Text()
			if res == "" {
				ctx.App.Logger.Html(ctx.Page, "Empty product Jan.")
			}
			return res
		},
		PageTitle: &ninjacrawler.SingleSelector{
			Selector: "title",
		},
		Url: func(ctx ninjacrawler.CrawlerContext) string { return ctx.UrlCollection.Url },
		Images: &ninjacrawler.MultiSelectors{
			Selectors: []ninjacrawler.Selector{
				{Query: ".u-image-watcher.u-image-watcher__select span img", Attr: "src"},
			},
		},
		ProductCodes: func(ctx ninjacrawler.CrawlerContext) []string { return []string{} },
		Maker:        &ninjacrawler.SingleSelector{Selector: ".p-product-detail__maker span"},
		Brand:        &ninjacrawler.SingleSelector{Selector: "div.u-hidden-sp > ul > li:nth-child(2) > a"},
		ProductName:  &ninjacrawler.SingleSelector{Selector: "div.p-product-detail > div.p-product-detail__maker > p"},
		Category:     &ninjacrawler.SingleSelector{Selector: "div.u-hidden-sp > ul > li:nth-child(1) > a"},
		Description: func(ctx ninjacrawler.CrawlerContext) string {
			res := strings.TrimSpace(ctx.Document.Find("div.p-product-detail > div.p-product-detail__description > p").Text())
			return res
		},
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
		SellingPrice:     &ninjacrawler.SingleSelector{Selector: "div.p-product-detail > div.p-product-detail__price > div.p-product-detail__price__main span.c-text-text.c-text-text--en.c-text-text--bold"},
		Attributes: func(ctx ninjacrawler.CrawlerContext) []ninjacrawler.AttributeItem {
			return []ninjacrawler.AttributeItem{}
		},
	}
	crawler.Collection(constant.ProductDetails).SetConcurrentLimit(50).DisableRendering().IsDynamicPage(false).CrawlPageDetail(constant.Products)
}
