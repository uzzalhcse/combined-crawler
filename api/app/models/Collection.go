package models

type Collection struct {
	CollectionID   string          `json:"collection_id" bson:"collection_id"`
	SiteID         string          `json:"site_id" bson:"site_id"`
	Name           string          `json:"name" bson:"name"`
	UrlCollections []UrlCollection `json:"url_collections" bson:"url_collections"`
}

func (c *Collection) GetTableName() string {
	return "collections"
}
