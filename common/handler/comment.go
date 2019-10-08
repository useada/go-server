package handler

import (
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"serve/models"

	"github.com/gin-gonic/gin"
)

// List Comments
func GetComments(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	referToStr := c.DefaultQuery("referTo", "")
	indexStr := c.Query("index")
	countStr := c.Query("count")

	if referToStr == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 499,
			"msg":    "referTo empty",
		})
		return
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		index = 0
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		count = 10
	}

	query := bson.M{
		"referTo": bson.ObjectIdHex(referToStr),
	}

	sort := "-createdAt"

	commentsCount, err := db.C(models.CollectionComment).Find(query).Count()
	if err != nil {
		log.WithFields(log.Fields{
			"referTo": referToStr,
			"index":   index,
			"count":   count,
			"err":     err,
		}).Error("open db error")

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	var comments []models.Comment
	err = db.C(models.CollectionComment).Find(query).Sort(sort).Skip(index).Limit(count).All(&comments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "Success",
		"data": gin.H{
			"count":    commentsCount,
			"comments": comments,
		},
	})
}

// Create a Comment
func CreateComment(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	var comment models.Comment
	err := c.BindJSON(&comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	var link models.Link
	err = db.C(models.CollectionLink).
		FindId(comment.ReferTo).
		One(&link)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	comment.Type = 1
	comment.CreatedAt = time.Now()

	comment.ID = bson.NewObjectId()

	err = db.C(models.CollectionComment).Insert(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	var recentComments []models.Comment
	recentComments = append(recentComments, comment)

	for k, v := range link.RecentComments {
		if k < 2 {
			recentComments = append(recentComments, v)
		}
	}
	link.RecentComments = recentComments

	err = db.C(models.CollectionLink).Update(bson.M{"_id": comment.ReferTo},
		bson.M{"$inc": bson.M{"commentCount": 1}, "$set": bson.M{"recentComments": recentComments}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "Success",
		"data":   comment,
	})
}

// Delete comment
func DeleteComment(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	var comment models.Comment

	err := db.C(models.CollectionComment).
		FindId(bson.ObjectIdHex(c.Param("_id"))).
		One(&comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	err = db.C(models.CollectionComment).Update(bson.M{"_id": bson.ObjectIdHex(c.Param("_id"))},
		bson.M{"$set": bson.M{"delFlag": 1}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "Success",
	})
}
