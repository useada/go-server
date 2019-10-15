package models

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionLike = "likes"
)

type Like struct {
	ID        bson.ObjectId   `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    bson.ObjectId   `json:"userID,omitempty" bson:"userID,omitempty"`
	LinkList  []bson.ObjectId `json:"linkList,omitempty" bson:"linkList,omitempty"`
	CreatedAt time.Time       `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

func GetLikeByUserID(db *mgo.Database, userID bson.ObjectId) (*Like, error) {
	var like Like

	err := db.C(CollectionLike).
		Find(bson.M{"userID": userID}).
		One(&like)
	if err != nil {
		if !strings.Contains(err.Error(), `not found`) {
			return nil, err
		} else {
			return nil, nil
		}
	}

	return &like, nil
}

func UpdateLikeByUserID(db *mgo.Database, op string, userID bson.ObjectId, linkID bson.ObjectId) (*Like, error) {
	var like Like

	err := db.C(CollectionLike).
		Find(bson.M{"userID": userID}).
		One(&like)
	if err != nil {
		if !strings.Contains(err.Error(), `not found`) {
			return nil, err
		} else {
			like.ID = bson.NewObjectId()
			like.UserID = userID
			like.CreatedAt = time.Now()
		}
	}

	if op == "add" && len(like.LinkList) >= 18 {
		return nil, fmt.Errorf("已满18个，无法继续添加")
	}
	if op == "del" && len(like.LinkList) <= 0 {
		return nil, nil
	}

	if op == "add" {
		for _, v := range like.LinkList {
			if v == linkID {
				return &like, nil
			}
		}

		like.LinkList = append(like.LinkList, linkID)
	}

	if op == "del" {
		var tmp []bson.ObjectId
		for _, v := range like.LinkList {
			if v != linkID {
				tmp = append(tmp, v)
			}
		}

		like.LinkList = tmp
	}

	// 查找原来的文档
	//query := bson.M{
	//	"_id": like.ID,
	//}

	// 更新
	//err = db.C(CollectionLike).Update(query, like)
	_, err = db.C(CollectionLike).UpsertId(like.ID, &like)

	return &like, err
}
