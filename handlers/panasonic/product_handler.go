package panasonic

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
)

func ProductHandler(crawler *ninjacrawler.Crawler) {
	productDetailSelector := ninjacrawler.ProductDetailSelector{
		Jan: "",
		PageTitle: &ninjacrawler.SingleSelector{
			Selector: "title",
		},
		Url: func(ctx ninjacrawler.CrawlerContext) string {
			return ctx.UrlCollection.Url
		},
		Images:       getImagesService,
		ProductCodes: getProductCodesService,
		Maker:        "",
		Brand:        "",
		ProductName:  ProductNameHandler,
		Category:     getCategoryService,
		Description:  getDescriptionService,
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
		Attributes:       GetProductAttribute,
	}
	crawler.CrawlPageDetail([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.ProductDetails,
			OriginCollection: constant.Products,
			Processor:        productDetailSelector,
			Preference: ninjacrawler.Preference{
				ValidationRules: []string{"PageTitle|required|blacklist:ページが見つかりません | Panasonic,エラー,URL変更のお知らせ,Redirect"},
				PreHandlers: []func(c ninjacrawler.PreHandlerContext) error{
					ValidHost,
					HandleUrlExtension,
				},
			},
		},
	})
}

func ValidHost(c ninjacrawler.PreHandlerContext) error {
	if !isValidHost(c.UrlCollection.Url) {
		_ = c.App.MarkAsMaxErrorAttempt(c.UrlCollection.Url, constant.Products, "Invalid Host")
		return fmt.Errorf("invalid host %s", c.UrlCollection.Url)
	}
	return nil
}
func HandleUrlExtension(c ninjacrawler.PreHandlerContext) error {
	ext := GetUrlFileExtension(c.UrlCollection.Url)
	if ext != "" && ext != ".html" && ext != ".htm" {
		_ = c.App.MarkAsMaxErrorAttempt(c.UrlCollection.Url, constant.Products, "Invalid Url Extension")
		return fmt.Errorf("invalid Url Extension %s", c.UrlCollection.Url)
	}
	return nil
}
func GetUrlFileExtension(urlString string) string {
	u, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return ""
	}
	fileName := path.Base(u.Path)
	extension := filepath.Ext(fileName)

	return extension
}
