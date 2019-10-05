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
	ID           bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string        `json:"name" bson:"name"`
	Title        string        `json:"title" binding:"required" bson:"title"` // 标题
	Url          string        `json:"url" bson:"url"`
	Desc         string        `json:"desc" bson:"desc"`
	Content      string        `json:"content" bson:"content"` // 内容
	ImgName      string        `json:"imgName" bson:"imgName"`
	Tags         []string      `json:"tags" bson:"tags"`     // 标签
	Rating       int           `json:"rating" bson:"rating"` // 评分
	OpenCount    int           `json:"openCount" bson:"openCount"`
	CommentCount int           `json:"commentCount" bson:"commentCount"`

	CreatedBy bson.ObjectId `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	CreatedAt time.Time     `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
