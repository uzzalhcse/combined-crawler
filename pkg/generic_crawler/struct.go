package generic_crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

type SiteConfig struct {
	Name        string       `json:"name"`
	Url         string       `json:"url"`
	Collections []Collection `json:"collections"`
}
type Collection struct {
	Entity           string     `json:"entity"`
	OriginCollection string     `json:"originCollection"`
	Preference       Preference `json:"preference"`
	Engine           Engine     `json:"engine"`
	Processor        Processor  `json:"processor"`
}
type Processor struct {
	UrlProcessor     UrlProcessor     `json:"url_processor"`
	ElementProcessor ElementProcessor `json:"element_processor"`
}
type UrlProcessor struct {
	Plugin   string   `json:"plugin"`
	Selector Selector `json:"selector"`
}
type ElementProcessor struct {
	Plugin   string        `json:"plugin"`
	Elements []ElementType `json:"elements"`
}
type ElementType struct {
	ElementID     string        `json:"element_id"`
	Plugin        string        `json:"plugin"`
	Selector      Selector      `json:"selector"`
	MultiSelector MultiSelector `json:"multi_selector"`
	Formatters    []Formatter   `json:"formatters"` //
}

type Selector struct {
	Query string `json:"query"` // CSS selector
	Attr  string `json:"attr"`  // Attribute to extract (e.g., "src" or "href")
}

type Formatter struct {
	Type string                 `toml:"type"`
	Args map[string]interface{} `toml:"args"`
}

type MultiSelector struct {
	Selectors     []Selector `json:"selectors"`      // Array of selectors
	ExcludeString []string   `json:"exclude_string"` //
	IsUnique      bool       `json:"is_unique"`      //
}

type CrawlerContext struct {
	App           *GenericCrawler
	Document      *goquery.Document
	UrlCollection UrlCollection
	Page          playwright.Page
	Scope         Map
}
type AppPreference struct {
	ExcludeUniqueUrlEntities []string
}
type Preference struct {
	DoNotMarkAsComplete bool
	ValidationRules     []string
}
type Engine struct {
}
type Map map[string]interface{}

func (m Map) Get(key string) interface{} {
	return m[key]
}
func (m Map) Set(key string, value interface{}) {
	m[key] = value
}
