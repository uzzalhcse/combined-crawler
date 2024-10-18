package sony

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {
	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        handleCategory,
		},
		{
			Entity:           constant.SubCategories,
			OriginCollection: constant.Categories,
			Processor:        handleSubCategory,
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.SubCategories,
			Processor:        handleProducts,
			Engine: ninjacrawler.Engine{
				IsDynamic: ninjacrawler.Bool(true),
			},
		},
	})
}

func handleCategory(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	items := []ninjacrawler.UrlCollection{}
	ctx.Document.Find(".s5-categoryList__itemInner a").Each(func(_ int, s *goquery.Selection) {
		href, exist := s.Attr("href")
		if exist {
			fullUrl := ctx.App.GetFullUrl(href)
			if isValidHost(fullUrl) {
				items = append(items, ninjacrawler.UrlCollection{
					Url:    fullUrl,
					Parent: ctx.UrlCollection.Url,
				})
			}
		}
	})
	return items
}

func handleProducts(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	items := []ninjacrawler.UrlCollection{}
	ctx.Document.Find("a.GalleryListItem__Button,a.s5-ACTFOCUSListItemMin__button,a.s5-buttonV3").Each(func(_ int, s *goquery.Selection) {
		href, exist := s.Attr("href")
		if exist && !strings.HasPrefix(href, "#") {
			fullUrl := ctx.App.GetFullUrl(href)
			if isValidHost(fullUrl) {
				items = append(items, ninjacrawler.UrlCollection{
					Url:    fullUrl,
					Parent: ctx.UrlCollection.Url,
				})
			}
		}
	})
	return items
}

func handleSubCategory(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	items := []ninjacrawler.UrlCollection{}
	found := false
	ctx.Document.Find("ul.CategoryNav__MainList li a,.s5-categoryNav__item a").Each(func(_ int, s *goquery.Selection) {
		menu := s.Text()
		if strings.Contains(menu, "商品一覧") {
			href, exist := s.Attr("href")
			if exist {
				fullUrl := ctx.App.GetFullUrl(href)
				if isValidHost(fullUrl) {
					items = append(items, ninjacrawler.UrlCollection{
						Url:    fullUrl,
						Parent: ctx.UrlCollection.Url,
					})
				}
			}

			found = true
		}
	})
	if !found {
		promotionsAndOffers := ctx.Document.Find(".promotionsAndOffers")
		if promotionsAndOffers.Length() > 1 {
			// Remove the last item
			promotionsAndOffers.Last().Remove()
			// Now only the first item of .promotionsAndOffers will be processed
			ctx.Document.Find(".promotionsAndOffers a").Each(func(i int, s *goquery.Selection) {
				href, exist := s.Attr("data-link")
				if exist {
					fullUrl := ctx.App.GetFullUrl(href)
					if isValidHost(fullUrl) {
						items = append(items, ninjacrawler.UrlCollection{
							Url:    fullUrl,
							Parent: ctx.UrlCollection.Url,
						})
					}
				}
			})
		}
		ctx.Document.Find(".s5-list .s5-listItem4__image a").Each(func(_ int, s *goquery.Selection) {
			href, exist := s.Attr("href")
			if exist {
				fullUrl := ctx.App.GetFullUrl(href)
				if isValidHost(fullUrl) {
					items = append(items, ninjacrawler.UrlCollection{
						Url:    fullUrl,
						Parent: ctx.UrlCollection.Url,
					})
				}
			}
		})
	}
	//fmt.Println("invalidUrls", invalidUrls)
	return items
}
func isValidHost(urlString string) bool {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Url parsing error:", err)
		return false
	}

	hostname := parsedUrl.Hostname()
	if hostname == "sony.jp" || hostname == "www.sony.jp" {
		return true
	}

	return false
}
