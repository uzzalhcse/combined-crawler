package panasonic

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func getImagesService(ctx ninjacrawler.CrawlerContext) []string {
	baseUrl := ctx.App.BaseUrl
	images := []string{}

	url := strings.ReplaceAll(ctx.UrlCollection.Url, baseUrl, "")
	urlParts := strings.Split(url, "/")
	if len(urlParts) < 2 {
		return images // Return empty slice if urlParts has less than 2 elements
	}
	urlStarting := urlParts[1]

	mainSection := ctx.Document.Find("section#maincontents").First()
	if mainSection == nil {
		return images
	}
	productDetailsSections := getProductDetailsSections(ctx.Document)
	sections := []*goquery.Selection{}
	firstFigure := mainSection.Find("figure").First()

	sections = append(sections, firstFigure)
	sections = append(sections, productDetailsSections)

	uniqueImages := make(map[string]bool)
	for _, section := range sections {
		if section == nil {
			continue
		}
		section.Find("img").Each(func(i int, s *goquery.Selection) {
			dataSrc, dataSrcOk := s.Attr("data-src")
			src, _ := s.Attr("src")
			imgUrl := ""

			if dataSrcOk {
				imgUrl = dataSrc
			} else {
				imgUrl = src
			}
			if !strings.HasPrefix(imgUrl, "/"+urlStarting) {
				return
			}
			url := baseUrl + imgUrl

			if !uniqueImages[url] {
				uniqueImages[url] = true
				images = append(images, url)
			}
		})
	}

	return images
}

func getProductCodesService(ctx ninjacrawler.CrawlerContext) []string {
	productCodes := []string{}
	if !strings.Contains(ctx.UrlCollection.Url, "products/") {
		return productCodes
	}

	productUrlParts := strings.Split(ctx.UrlCollection.Url, "/")
	productCode := productUrlParts[len(productUrlParts)-1]
	productCode = strings.Split(productCode, ".")[0]

	productCodes = append(productCodes, productCode)

	return productCodes
}
func ProductNameHandler(ctx ninjacrawler.CrawlerContext) string {

	productMainSection := getProductMainSection(ctx.Document)
	productName := ""
	if productMainSection == nil {
		return productName
	}

	productNameSection := productMainSection.Find("h1").First()
	productName = strings.Trim(productNameSection.Text(), " \t\n")

	return productName
}

func getCategoryService(ctx ninjacrawler.CrawlerContext) string {
	category := ""
	productDetailsSections := getProductDetailsSections(ctx.Document)
	additionalPage, _ := getAdditionalPage(ctx, productDetailsSections)
	if additionalPage == nil {
		return category
	}

	categoryItems := []string{}
	categorySection := additionalPage.Find("ol").First()
	categorySection.Find("li").Each(func(i int, s *goquery.Selection) {
		txt := strings.Trim(s.Text(), " \n")
		categoryItems = append(categoryItems, txt)
	})
	category = strings.Join(categoryItems, " > ")

	return category
}
func getAdditionalPage(ctx ninjacrawler.CrawlerContext, productDetailsSection *goquery.Selection) (*goquery.Document, error) {
	url := getRelevantUrl(productDetailsSection, ctx)
	if url == "" || !ctx.App.IsValidPage(url) {
		return nil, nil
	}

	const maxAttempts = 2
	const retryDelay = 2 * time.Second

	var document *goquery.Document
	var err error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		document, err = ctx.App.NavigateToStaticURL(ctx.App.GetHttpClient(), url, ctx.App.CurrentProxy)
		if err == nil {
			ctx.App.Logger.Warn("Attempt %d: Successful navigation: %v", attempt)
			return document, nil // Successful navigation, return the document
		}

		ctx.App.Logger.Warn("Attempt %d: Error navigating to page: %v", attempt, err)

		if attempt == maxAttempts {
			_ = ctx.App.MarkAsError(ctx.UrlCollection.Url, constant.Products, err.Error(), 1)
			ctx.App.Logger.Error("Error navigating to page after %d attempts: %v", maxAttempts, err)
		} else {
			//ctx.App.HandleThrottling(attempt, 494)
			time.Sleep(retryDelay)
		}
	}

	return nil, err
}
func getRelevantUrl(productDetailsSection *goquery.Selection, ctx ninjacrawler.CrawlerContext) string {
	url := ""
	if productDetailsSection == nil {
		return url
	}

	urls := productDetailsSection.Find("a")
	urls.Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			urlParts := strings.Split(href, "/")
			for _, elem := range urlParts {
				if elem == "p-db" && strings.Contains(href, "_spec") {
					url = ctx.App.GetFullUrl(href)
				}
			}
		}
	})

	return url
}

func getDescriptionService(ctx ninjacrawler.CrawlerContext) string {
	description := ""

	if strings.Contains(ctx.UrlCollection.Url, "/products/") {
		descriptionSection := ctx.Document.Find("p.header2").First()
		if descriptionSection != nil {
			description = ctx.App.HtmlToText(descriptionSection)
		}
	} else {
		descriptionSection := ctx.Document.Find("meta[name='description']")
		if descriptionSection != nil {
			desc, ok := descriptionSection.Attr("content")
			if ok {
				description = desc
			}
		}
	}

	return description
}

