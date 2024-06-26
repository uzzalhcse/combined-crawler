package models

import (
	"time"
)

type UrlCollection struct {
	CollectionID   string                 `json:"collection_id" bson:"collection_id"`
	Url            string                 `json:"url" bson:"url"`
	Parent         string                 `json:"parent" bson:"parent"`
	CurrentPageUrl string                 `json:"current_page_url" bson:"current_page_url"`
	Status         bool                   `json:"status" bson:"status"`
	Error          bool                   `json:"error" bson:"error"`
	Attempts       int                    `json:"attempts" bson:"attempts"`
	MetaData       map[string]interface{} `json:"meta_data" bson:"meta_data"`
	CreatedAt      time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt      *time.Time             `json:"updated_at" bson:"updated_at"`
}

func (c *UrlCollection) GetTableName() string {
	return "url_collections"
}
