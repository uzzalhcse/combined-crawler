package generic_crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"plugin"
	"reflect"
	"regexp"
	"strings"
)

const (
	baseCollection = "sites"
)

type GenericCrawler struct {
	*mongo.Client
	Name       string
	Url        string
	BaseUrl    string
	httpClient *http.Client
}

func NewGenericCrawler() *GenericCrawler {
	return &GenericCrawler{}
}

var pkg *plugin.Plugin
var pluginMap []string

func (gc *GenericCrawler) RunAutoPilot() {
	sites, err := loadSites("sites.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var siteConfig []SiteConfig

	for _, site := range sites {
		gc.Name = site.Name
		gc.Url = site.Url
		gc.BaseUrl = getBaseUrl(site.Url)
		gc.Client = gc.mustGetClient()
		gc.httpClient = gc.GetHttpClient()
		gc.newSite()
		for _, collection := range site.Collections {
			urlCollections := gc.getUrlCollections(collection.OriginCollection)
			for _, urlCollection := range urlCollections {
				doc, err := gc.NavigateToStaticURL(urlCollection.Url)
				if err != nil {
					fmt.Println("NavigateToStaticURL failed", err)
				}
				crawlerCtx := CrawlerContext{
					App:           gc,
					Document:      doc,
					UrlCollection: urlCollection,
					Scope:         Map{},
				}
				crawlerCtx.Scope.Set("crawler", "generic crawler")
				handleCollection(collection, crawlerCtx)
			}
		}

		siteConfig = append(siteConfig, site)
	}
}
func handleCollection(collection Collection, crawlerCtx CrawlerContext) {
	productDetail := &ProductDetail{}
	// Url Selector
	if collection.Processor.UrlProcessor.Plugin != "" {
		fnSymbol, err := handlePlugin("urls", collection.Processor.UrlProcessor.Plugin, crawlerCtx.App.Name)
		if err != nil {
			fmt.Println("handlePlugin error:", err)
			return
		}
		switch v := fnSymbol.(type) {
		case func(CrawlerContext, func([]UrlCollection, string)) error:
			fn, ok := fnSymbol.(func(CrawlerContext, func([]UrlCollection, string)) error)
			if !ok {
				fmt.Printf("Function %s has unexpected type in package\n", collection.Processor.UrlProcessor.Plugin)
			}
			handleErr := fn(crawlerCtx, func(collections []UrlCollection, currentPageUrl string) {
				fmt.Println("Insert to DB", len(collections))
			})
			if handleErr != nil {
				fmt.Println("Error inserting")
			}

		default:
			fmt.Printf("%s Invalid Signature: %T\n", collection.Processor.UrlProcessor.Plugin, v)
		}

	}
	// Element selector

	if collection.Processor.ElementProcessor.Plugin != "" {
		//fn, err := handlePlugin("elements", collection.Processor.ElementProcessor.Plugin, crawlerCtx.App.Name, crawlerCtx)
		//if err != nil {
		//	fmt.Println("handlePlugin error:", err)
		//	return
		//}
		//collection.Processor.ElementProcessor.Processor = fn
	} else {
		for _, element := range collection.Processor.ElementProcessor.Elements {
			var result interface{}
			if element.Plugin != "" {
				fnSymbol, err := handlePlugin("elements", element.Plugin, crawlerCtx.App.Name)
				if err != nil {
					fmt.Println("handlePlugin error:", err)
					return
				}
				switch v := fnSymbol.(type) {
				case func(CrawlerContext) interface{}:
					fn, ok := fnSymbol.(func(CrawlerContext) interface{})
					if !ok {
						fmt.Printf("Function %s has unexpected type in package\n", element.Plugin)
					}
					result = fn(crawlerCtx)

				default:
					fmt.Printf("%s Invalid Signature: %T\n", element.Plugin, v)
				}
			} else {
				result = processSelector(crawlerCtx, element)
			}

			field := reflect.ValueOf(productDetail).Elem().FieldByName(element.ElementID)
			fmt.Printf("%s %v\n", element.ElementID, result)

			switch result.(type) {
			case string:
				field.SetString(result.(string))
			case []string:
				field.Set(reflect.ValueOf(result))
			case []AttributeItem:
				field.Set(reflect.ValueOf(result))
			}
		}
	}
	//fmt.Println("productDetail", productDetail)
}
func handlePlugin(path, functionName, siteName string) (interface{}, error) {
	ElementPluginPath := fmt.Sprintf("plugins/%s/%s", siteName, path)
	pluginPath, err := buildPlugin(ElementPluginPath, functionName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !contains(pluginMap, ElementPluginPath) {
		pkg, err = plugin.Open(pluginPath)
		if err != nil {
			fmt.Println("err", err)
			os.Exit(1)
		}
		pluginMap = append(pluginMap, ElementPluginPath)
	}

	var fnSymbol interface{}
	// Look for the function by name in the package
	fnSymbol, err = pkg.Lookup(functionName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return fnSymbol, nil
}

func processSelector(crawlerCtx CrawlerContext, element ElementType) interface{} {

	var data interface{}
	isMultiSelector := !reflect.DeepEqual(element.MultiSelector, MultiSelector{})
	if isMultiSelector {
		//
		data = handleMultiSelectors(crawlerCtx, element.MultiSelector)
	} else {
		//single
		data = handleSingleSelector(crawlerCtx, element.Selector)
	}
	formattedData := formatProcess(data, element.Formatters)
	return formattedData
}
func formatProcess(data interface{}, formatters []Formatter) interface{} {
	for _, postProcess := range formatters {
		switch postProcess.Type {
		case "trim":
			data = trimProcess(data, postProcess.Args)
		case "replace":
			data = replaceProcess(data, postProcess.Args)
		default:
			log.Printf("Unknown post process type: %s", postProcess.Type)
		}
	}
	return data
}
func trimProcess(value interface{}, args map[string]interface{}) interface{} {
	switch v := value.(type) {
	case string:
		if args != nil {
			return strings.Trim(v, args["cutset"].(string))
		}
		return strings.TrimSpace(v)
	case []string:
		for i, str := range v {
			if args != nil {
				v[i] = strings.Trim(str, args["cutset"].(string))
			} else {
				v[i] = strings.TrimSpace(str)
			}
		}
		return v
	default:
		return value
	}
}

func replaceProcess(value interface{}, args map[string]interface{}) interface{} {
	switch v := value.(type) {
	case string:
		pattern := regexp.MustCompile(args["regex"].(string))
		return pattern.ReplaceAllString(v, args["repl"].(string))
	case []string:
		for i, str := range v {
			pattern := regexp.MustCompile(args["regex"].(string))
			v[i] = pattern.ReplaceAllString(str, args["repl"].(string))
		}
		return v
	default:
		return value
	}
}

func handleSingleSelector(crawlerCtx CrawlerContext, selector Selector) interface{} {
	txt := crawlerCtx.Document.Find(selector.Query).Text()
	return txt
}
func handleMultiSelectors(crawlerCtx CrawlerContext, selectors MultiSelector) interface{} {
	items := []string{}
	itemSet := make(map[string]struct{})

	// Helper function to append images if the specified attribute exists
	appendImages := func(selection *goquery.Selection, attr string) {
		selection.Each(func(i int, s *goquery.Selection) {
			if url, ok := s.Attr(attr); ok {
				fullUrl := crawlerCtx.App.GetFullUrl(url)

				// Check if the Url contains any excluded strings
				excluded := false
				for _, exclude := range selectors.ExcludeString {
					if strings.Contains(fullUrl, exclude) {
						excluded = true
						break
					}
				}
				if excluded {
					return
				}

				// Add to items if unique or uniqueness is not enforced
				if selectors.IsUnique {
					if _, exists := itemSet[fullUrl]; !exists {
						itemSet[fullUrl] = struct{}{}
						items = append(items, fullUrl)
					}
				} else {
					items = append(items, fullUrl)
				}
			}
		})
	}

	// Process each selector in the array
	for _, selector := range selectors.Selectors {
		appendImages(crawlerCtx.Document.Find(selector.Query), selector.Attr)
	}

	return items
}