func GetProductAttribute(ctx ninjacrawler.CrawlerContext) []ninjacrawler.AttributeItem {
	attributes := []ninjacrawler.AttributeItem{}

	productDetailsSections := getProductDetailsSections(ctx.Document)
	additionalPage, _ := getAdditionalPage(ctx, productDetailsSections)

	getProductDetailsAttributeService(ctx, productDetailsSections, &attributes)
	getRelevantUrlAttributeService(ctx, productDetailsSections, &attributes)
	getSpecialKeyTableAttributeService(additionalPage, &attributes)
	getSpecialKeyAttributeService(additionalPage, &attributes)
	getLinkedUrlsAttributeService(ctx, ctx.Document, &attributes)
	getMentionedProductCodes(ctx.Document, &attributes)

	return attributes
}
func getProductDetailsAttributeService(ctx ninjacrawler.CrawlerContext, productDetailsSection *goquery.Selection, attributes *[]ninjacrawler.AttributeItem) {
	productDetails := ""
	if productDetailsSection == nil {
		return
	}

	productDetailsSectionCopy := productDetailsSection
	productDetailsSectionCopy.Find("a").Remove()
	productDetailsSectionCopy.Find("img").Remove()

	productDetails = ctx.App.HtmlToText(productDetailsSectionCopy)
	productDetails = strings.Trim(productDetails, " \t\u00a0\n")
	productDetails = strings.ReplaceAll(productDetails, "\t", "\n")
	productDetails = RemoveConsecutiveCharacter(productDetails, " ")
	productDetails = RemoveConsecutiveCharacter(productDetails, "\n")
	productDetails = RemoveConsecutiveCharacter(productDetails, " \n")
	productDetails = RemoveConsecutiveCharacter(productDetails, "\t\n")
	productDetails = RemoveConsecutiveCharacter(productDetails, "\n\t")

	if len(productDetails) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   "product_details",
			Value: productDetails,
		}
		*attributes = append(*attributes, attribute)
	}
}
func getRelevantUrlAttributeService(ctx ninjacrawler.CrawlerContext, productDetailsSection *goquery.Selection, attributes *[]ninjacrawler.AttributeItem) {
	url := getRelevantUrl(productDetailsSection, ctx)

	if len(url) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   "relevant_url",
			Value: url,
		}
		*attributes = append(*attributes, attribute)
	}
}
func getSpecialKeyTableAttributeService(additionalPage *goquery.Document, attributes *[]ninjacrawler.AttributeItem) {
	tableData := ""
	if additionalPage == nil {
		return
	}

	mainContent := additionalPage.Find("section#maincontents").First()
	table := mainContent.Find("table").First()
	if len(table.Text()) > 0 {
		keyText := table.Parents().Find("div.bgWhite").First().Text()
		keyText = strings.Trim(keyText, " \t\n")

		table.Find("tr").Each(func(i int, s *goquery.Selection) {
			tableText := ""
			s.Children().Each(func(j int, h *goquery.Selection) {
				txt := strings.Trim(h.Text(), " \t\n")
				if len(tableText) > 0 && len(txt) > 0 {
					tableText += "/"
				}
				tableText += txt
			})

			if len(tableData) > 0 && len(tableText) > 0 {
				tableData += "\n"
			}
			tableData += tableText
		})

		if len(keyText) > 0 && len(tableData) > 0 {
			attribute := ninjacrawler.AttributeItem{
				Key:   keyText,
				Value: tableData,
			}
			*attributes = append(*attributes, attribute)
		}
	}
}

func getSpecialKeyAttributeService(additionalPage *goquery.Document, attributes *[]ninjacrawler.AttributeItem) {
	if additionalPage == nil {
		return
	}

	additionalPage.Find("div.pagesection").Each(func(i int, s *goquery.Selection) {
		h3Text := s.Find("h3").First().Text()
		h3Text = strings.Trim(h3Text, " \n")
		if h3Text == "付属品" {
			nextSectionText := s.Next().First().Text()
			nextSectionText = strings.Trim(nextSectionText, " \t\n")

			if len(nextSectionText) > 0 {
				attribute := ninjacrawler.AttributeItem{
					Key:   h3Text,
					Value: nextSectionText,
				}
				*attributes = append(*attributes, attribute)
			}
		}
	})
}

