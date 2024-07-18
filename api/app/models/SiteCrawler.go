package models

type SiteCrawler struct {
	SiteID       string `json:"site_id" bson:"site_id"`
	InstanceName string `json:"instance_name" bson:"instance_name"`
	Zone         string `json:"zone"  bson:"zone"`
	Status       string `json:"status" bson:"status"` // running, deleted, or stopped
}

func (c *SiteCrawler) GetTableName() string {
	return "site_crawler"
}
