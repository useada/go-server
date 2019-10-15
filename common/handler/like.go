package handler

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	myjwt "serve/middleware"
	"serve/models"

	"github.com/gin-gonic/gin"
)

func ListUserLike(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	claims := c.MustGet("claims").(*myjwt.CustomClaims)
	userID := claims.ID

	like, err := models.GetLikeByUserID(db, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	if like == nil || len(like.LinkList) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "success",
			"data": gin.H{
				"count": 0,
				"links": make([]models.ResultLink, 0),
			},
		})
		return
	}

	var links []models.Link
	query := bson.M{
		"_id": bson.M{
			"$in": like.LinkList,
		},
	}
	err = db.C(models.CollectionLink).Find(query).All(&links)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	resultLinks := make([]models.ResultLink, 0)
	for _, v := range links {
		resultLink := models.LinkToResultLink(&v)
		resultLink.Liked = true
		resultLinks = append(resultLinks, *resultLink)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "Success",
		"data": gin.H{
			"count": len(resultLinks),
			"links": resultLinks,
		},
	})
}

func UpdateUserLike(c *gin.Context) {
	var req struct {
		OP     string        `json:"op" binding:"required" bson:"op"`
		LinkID bson.ObjectId `json:"linkID" binding:"required" bson:"linkID"`
	}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	op := req.OP
	if op == "" {
		op = "add"
	}

	claims := c.MustGet("claims").(*myjwt.CustomClaims)
	userID := claims.ID

	db := c.MustGet("db").(*mgo.Database)

	_, err = models.UpdateLikeByUserID(db, op, userID, req.LinkID)
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
