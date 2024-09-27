package panasonic_ec

import (
	"bytes"
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {
	crawler.Crawl([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        handleCategory,
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Categories,
			Processor:        handleProduct,
		},
	})

}
func handleCategory(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	items := []ninjacrawler.UrlCollection{}
	requestBody := map[string]interface{}{
		"RANK":     "",
		"CART":     "",
		"CATEGORY": "",
		"LEVEL":    2,
		"COIN":     "",
		"COUPON":   "",
		"MYPLAN":   "",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		ctx.App.Logger.Error("Error marshaling JSON:", "error", err.Error())
	}

	URL := ctx.UrlCollection.Url + "/store/ap/storeaezc/a2A/Common"

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		panic("Couldn't make category API request.")
	}
	req.Header.Set("User-Agent", ctx.App.Config.GetString("USER_AGENT", "PostmanRuntime/7.39.0"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.App.Logger.Error("Error sending request:", "error", err.Error())
	}
	defer resp.Body.Close()

	var response struct {
		Coupon   string `json:"COUPON"`
		MyPlan   string `json:"MYPLAN"`
		Category []struct {
			CategoryNameSP     string `json:"CATEGORY_NAME_SP"`
			CategoryCode       string `json:"CATEGORY_CODE"`
			IsLeaf             string `json:"IS_LEAF"`
			CategoryNamePC     string `json:"CATEGORY_NAME_PC"`
			ParentCategoryCode string `json:"PARENT_CATEGORY_CODE"`
			Level              string `json:"LEVEL"`
			FilePath           string `json:"FILE_PATH"`
			GzFPss             string `json:"GZ_F_PSS"`
			FileName           string `json:"FILE_NAME"`
		} `json:"CATEGORY"`
	}

	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		ctx.App.Logger.Error("Error parsing JSON:", "error", err.Error())
	}

	for _, item := range response.Category {
		if item.Level == "2" {
			categoryCode := item.CategoryCode
			categoryUrl := fmt.Sprintf("https://ec-plus.panasonic.jp/store/ap/storeaez/a2A/Product?CATEGORY_CODE=%s", categoryCode)

			items = append(items, ninjacrawler.UrlCollection{
				Url:    categoryUrl,
				Parent: ctx.UrlCollection.Url,
			})
		}
	}
	return items

}
func handleProduct(ctx ninjacrawler.CrawlerContext, next func([]ninjacrawler.UrlCollection, string)) error {
	productUrls := []ninjacrawler.UrlCollection{}
	totalPages := calculateTotalPage(ctx.Document)
	currentPage, _ := strconv.Atoi(ctx.Document.Find("li.pd_current a").First().Text())

	ctx.Document.Find("ul.pd_b-searchResultPanel_list > li").Each(func(i int, s *goquery.Selection) {
		aTag := s.Find("a").First()
		href, ok := aTag.Attr("href")
		if !ok {
			ctx.App.Logger.Warn("Product URL not found.")
		}

		productUrls = append(productUrls, ninjacrawler.UrlCollection{
			Url:    ctx.App.GetFullUrl(href),
			Parent: ctx.UrlCollection.Url,
		})
	})
	if currentPage == totalPages {
		next(productUrls, "")
		return nil
	} else {
		next(productUrls, generatePaginatedUrl(ctx.UrlCollection.Url, currentPage))
	}
	err := ctx.App.SaveHtml(ctx.Document, ctx.UrlCollection.Url)
	if err != nil {
		return err
	}
	return nil
}
func calculateTotalPage(document *goquery.Document) int {
	lastPage := document.Find("ol.pd_m-pagenation_list li a").Last().Text()
	totalPages, _ := strconv.Atoi(lastPage)
	return totalPages
}

func generatePaginatedUrl(urlStr string, pageNumber int) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	queryParams := parsedURL.Query()
	queryParams.Set("PAGE", strconv.Itoa(pageNumber+1))
	parsedURL.RawQuery = queryParams.Encode()
	nextPageUrl := parsedURL.String()
	return nextPageUrl
}
