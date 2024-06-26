package models

type CrawlingHistory struct {
	SiteID      string        `json:"site_id" bson:"site_id"`
	Duration    string        `json:"duration" bson:"duration"`
	DataChanges []interface{} `json:"changes" bson:"changes"` // track the history of product detail property changes
}
type DataChanges struct {
	Name    string    `json:"name" bson:"name"`       // selling_price
	Changes []Changes `json:"changes" bson:"changes"` // [{},{}]
}
type Changes struct {
	Value string `json:"value" bson:"value"` // 185.06
	Date  string `json:"date" bson:"date"`   // 2024-05-25
}

func (c *SiteCollection) CrawlingHistory() string {
	return "sites"
}
