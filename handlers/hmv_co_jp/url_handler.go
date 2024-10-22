package hmv_co_jp

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
)

var totalCount = 0

func UrlHandler(crawler *ninjacrawler.Crawler) {
	//categorySelector := ninjacrawler.UrlSelector{
	//	Selector:     "ul.listSubInnerList li.listSmallSub",
	//	FindSelector: "a",
	//	Attr:         "href",
	//}
	//productSelector := ninjacrawler.UrlSelector{
	//	Selector:     "ul.resultList li .thumbnailBlock",
	//	FindSelector: "a",
	//	Attr:         "href",
	//}
	crawler.Crawl([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        categoryHandler,
			Engine:           ninjacrawler.Engine{
				//IsDynamic: ninjacrawler.Bool(false),
			},
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Categories,
			Processor:        productHandler,
			Engine:           ninjacrawler.Engine{
				//IsDynamic: ninjacrawler.Bool(false),
			},
		},
	})

}
func categoryHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	items := []ninjacrawler.UrlCollection{}
	ctx.Document.Find("ul.listInner li.listSmall").Each(func(_ int, s *goquery.Selection) {
		s.Find("a").Each(func(_ int, s *goquery.Selection) {
			url, exist := s.Attr("href")
			if exist {
				items = append(items, ninjacrawler.UrlCollection{
					Url:    ctx.App.GetFullUrl(url),
					Parent: ctx.UrlCollection.Url,
				})
			}
		})
	})
	return items
}
func productHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	items := []ninjacrawler.UrlCollection{}
	doc := ctx.Document.Find(".view")
	if doc.Length() == 0 {
		ctx.App.Logger.Warn("searchContents not found %s", ctx.UrlCollection.Url)
		return items
	}
	totalText := ctx.Document.Find(".view").Text()

	// Use regex to extract the total number before '件中'
	re := regexp.MustCompile(`(\d+)件中`)
	matches := re.FindStringSubmatch(totalText)
	total := 0

	if len(matches) > 1 {
		total, _ = strconv.Atoi(matches[1]) // First capturing group is the total number
		totalCount += total
		ctx.App.Logger.Info("Total results: %d", totalCount)
	} else {
		ctx.App.Logger.Error("Could not extract total results")
	}
	ctx.Document.Find("ul.resultList li .thumbnailBlock").Each(func(_ int, s *goquery.Selection) {
		url, exist := s.Find("a").Attr("href")
		if exist {
			items = append(items, ninjacrawler.UrlCollection{
				Url:    ctx.App.GetFullUrl(url),
				Parent: ctx.UrlCollection.Url,
				MetaData: map[string]interface{}{
					"total": total,
				},
			})
		}
	})
	return items
}
