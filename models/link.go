package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionLink = "links"
)

// Link model
type Link struct {
	ID        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string        `json:"name" binding:"required" bson:"name"`
	Url       string        `json:"url" binding:"required" bson:"url"`
	Desc      string        `json:"desc" bson:"desc"`
	ImgName   string        `json:"imgName" bson:"imgName"`
	Tags      []string      `json:"tags" bson:"tags"`
	OpenCount int           `json:"openCount" bson:"openCount"`
	CreatedBy bson.ObjectId `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	CreatedAt time.Time     `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
