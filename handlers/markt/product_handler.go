package markt

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"strings"
)

func ProductHandler(crawler *ninjacrawler.Crawler) {

	//productDetailApi := ninjacrawler.ProductDetailApi{
	//	Jan: "item_code",
	//	PageTitle: func(ctx ninjacrawler.CrawlerContext) string {
	//		title := fmt.Sprintf("%v-%v【マルクト】-%v 【マルクト】", ctx.ApiResponse.Get("name"), ctx.ApiResponse.Get("shop.name"), ctx.ApiResponse.Get("shop.name"))
	//		return title
	//	},
	//	Url: func(ctx ninjacrawler.CrawlerContext) string { return ctx.UrlCollection.Url },
	//	Images: func(ctx ninjacrawler.CrawlerContext) []string {
	//		images, ok := ctx.ApiResponse["product_images"].([]interface{})
	//		if !ok {
	//			fmt.Println("product_images is not an array")
	//			return nil
	//		}
	//
	//		items := make([]string, len(images))
	//		for index, image := range images {
	//			items[index] = ctx.App.GetFullUrl("/html/upload/save_image/" + image.(string))
	//		}
	//		return items
	//	},
	//	ProductCodes: func(ctx ninjacrawler.CrawlerContext) []string {
	//		return []string{}
	//	},
	//	Maker: "place_of_origin",
	//	Brand: "middleCategories.1.name",
	//	ProductName: func(ctx ninjacrawler.CrawlerContext) string {
	//		return fmt.Sprintf("%v|%v", ctx.ApiResponse.Get("name"), ctx.ApiResponse.Get("amount"))
	//	},
	//	Category:    "middleCategories.0.name",
	//	Description: "description_detail",
	//	Reviews: func(ctx ninjacrawler.CrawlerContext) []string {
	//		return []string{}
	//	},
	//	ItemTypes: func(ctx ninjacrawler.CrawlerContext) []string {
	//		return []string{}
	//	},
	//	ItemSizes: func(ctx ninjacrawler.CrawlerContext) []string {
	//		return []string{}
	//	},
	//	ItemWeights: func(ctx ninjacrawler.CrawlerContext) []string {
	//		return []string{}
	//	},
	//	SingleItemSize:   "",
	//	SingleItemWeight: "",
	//	NumOfItems:       "",
	//	ListPrice:        "",
	//	SellingPrice:     "price",
	//	Attributes: func(ctx ninjacrawler.CrawlerContext) []ninjacrawler.AttributeItem {
	//		return []ninjacrawler.AttributeItem{}
	//	},
	//}
	productDetailSelector := ninjacrawler.ProductDetailSelector{
		Jan: &ninjacrawler.SingleSelector{
			Selector: "div.p-product-detail > div.p-product-detail__code > dl:nth-child(1) > span",
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
	crawler.Collection(constant.ProductDetails).IsDynamicPage(false).CrawlPageDetail(constant.Products, productDetailSelector, "PageTitle")
}
