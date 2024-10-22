package asics

import (
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

func CategoryHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var urls []ninjacrawler.UrlCollection
	categories := []string{"men", "women", "kids", "sports", "sportstyle", "sale-outlet"}
	for _, cat := range categories {
		selector := fmt.Sprintf("li[data-menu='%s'] ul.menu-vertical a", cat)

		ctx.Document.Find(selector).Each(func(i int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			if shouldSkipURL(url) {
				return
			}
			url = strings.ReplaceAll(url, "&sz=24", "")
			urls = append(urls, ninjacrawler.UrlCollection{
				Url:    url + "?start=0&sz=" + strconv.Itoa(ProductUrlPerPage),
				Parent: ctx.Page.URL(),
			})
		})
		//subCategories, _ := ctx.Page.Locator(selector).All()
		//for _, subCategory := range subCategories {
		//	url, _ := subCategory.GetAttribute("href")
		//	if shouldSkipURL(url) {
		//		continue
		//	}
		//	url = strings.ReplaceAll(url, "&sz=24", "")
		//	urls = append(urls, ninjacrawler.UrlCollection{
		//		Url:    url + "?start=0&sz=" + strconv.Itoa(ProductUrlPerPage),
		//		Parent: ctx.Page.URL(),
		//	})
		//}
	}
	return urls
}

func shouldSkipURL(href string) bool {
	skipPatterns := []string{
		"/mk/", "/shoe-finder/", "/search/", "/ja50000000/", "/ja10000000/", "gnavi-sports",
		"?srule=top-selling_2", "?srule=new-arrivals", "?prefn1=productArea",
	}
	for _, pattern := range skipPatterns {
		if strings.Contains(href, pattern) {
			return true
		}
	}
	return false
}

func ProductUrlHandler(ctx ninjacrawler.CrawlerContext, next func([]ninjacrawler.UrlCollection, string)) error {

	fmt.Println("Product Url Handler")
	productCountSelector := "#primary > div.search-result-options.search-result-multiple-filters > h1 > span"
	//ctx.Page.WaitForSelector(productCountSelector)
	productCount := ctx.Document.Find(productCountSelector).Text()
	totalProductCount, _ := strconv.Atoi(ctx.App.ToNumericsString(productCount))

	pages := totalProductCount / ProductUrlPerPage
	if totalProductCount%ProductUrlPerPage != 0 {
		pages += 1
	}

	productSelector := "[id=\"search-result-items\"] li.grid-tile .product-tile__link"

	hasProducts := ctx.Document.Find(productSelector).Length()
	for page := 1; page <= pages && hasProducts != 0; page++ {
		var urls []ninjacrawler.UrlCollection

		ctx.Document.Find(productSelector).Each(func(i int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			urls = append(urls, ninjacrawler.UrlCollection{
				Url:    url,
				Parent: Url,
			})
		})
		next(urls, "")
		break
		//ctx.App.Logger.Info("Next Page -> ", page)
		//nextPageUrl := ctx.Page.URL()
		//if strings.Contains(nextPageUrl, "?start=") {
		//	nextPageUrl = nextPageUrl[:strings.Index(nextPageUrl, "?start=")]
		//}
		//start := strconv.Itoa(page) + "00"
		//nextPageUrl = nextPageUrl + "?start==" + start + "&sz=" + strconv.Itoa(ProductUrlPerPage)
		//ctx.App.Logger.Info("Total Product Url's-> ", len(urls))
		//ctx.App.Logger.Info("Next Page URL -> ", nextPageUrl)
		//_, err := ctx.Page.Goto(nextPageUrl)
		//if err != nil {
		//	ctx.App.Logger.Error("Next page error.", err)
		//}
		////ctx.Page.WaitForSelector(productSelector)
		//hasProducts, _ = ctx.Page.Locator(productSelector).Count()
	}
	return nil
}
