package model

type Category struct {
	Site       string `bson:"site,omitempty" json:"site"`
	Url        string `bson:"url,omitempty" json:"url"`
	Tier       int    `bson:"tier,omitempty" json:"tier"`
	CategoryId string `bson:"_id,omitempty" json:"_id"`
}
