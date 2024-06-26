package models

import "combined-crawler/pkg/ninjacrawler"

type CrawlerConfig struct {
	SiteID     string                         `json:"site_id" bson:"site_id"`
	Engine     ninjacrawler.Engine            `json:"engine"`
	Processors []ninjacrawler.ProcessorConfig `json:"processors"`
	Preference ninjacrawler.AppPreference     `json:"preferences"`
}

func (c *CrawlerConfig) GetTableName() string {
	return "crawler_configs"
}
