package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionComment = "comment"
)

type Author struct {
	ID       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	NickName string        `json:"nickName" bson:"nickName"` // 昵称
}

type AtData struct {
	ID     bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Author Author        `json:"author" bson:"author"`
}

// Link model
type Comment struct {
	ID      bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Type    int           `json:"type" bson:"type"`       // 实名 0， 匿名 1
	Content string        `json:"content" bson:"content"` // 内容
	ReferTo bson.ObjectId `json:"referTo,omitempty" bson:"referTo,omitempty"`
	Author  Author        `json:"author" bson:"author"`                     // 作者
	AtData  AtData        `json:"atData,omitempty" bson:"atData,omitempty"` // at数据
	//CreatedBy bson.ObjectId `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