func getLinkedUrlsAttributeService(ctx ninjacrawler.CrawlerContext, pageData *goquery.Document, attributes *[]ninjacrawler.AttributeItem) {
	linkedUrls := []ninjacrawler.UrlCollection{}
	linkedUrlsStr := ""
	pageData.Find("a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok || strings.Contains(href, "#") || strings.Contains(href, "products") || strings.Contains(href, "javascript") || href == "" {
			return
		}

		link := ""
		href = strings.ReplaceAll(href, "\n", "")

		if strings.HasPrefix(href, "//") {
			href = "https:" + href
		}
		if strings.HasPrefix(href, "http") {
			link = href
		} else {
			link = JoinURLs(ctx.App.CurrentUrl, href)
		}
		if isValidHost(link) {
			linkedUrls = append(linkedUrls, ninjacrawler.UrlCollection{
				Url:    link,
				Parent: ctx.UrlCollection.Url,
			})
		}
	})
	linkedUrls, linkedUrlsStr = MakeUrlsUnique(linkedUrls)

	if len(linkedUrls) > 0 {
		ctx.App.InsertUrlCollections(constant.Products, linkedUrls, ctx.UrlCollection.Url)
		attribute := ninjacrawler.AttributeItem{
			Key:   "linked_urls",
			Value: linkedUrlsStr,
		}
		*attributes = append(*attributes, attribute)
	}
}

func getMentionedProductCodes(pageData *goquery.Document, attributes *[]ninjacrawler.AttributeItem) {
	fullPageText := pageData.Text()

	re, err := regexp.Compile(`\b[A-Z0-9]{2,}-?[A-Z0-9\/]{1,}-?[A-Z0-9]{1}\b`)
	if err != nil {
		fmt.Println("regex error.")
		return
	}

	validProductCode := []string{}
	match := re.FindAllString(fullPageText, -1)
	for _, item := range match {
		if strings.Contains(item, "-") && !strings.HasPrefix(item, "GTM-") && containsUpperCase(item) {
			validProductCode = append(validProductCode, item)
		}
	}
	validProductCode = MakeStringSliceUnique(validProductCode)
	mentionedCodes := strings.Join(validProductCode, ", ")

	if len(mentionedCodes) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   "mentioned_product_codes",
			Value: mentionedCodes,
		}
		*attributes = append(*attributes, attribute)
	}
}

func getProductDetailsSections(pageData *goquery.Document) *goquery.Selection {
	productDetailsSectionHeaders := pageData.Find("div.bgLightGray")

	var productDetailsSections *goquery.Selection = nil
	productDetailsSectionHeaders.Each(func(i int, s *goquery.Selection) {
		h2Text := s.Find("h2").First().Text()
		text := strings.Trim(h2Text, " \n")
		if text == "新商品のおすすめポイント" || text == "ご購入をお考えのお客様へ" {
			return
		}

		selection := s
		for {
			selection = selection.Next()
			if selection.Text() == "" || selection.HasClass("bgLightGray") {
				break
			}
			if productDetailsSections == nil {
				productDetailsSections = selection
			} else {
				productDetailsSections = productDetailsSections.AddSelection(selection)
			}
		}
	})

	return productDetailsSections
}
func getProductMainSection(pageData *goquery.Document) *goquery.Selection {
	mainSection := pageData.Find("section#maincontents").First()

	var productMainSection *goquery.Selection
	mainSection.Find("div.pagesection").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("bgLightGray") {
			return
		}
		if productMainSection == nil {
			productMainSection = s
		} else {
			productMainSection = productMainSection.AddSelection(s)
		}
	})

	return productMainSection
}
func RemoveConsecutiveCharacter(fullString string, subString string) string {
	for strings.Contains(fullString, subString+subString) {
		fullString = strings.ReplaceAll(fullString, subString+subString, subString)
	}

	return fullString
}
func JoinURLs(url1 string, url2 string) string {
	if url2[0] == '/' {
		return GetBaseUrl(url1) + url2
	}

	url1Parsed, err := url.Parse(url1)
	if err != nil {
		return ""
	}
	url2Parsed, err := url.Parse(url2)
	if err != nil {
		return ""
	}
	joinedURL := url1Parsed.ResolveReference(url2Parsed)

	return joinedURL.String()
}
func GetBaseUrl(urlString string) string {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("failed to parse URL:", err)
		return ""
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host
	return baseURL
}
func MakeUrlsUnique(itemList []ninjacrawler.UrlCollection) ([]ninjacrawler.UrlCollection, string) {
	// Convert the list to a set (remove duplicates)
	uniqueSet := make(map[string]bool)
	var uniqueList []ninjacrawler.UrlCollection
	var linkedUrlsStr []string
	for _, item := range itemList {
		if !uniqueSet[item.Url] {
			uniqueSet[item.Url] = true
			uniqueList = append(uniqueList, item)
			linkedUrlsStr = append(linkedUrlsStr, item.Url)
		}
	}
	return uniqueList, strings.Join(linkedUrlsStr, ", ")
}
func MakeStringSliceUnique(itemList []string) []string {
	// Convert the list to a set (remove duplicates)
	uniqueSet := make(map[string]bool)
	var uniqueList []string
	for _, item := range itemList {
		if !uniqueSet[item] {
			uniqueSet[item] = true
			uniqueList = append(uniqueList, item)
		}
	}
	return uniqueList
}
func containsUpperCase(s string) bool {
	for _, char := range s {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}
