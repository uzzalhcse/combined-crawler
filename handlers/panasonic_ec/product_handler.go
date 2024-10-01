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
				Adapter:         ninjacrawler.String(ninjacrawler.PlayWrightEngine),
				ConcurrentLimit: 6,
				IsDynamic:       ninjacrawler.Bool(false),
				//BlockResources:  true,
				WaitForSelector: ninjacrawler.String("div.pd_c-price"),
				//WaitForDynamicRendering: true,
				//ProxyStrategy: ninjacrawler.ProxyStrategyRotation,
				//ProxyServers: []ninjacrawler.Proxy{
				//	{
				//		Server:   "http://5.59.251.78:6117",
				//		Username: "lnvmpyru",
				//		Password: "5un1tb1azapa",
				//	},
				//	{
				//		Server:   "http://5.59.251.19:6058",
				//		Username: "lnvmpyru",
				//		Password: "5un1tb1azapa",
				//	},
				//	{
				//		Server:   "http://62.164.231.7:9319",
				//		Username: "lnvmpyru",
				//		Password: "5un1tb1azapa",
				//	},
				//	{
				//		Server:   "http://192.46.190.170:6763",
				//		Username: "lnvmpyru",
				//		Password: "5un1tb1azapa",
				//	},
				//	{
				//		Server:   "http://130.180.233.112:7683",
				//		Username: "lnvmpyru",
				//		Password: "5un1tb1azapa",
				//	},
				//},
			},
			Preference: ninjacrawler.Preference{ValidationRules: []string{"PageTitle", "ProductName"}},
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
