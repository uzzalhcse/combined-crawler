package models

type SiteSecret struct {
	SiteID  string                 `json:"site_id" bson:"site_id"`
	Secrets map[string]interface{} `json:"secrets" bson:"secrets"`
}

func (c *SiteSecret) GetTableName() string {
	return "site_secrets"
}
