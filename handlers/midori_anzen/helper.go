package midori_anzen

import (
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func productCategoryHandler(ctx ninjacrawler.CrawlerContext) string {
	categoryItems := make([]string, 0)
	ctx.Document.Find(".breadcrumb li").Each(func(i int, s *goquery.Selection) {
		totalItems := ctx.Document.Find(".breadcrumb li").Length()
		// Skip the first and last items
		if i > 0 && i < totalItems-1 {
			txt := strings.TrimSpace(s.Text())
			categoryItems = append(categoryItems, txt)
		}
	})
	return strings.Join(categoryItems, " > ")
}
func productNameHandler(ctx ninjacrawler.CrawlerContext) string {
	productName := ctx.Document.Find("h2.example").Text()
	productName = strings.Trim(productName, " \n")

	return productName
}

func getUrlHandler(ctx ninjacrawler.CrawlerContext) string {
	return ctx.UrlCollection.Url
}

func getProductDescription(ctx ninjacrawler.CrawlerContext) string {
	description := ""
	ctx.Document.Find(".goods_a_wrapper div").Each(func(i int, s *goquery.Selection) {
		description += fmt.Sprintf("%s\n", strings.TrimSpace(s.Text()))
	})
	return description
}

func getProductAttribute(ctx ninjacrawler.CrawlerContext) []ninjacrawler.AttributeItem {
	attributes := []ninjacrawler.AttributeItem{}
	attributes = append(attributes, ninjacrawler.AttributeItem{Key: "selling_price_tax", Value: "1"}) // Put it in to determine that it is tax included

	extractAttributes(ctx, &attributes)

	getFunctionAttributes(ctx, &attributes)
	parseMeasurementTable(ctx, &attributes)
	getTagAttributes(ctx, &attributes)
	getCoordinationAttributes(ctx, &attributes)
	getHelmetBodiesAttributes(ctx, &attributes)
	return attributes
}

func extractAttributes(ctx ninjacrawler.CrawlerContext, attributes *[]ninjacrawler.AttributeItem) {
	salesUnitAttr := findTextByThContent(ctx.Document, "販売単位")
	if salesUnitAttr != "" {
		*attributes = append(*attributes, ninjacrawler.AttributeItem{Key: "販売単位", Value: salesUnitAttr})
	}

	quantityAttr := findTextByThContent(ctx.Document, "入数")
	if quantityAttr != "" {
		*attributes = append(*attributes, ninjacrawler.AttributeItem{Key: "入数", Value: quantityAttr})
	}

	shippingDateAttr := findTextByThContent(ctx.Document, "出荷予定日")
	if shippingDateAttr != "" {
		*attributes = append(*attributes, ninjacrawler.AttributeItem{Key: "出荷予定日", Value: shippingDateAttr})
	}
	// Extracting maker
	policyAttr := findTextByThContent(ctx.Document, "返品可否", "img", "alt")
	if policyAttr != "" {
		*attributes = append(*attributes, ninjacrawler.AttributeItem{Key: "返品可否", Value: policyAttr})
	}

	ctx.Document.Find(".goods-extendeditem .table-bordered tr").Each(func(i int, s *goquery.Selection) {
		key := strings.TrimSpace(s.Find("th").Text())
		value := strings.TrimSpace(s.Find("td").Text())
		if key != "" && value != "" {
			if key == "使用区分" {
				// Extract the document file URL from the <a> tag within <td>
				docURL, exists := s.Find("td p a").Attr("onclick")
				if exists {
					// Find the start and end indices of the URL within the onclick attribute
					urlStart := strings.Index(docURL, "'https://") + 1
					urlEnd := strings.Index(docURL[urlStart:], "'") + urlStart

					// Extract the document URL
					documentURL := docURL[urlStart:urlEnd]
					value = documentURL

					// Here you can handle the document URL as needed, such as storing it in a struct or variable
				}
			}
			attribute := ninjacrawler.AttributeItem{
				Key:   key,
				Value: value,
			}
			*attributes = append(*attributes, attribute)
		}
	})
}
func getFunctionAttributes(ctx ninjacrawler.CrawlerContext, attributes *[]ninjacrawler.AttributeItem) {
	key := strings.TrimSpace(ctx.Document.Find(".row_contents.goods_contents .item_title").Last().Text())
	value := ""
	ctx.Document.Find(".row_contents.goods_contents img:not(.img-responsive)").Each(func(i int, s *goquery.Selection) {
		val, ok := s.Attr("src")
		if key != "" && val != "" && ok {
			val = ctx.App.GetFullUrl(val)
		}
		value += fmt.Sprintf("%s \n", val)
	})
	if key != "" && value != "" {
		attribute := ninjacrawler.AttributeItem{
			Key:   key,
			Value: value,
		}
		*attributes = append(*attributes, attribute)

	}
}
func getTagAttributes(ctx ninjacrawler.CrawlerContext, attributes *[]ninjacrawler.AttributeItem) {
	key := "tags"
	tags := ctx.Document.Find("div.event-season.event-season_ss span").First().Text()
	if tags != "" {
		attribute := ninjacrawler.AttributeItem{
			Key:   key,
			Value: tags,
		}
		*attributes = append(*attributes, attribute)
	}
}
func getCoordinationAttributes(ctx ninjacrawler.CrawlerContext, attributes *[]ninjacrawler.AttributeItem) {
	key := ctx.Document.Find(".top_title").Last().Text()
	if key == "コーディネート" {

		value := ""
		ctx.Document.Find(".event-contents .event-contents").Last().Each(func(i int, s *goquery.Selection) {
			s.Find(".event-goods .event-price-img a").Each(func(i int, s *goquery.Selection) {
				url, _ := s.Attr("href")
				title, _ := s.Attr("title")
				url = ctx.App.GetFullUrl(url)
				value += fmt.Sprintf("%s\n/%s\n", title, url)
			})
		})

		attribute := ninjacrawler.AttributeItem{
			Key:   key,
			Value: value,
		}
		*attributes = append(*attributes, attribute)
	}
}
func getHelmetBodiesAttributes(ctx ninjacrawler.CrawlerContext, attributes *[]ninjacrawler.AttributeItem) {
	count := ctx.Document.Find(".goods_c .entry_p .table-responsive").Length()
	if count == 0 {
		return
	}
	key := ""
	ctx.Document.Find(".entry_p").Last().Contents().Each(func(i int, node *goquery.Selection) {
		if goquery.NodeName(node) == "#text" {
			key += node.Text()
		}
	})

	var value = ""
	ctx.Document.Find(".table-responsive table tbody tr").Each(func(i int, s *goquery.Selection) {
		// Extract value from each column
		columns := s.Find("td")
		// Skip the header row
		if i == 0 {
			value += fmt.Sprintf("%s / %s\n",
				columns.Eq(0).Text(),
				columns.Eq(1).Text(),
			)
		}

		if columns.Length() == 3 {
			value += fmt.Sprintf("%s %s / %s\n",
				columns.Eq(0).Text(),
				columns.Eq(1).Find("img").AttrOr("src", ""),
				columns.Eq(2).Text(),
			)
		}
	})

	// Print the extracted value

	if key != "" && value != "" {
		attribute := ninjacrawler.AttributeItem{
			Key:   key,
			Value: value,
		}
		*attributes = append(*attributes, attribute)
	}
}
func parseMeasurementTable(ctx ninjacrawler.CrawlerContext, attributes *[]ninjacrawler.AttributeItem) {
	key := "measurement_information"
	// Initialize slices for sizes and measurement types
	sizes := []string{}
	measurementTypes := []string{}
	measurements := make(map[string][]string)

	// Extract headers (sizes)
	headerLabel := ""
	ctx.Document.Find(".table-responsive:last-of-type table:first-of-type tr:first-child").Each(func(index int, rowHtml *goquery.Selection) {
		rowHtml.Find("td").Each(func(cellIdx int, cellHtml *goquery.Selection) {
			if cellIdx == 0 {
				headerLabel = strings.TrimSpace(cellHtml.Text())
			} else { // Skip the first header cell (headerLabel)
				sizes = append(sizes, strings.TrimSpace(cellHtml.Text()))
			}
		})
	})
	if len(sizes) < 2 {
		return
	}
	// Extract measurement types and their corresponding data
	ctx.Document.Find(".table-responsive:last-of-type table:first-of-type tr").Each(func(rowIdx int, rowHtml *goquery.Selection) {
		if rowIdx > 0 { // Skip the first row (headers)
			var measurementType string
			rowHtml.Find("td").Each(func(cellIdx int, cellHtml *goquery.Selection) {
				if cellIdx == 0 { // First column contains the measurement category
					measurementType = strings.TrimSpace(cellHtml.Text())
					if _, ok := measurements[measurementType]; !ok {
						measurements[measurementType] = make([]string, len(sizes))
						measurementTypes = append(measurementTypes, measurementType)
					}
				} else {
					measurements[measurementType][cellIdx-1] = strings.TrimSpace(cellHtml.Text())
				}
			})
		}
	})
	unitInfo := ""
	ctx.Document.Find(".table-responsive:last-of-type table:last-of-type tr td").Each(func(index int, cellHtml *goquery.Selection) {
		unitInfo = strings.TrimSpace(cellHtml.Text())
	})
	measurement_information := ""
	// Print the header row (measurement types)
	headerOutput := headerLabel
	for _, measurementType := range measurementTypes {
		headerOutput += fmt.Sprintf("/%s", measurementType)
	}

	measurement_information = headerOutput + "\n"
	// Print the extracted data in the "y:x" format
	for i, size := range sizes {
		output := size
		for _, measurementType := range measurementTypes {
			output += fmt.Sprintf("/%s", measurements[measurementType][i])
		}
		measurement_information += output + "\n"
	}
	measurement_information += unitInfo
	attribute := ninjacrawler.AttributeItem{
		Key:   key,
		Value: measurement_information,
	}
	*attributes = append(*attributes, attribute)
}

func getProductCode(ctx ninjacrawler.CrawlerContext) []string {
	return []string{findTextByThContent(ctx.Document, "商品コード")}
}
func getMaker(ctx ninjacrawler.CrawlerContext) string {
	return findTextByThContent(ctx.Document, "メーカー")
}
func findTextByThContent(doc *goquery.Document, thContent string, attr ...string) string {
	th := doc.Find("th").FilterFunction(func(i int, sel *goquery.Selection) bool {
		return strings.TrimSpace(sel.Text()) == thContent
	})

	// If found, get the text content of the corresponding <td>
	if th.Length() > 0 {
		td := th.Parent().Find("td")
		if td.Length() > 0 {
			if attr != nil && len(attr) > 0 {
				attrValue, ok := td.Find(attr[0]).Attr(attr[1])
				if ok {
					return attrValue
				}
			}
			text := strings.TrimSpace(td.Text())
			// Replace multiple whitespace characters with a single space
			text = strings.Join(strings.Fields(text), " ")
			return text

		}
	}

	return ""
}
func getItemSizes(ctx ninjacrawler.CrawlerContext) []string {
	sizes := extractAllThValues(ctx.Document)
	return sizes
}
func extractAllThValues(doc *goquery.Document) []string {
	var thValues []string

	// Find all <th> elements within the <thead>
	doc.Find(".goods_sizechart.table-responsive table.table-bordered thead th:not(:first-child)").Each(func(i int, s *goquery.Selection) {
		// Get the trimmed text content of each <th>
		thText := strings.TrimSpace(s.Text())
		thValues = append(thValues, thText)
	})

	return thValues
}
