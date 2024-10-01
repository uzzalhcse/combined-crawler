package panasonic_ec

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func ProductHandler(crawler *ninjacrawler.Crawler) {
	productDetailSelector := ninjacrawler.ProductDetailSelector{
		Jan: getJanService,
		PageTitle: &ninjacrawler.SingleSelector{
			Selector: "title",
		},
		Url:              GetUrlHandler,
		Images:           getImagesService,
		ProductCodes:     getProductCodesService,
		Maker:            "Panasonic",
		Brand:            getBrandService,
		ProductName:      getProductNameService,
		Category:         GetProductCategory,
		Description:      GetProductDescription,
		Reviews:          []string{},
		ItemTypes:        []string{},
		ItemSizes:        []string{},
		ItemWeights:      []string{},
		SingleItemSize:   "",
		SingleItemWeight: "",
		NumOfItems:       "",
		ListPrice:        "",
		SellingPrice:     getSellingPriceService,
		Attributes:       GetProductAttribute,
	}
	crawler.Crawl([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.ProductDetails,
			OriginCollection: constant.Products,
			Processor:        productDetailSelector,
			StateHandler:     handleState,
			Engine: ninjacrawler.Engine{
				ConcurrentLimit: 10,
				Adapter:         ninjacrawler.String(ninjacrawler.RodEngine),
				IsDynamic:       ninjacrawler.Bool(true),
				BlockResources:  true,
				//BlockedURLs: []string{
				//	"https://sprocket-ping.s3.amazonaws.com/",
				//},
				WaitForSelector: ninjacrawler.String("script[type='application/ld+json']"),
				//WaitForDynamicRendering: true,
			},
			Preference: ninjacrawler.Preference{ValidationRules: []string{"PageTitle"}},
		},
	})
}

func handleState(ctx ninjacrawler.CrawlerContext) ninjacrawler.Map {
	condition := "script[type='application/ld+json']"
	jsonData, specificationData := getSpecificationData(ctx, condition)
	data := ninjacrawler.Map{}
	data["jsonData"] = jsonData
	data["specificationData"] = specificationData
	return data
}

func getSpecificationData(ctx ninjacrawler.CrawlerContext, loadingCondition string) (ninjacrawler.Map, map[string]string) {
	scriptText := ""
	scriptElement := ctx.Document.Find(loadingCondition).First()
	if scriptElement.Length() > 0 {
		scriptText = scriptElement.Text()
	}

	var jsonData ninjacrawler.Map
	err := json.Unmarshal([]byte(scriptText), &jsonData)
	if err != nil {
		ctx.App.Logger.Error(err.Error())
	}

	specificationData := make(map[string]string)
	specTableDiv := ctx.Document.Find("div.specarea").First()
	if specTableDiv != nil {
		ths := []string{}
		tds := []string{}

		specTable := specTableDiv.Find("table").First()
		specTable.Find("th").Each(func(i int, s *goquery.Selection) {
			text := ctx.App.HtmlToText(s)
			ths = append(ths, text)
		})
		specTable.Find("td").Each(func(i int, s *goquery.Selection) {
			text := ctx.App.HtmlToText(s)
			text = strings.ReplaceAll(text, "\t", "")
			tds = append(tds, text)
		})

		if len(ths) != len(tds) {
			return jsonData, specificationData
		}

		for i := 0; i < len(ths); i++ {
			specificationData[ths[i]] = tds[i]
		}
	}

	return jsonData, specificationData
}
