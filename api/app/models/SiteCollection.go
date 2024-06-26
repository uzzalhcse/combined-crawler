package models

type SiteCollection struct {
	SiteID        string        `json:"site_id" bson:"site_id"`
	Name          string        `json:"name" bson:"name"`
	Url           string        `json:"url" bson:"url"`
	BaseUrl       string        `json:"base_url" bson:"base_url"`
	Status        string        `json:"status" bson:"status"`
	Collections   []Collection  `json:"collections" bson:"collections"`
	CrawlerConfig CrawlerConfig `json:"crawler_config" bson:"crawler_config"`
}

func (c *SiteCollection) GetTableName() string {
	return "sites"
}
