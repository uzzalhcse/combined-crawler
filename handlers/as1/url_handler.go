package as1

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {
	subCategorySelector := ninjacrawler.UrlSelector{
		Selector:     "#af-categories-list > li",
		SingleResult: false,
		FindSelector: "a",
		Attr:         "href",
	}
	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        categoryHandler,
		},
		{
			Entity:           constant.SubCategories,
			OriginCollection: constant.Categories,
			Processor:        subCategorySelector,
		},
		{
			Entity:           constant.Series,
			OriginCollection: constant.SubCategories,
			Processor:        seriesAndProductHandler,
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Series,
			Processor:        productHandler,
		},
	})
}
func categoryHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var urls []ninjacrawler.UrlCollection
	categoryDiv := ctx.Document.Find("#menuList > div.submenu > div.accordion > dl")

	// Get the total number of items in the selection
	totalCats := categoryDiv.Length()

	// Iterate over all items except the last three
	categoryDiv.Slice(0, totalCats-3).Each(func(i int, cat *goquery.Selection) {
		cat.Find("dd > div > ul > li").Each(func(j int, li *goquery.Selection) {
			li.Find("ul > li").Each(func(j int, l *goquery.Selection) {
				l.Find("ul > li").Each(func(j int, lMain *goquery.Selection) {
					href, ok := lMain.Find("a").Attr("href")
					if ok {
						urls = append(urls, ninjacrawler.UrlCollection{
							Url:    ctx.App.GetFullUrl(href),
							Parent: ctx.UrlCollection.Url,
						})
						ctx.App.Logger.Info("Category URL %s", href)
					} else {
						ctx.App.Logger.Error("Category URL not found")
					}
				})
			})
		})
	})
	return urls
}

func seriesAndProductHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var seriesUrls []ninjacrawler.UrlCollection
	var prdUrls []ninjacrawler.UrlCollection

	scrapeAnchors := func(s *goquery.Selection, urlType string) {
		href, ok := s.Find("a").Attr("href")
		fullUrl := ctx.App.GetFullUrl(href)
		if ok {
			if urlType == "series" {
				seriesUrls = append(seriesUrls, ninjacrawler.UrlCollection{
					Url:    fullUrl,
					Parent: ctx.UrlCollection.Url,
				})
			} else if urlType == "product" {
				prdUrls = append(prdUrls, ninjacrawler.UrlCollection{
					Url:    fullUrl,
					Parent: ctx.UrlCollection.Url,
				})
			}
		} else {
			ctx.App.Logger.Warn("URL not found")
		}
	}

	mainFunc := func(pageDt *goquery.Document) {
		pageDt.Find("#af-product-list > ul > li").Each(func(i int, li *goquery.Selection) {
			pElements := li.Find("div > div > p")
			if pElements.Length() > 2 {
				thirdP := pElements.Eq(2)
				productsLength := thirdP.Find("a").Length()
				if productsLength > 1 {
					scrapeAnchors(li, "series")
				} else {
					scrapeAnchors(thirdP, "product")
				}
			}
		})
	}

	totalItems := ctx.Document.Find("#af-result-count > span").Text()
	ctx.App.Logger.Info("Total items found: %s", totalItems)
	re := regexp.MustCompile(`[-+]?(?:\d*\.\d+|\d+)`)
	totalNo := re.FindString(totalItems)
	// Check if total number is numeric
	perBlockItem := 40
	// Convert total number to float64
	total, err := strconv.ParseFloat(totalNo, 64)
	if err != nil {
		ctx.App.Logger.Error("Failed to convert total number to float64 %s", err.Error())
	}
	pageTotal := int(math.Ceil(total / float64(perBlockItem)))

	for i := 1; i <= pageTotal; i++ {
		pageURL := fmt.Sprintf("%s&page=%d", ctx.UrlCollection.Url, i)
		ctx.App.Logger.Info("Fetching page %s", pageURL)
		doc, pwErr := ctx.App.NavigateToStaticURL(ctx.App.GetHttpClient(), pageURL, ctx.App.CurrentProxy)
		if pwErr != nil {
			ctx.App.Logger.Error("Failed to fetch page data page_url %s error %s", pageURL, pwErr.Error())
		}
		mainFunc(doc)
	}

	ctx.App.InsertUrlCollections(constant.Products, prdUrls, ctx.UrlCollection.Url)
	return seriesUrls
}
func productHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var productUrls []ninjacrawler.UrlCollection

	series := getSeries(ctx.Document)
	category := getCategory(ctx.Document)
	description := getDescription(ctx.Document)

	ctx.Document.Find("#af-product-list-body > tr").Each(func(i int, tr *goquery.Selection) {
		href, ok := tr.Find("td.af-group-item-productno > a").First().Attr("href")
		fullUrl := ctx.App.GetFullUrl(href)
		if ok {
			productUrls = append(productUrls, ninjacrawler.UrlCollection{
				Url:    fullUrl,
				Parent: ctx.UrlCollection.Url,
				MetaData: map[string]interface{}{
					"series":      series,
					"category":    category,
					"description": description,
				},
			})
		} else {
			ctx.App.Logger.Error("Product URL not found")
		}
	})
	return productUrls
}
func getSeries(doc *goquery.Document) string {
	series := doc.Find("#af-groupcd > div.groupdetail__name > h1").Text()
	return series
}

func getCategory(doc *goquery.Document) string {
	var categoryTexts []string
	lis := doc.Find("body > div.container > nav > ul").First().Find("li")
	liCount := lis.Length()
	if liCount < 3 {
		return ""
	}
	// Iterate through the list items, excluding the first, second, and last items
	lis.Each(func(i int, li *goquery.Selection) {
		if i > 1 && i < liCount-1 { // Exclude the first, second, and last items
			text := li.Text()
			categoryTexts = append(categoryTexts, strings.TrimSpace(text))
		}
	})
	category := strings.Join(categoryTexts, " > ")
	return category
}

func getDescription(doc *goquery.Document) string {
	uls := doc.Find("body > div.container > nav > ul")
	ulCount := uls.Length()
	if ulCount <= 1 {
		return ""
	} else {
		var descriptionTexts []string
		uls.Each(func(j int, ul *goquery.Selection) {
			if j > 0 {
				var categoryTexts []string
				lis := ul.Find("li")
				liCount := lis.Length()
				if liCount > 2 {
					// Iterate through the list items, excluding the first, second, and last items
					lis.Each(func(i int, li *goquery.Selection) {
						if i > 1 && i < liCount-1 { // Exclude the first, second, and last items
							text := li.Text()
							categoryTexts = append(categoryTexts, strings.TrimSpace(text))
						}
					})
				}
				category := strings.Join(categoryTexts, " > ")
				descriptionTexts = append(descriptionTexts, strings.TrimSpace(category))
			}
		})
		description := strings.Join(descriptionTexts, " | ")
		return description
	}
}
