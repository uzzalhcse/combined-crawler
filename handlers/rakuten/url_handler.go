package rakuten

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {
	categorySelector := ninjacrawler.UrlSelector{
		Selector:     "ul.header-globalNav__inner__nav li ul li.sub-genre-popup__inner__list__item",
		FindSelector: "a",
		Attr:         "href",
	}
	//productSelector := ninjacrawler.UrlSelector{
	//	Selector:     ".item-title,.item,.slickItem,.carousel__inner-slide",
	//	FindSelector: "a",
	//	Attr:         "href",
	//}
	crawler.Crawl([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        categorySelector,
			Engine:           ninjacrawler.Engine{},
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Categories,
			Processor:        handleProduct,
			Engine:           ninjacrawler.Engine{
				//IsDynamic: ninjacrawler.Bool(false),
			},
		},
	})

}

func handleProduct(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	items := []ninjacrawler.UrlCollection{}
	ctx.Document.Find(".item-title,dl.details,.slickItem,.carousel__inner-slide,.rankingContents,.rbcomp__item-list__item__details__lead h3").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a").Attr("href")
		href = ctx.App.GetFullUrl(href)
		items = append(items, ninjacrawler.UrlCollection{
			Url:    href,
			Parent: ctx.UrlCollection.Url,
		})
	})
	if len(items) == 0 {
		handleAuthorPublisherCategories(ctx)
	}
	return items
}

func handleAuthorPublisherCategories(ctx ninjacrawler.CrawlerContext) {
	// Find the author categories using the selector
	selection := ctx.Document.Find(".twoColumn td.etc_rank a")

	// Check if the selection contains any elements
	if selection.Length() == 0 {
		ctx.App.Logger.Warn("No author categories found %s", ctx.UrlCollection.Url)
		return
	}

	// Proceed if elements are found
	items := []ninjacrawler.UrlCollection{}
	selection.Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		href = ctx.App.GetFullUrl(href)
		items = append(items, ninjacrawler.UrlCollection{
			Url:    href,
			Parent: ctx.UrlCollection.Url,
		})
	})

	// Insert the URL collections and log the count if items exist
	if len(items) > 0 {
		ctx.App.InsertUrlCollections(constant.Categories, items, ctx.UrlCollection.Url)
		ctx.App.Logger.Info("Total Author Categories: %d", len(items))
	} else {
		ctx.App.Logger.Warn("Unknown Layout %s", ctx.UrlCollection.Url)
	}
}
