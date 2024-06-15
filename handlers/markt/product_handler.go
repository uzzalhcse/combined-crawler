package markt

import (
	"combined-crawler/constant"
	"github.com/PuerkitoBio/goquery"
	"github.com/lazuli-inc/ninjacrawler"
	"strings"
)

func ProductHandler(crawler *ninjacrawler.Crawler) {

	crawler.ProductDetailSelector = ninjacrawler.ProductDetailSelector{
		Jan: "",
		PageTitle: &ninjacrawler.SingleSelector{
			Selector: "title",
		},
		Url: getUrlHandler,
		Images: &ninjacrawler.MultiSelectors{
			Selectors: []ninjacrawler.Selector{
				{Query: "img#image-item", Attr: "src"},
				{Query: "section.ProductDetail_Section_Function a", Attr: "href"},
				{Query: "section.ProductDetail_Section_Spec img", Attr: "src"},
			},
		},
		ProductCodes: productCodeHandler,
		Maker:        "",
		Brand:        "",
		ProductName:  productNameHandler,
		Category:     "",
		Description:  "",
	}
	crawler.Collection(constant.ProductDetails).CrawlPageDetail(constant.Products)
}
func productCodeHandler(app ninjacrawler.Crawler, document *goquery.Document, urlCollection ninjacrawler.UrlCollection) []string {
	urlParts := strings.Split(strings.Trim(urlCollection.Url, "/"), "/")
	return []string{urlParts[len(urlParts)-1]}
}

func productNameHandler(app ninjacrawler.Crawler, document *goquery.Document, urlCollection ninjacrawler.UrlCollection) string {
	return strings.Trim(document.Find("h2.ProductInfo_Head_Main_ProductName").Text(), " \n")
}

func getUrlHandler(app ninjacrawler.Crawler, document *goquery.Document, urlCollection ninjacrawler.UrlCollection) string {
	return urlCollection.Url
}
