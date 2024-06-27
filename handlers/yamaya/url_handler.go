package yamaya

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {
	categorySelector := ninjacrawler.UrlSelector{
		Selector:     ".col-md p",
		SingleResult: false,
		FindSelector: "a.btn.btn-lg.btn-block.btn-outline-primary",
		Attr:         "href",
	}
	SubCategorySelector := ninjacrawler.UrlSelector{
		Selector:     ".row .card-body",
		SingleResult: false,
		FindSelector: "a.mx-auto",
		Attr:         "href",
		Handler: func(urlCollection ninjacrawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {
			storeId := strings.ReplaceAll(fullUrl, "/stores", "")
			href := storeId + "catalog/index.php"
			return href, nil
		},
	}
	SubSubCategorySelector := ninjacrawler.UrlSelector{
		Selector:     "div.no-gutters",
		SingleResult: false,
		FindSelector: "a",
		Attr:         "href",
	}

	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        categorySelector,
		},
		{
			Entity:           constant.SubCategories,
			OriginCollection: constant.Categories,
			Processor:        SubCategorySelector,
		},
		{
			Entity:           constant.SubSubCategories,
			OriginCollection: constant.SubCategories,
			Processor:        SubSubCategorySelector,
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.SubSubCategories,
			Processor:        productListHandler,
		},
	})

}

func productListHandler(ctx ninjacrawler.CrawlerContext, next func([]ninjacrawler.UrlCollection, string)) error {
	productUrls := []ninjacrawler.UrlCollection{}
	ctx.Document.Find("div.card").Each(func(i int, s *goquery.Selection) {
		cardBody := s.Find("div.card-body").Last()
		aTag := cardBody.Find("a").First()
		href, ok := aTag.Attr("href")
		if ok {
			href = ctx.App.GetFullUrl(href)
			productUrls = append(productUrls, ninjacrawler.UrlCollection{
				Url:      href,
				MetaData: nil,
				Parent:   ctx.UrlCollection.Url,
			})
		} else {
			ctx.App.Logger.Warn("Product URL not found for %s", ctx.UrlCollection.Url)
		}
	})

	nextPageUrl := ""
	lastPageUrl := ""
	currentUrl := strings.Split(ctx.UrlCollection.Url, "?")[0]
	ctx.Document.Find("a.page-link").Each(func(i int, s *goquery.Selection) {
		txt := strings.Trim(s.Text(), " \n")
		if txt == ">" {
			href, _ := s.Attr("href")
			nextPageUrl = currentUrl + href
		}
		if txt == ">>" {
			href, _ := s.Attr("href")
			lastPageUrl = currentUrl + href
		}
	})
	if ctx.UrlCollection.CurrentPageUrl == lastPageUrl {
		next(productUrls, "")
		return nil
	} else {
		next(productUrls, nextPageUrl)
	}
	return nil

}
