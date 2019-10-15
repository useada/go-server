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
	ID             bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Type           int           `json:"type" bson:"type"` // 实名 0， 匿名 1
	Name           string        `json:"name" bson:"name"`
	Title          string        `json:"title" binding:"required" bson:"title"` // 标题
	Url            string        `json:"url" bson:"url"`
	Desc           string        `json:"desc" bson:"desc"`
	Content        string        `json:"content" bson:"content"` // 内容
	ImgName        string        `json:"imgName" bson:"imgName"`
	Tags           []string      `json:"tags" bson:"tags"`     // 标签
	Rating         int           `json:"rating" bson:"rating"` // 评分
	Author         Author        `json:"author" bson:"author"` // 作者
	OpenCount      int           `json:"openCount" bson:"openCount"`
	CommentCount   int           `json:"commentCount" bson:"commentCount"`
	RecentComments []Comment     `json:"recentComments" bson:"recentComments"` // 最近3条评论

	CreatedBy bson.ObjectId `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	CreatedAt time.Time     `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

type ResultLink struct {
	Link
	ImgUrl string `json:"imgUrl"`
	Liked  bool   `json:"liked"`
}

func LinkToResultLink(link *Link) *ResultLink {

	resultLink := ResultLink{}
	resultLink.Link = *link

	if resultLink.Title == "" {
		resultLink.Title = link.Name
	}
	if resultLink.Content == "" {
		resultLink.Content = link.Desc
	}

	if resultLink.Author.NickName == "" {
		resultLink.Author.NickName = "李白"
	}

	resultLink.ImgUrl = calcImgUrl(link.ImgName)
	return &resultLink
}

func calcImgUrl(imgName string) string {
	if imgName == "" {
		return ""
	}
	return "http://static.d36.net/funnylink/links/image/png/" + imgName + "-small"
}
