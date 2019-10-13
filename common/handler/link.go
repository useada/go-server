package handler

import (
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	myjwt "serve/middleware"
	"serve/models"

	"github.com/gin-gonic/gin"
)

// List links
func GetLinks(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	tag := c.DefaultQuery("tag", "")

	// home, sns, admin
	op := c.DefaultQuery("op", "home")

	indexStr := c.Query("index")
	countStr := c.Query("count")

	if tag == "" {
		tag = "recommend"
	}

	if op == "sns" {
		tag = ""
	}

	if op == "admin" {
		tag = ""
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		index = 0
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		count = 10
	}

	log.WithFields(log.Fields{
		"tag": tag,
		"op":  op,
	}).Debug("query")

	query := bson.M{}
	if tag != "" {
		query = bson.M{
			"tags": tag,
		}
	}

	sort := "-rating"

	linksCount, err := db.C(models.CollectionLink).Find(query).Count()
	if err != nil {
		log.WithFields(log.Fields{
			"tag":   tag,
			"index": index,
			"count": count,
			"err":   err,
		}).Error("open db error")

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	var links []models.Link
	err = db.C(models.CollectionLink).Find(query).Sort(sort).Skip(index).Limit(count).All(&links)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	resultLinks := make([]ResultLink, 0)

	for _, v := range links {
		resultLink := ResultLink{}
		resultLink.Link = v

		if resultLink.Title == "" {
			resultLink.Title = v.Name
		}
		if resultLink.Content == "" {
			resultLink.Content = v.Desc
		}

		if resultLink.Author.NickName == "" {
			resultLink.Author.NickName = "李白"
		}

		resultLink.ImgUrl = calcImgUrl(v.ImgName)

		resultLinks = append(resultLinks, resultLink)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "Success",
		"data": gin.H{
			"count": linksCount,
			"links": resultLinks,
		},
	})
}

// List all links
func ListLinks(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	var links []models.Link
	err := db.C(models.CollectionLink).Find(nil).All(&links)
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
		"data":   links,
	})
}

// Get a link
func GetLink(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	var link models.Link

	err := db.C(models.CollectionLink).
		FindId(bson.ObjectIdHex(c.Param("_id"))).
		One(&link)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	resultLink := ResultLink{}
	resultLink.Link = link

	if resultLink.Title == "" {
		resultLink.Title = link.Name
	}
	if resultLink.Content == "" {
		resultLink.Content = link.Desc
	}

	resultLink.ImgUrl = calcImgUrl(link.ImgName)

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "Success",
		"data":   resultLink,
	})
}

// Create a link
func CreateLink(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	var link models.Link
	err := c.BindJSON(&link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}
	link.CreatedAt = time.Now()

	claims := c.MustGet("claims").(*myjwt.CustomClaims)
	userID := claims.ID
	user, err := models.GetUserByID(db, userID)
	if err != nil || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	link.Author.ID = userID
	link.Author.NickName = user.NickName
	link.ID = bson.NewObjectId()

	err = db.C(models.CollectionLink).Insert(link)
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
		"data":   link,
	})
}

// Create a links
func CreateLinks(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	var links []models.Link
	err := c.BindJSON(&links)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	var pLinks []interface{}

	for i, _ := range links {
		links[i].CreatedAt = time.Now()
		pLinks = append(pLinks, &links[i])
	}

	err = db.C(models.CollectionLink).Insert(pLinks...)
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

// Delete link
func DeleteLink(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	var link models.Link

	err := db.C(models.CollectionLink).
		FindId(bson.ObjectIdHex(c.Param("_id"))).
		One(&link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	//err = db.C(models.CollectionLink).Remove(bson.M{"_id": bson.ObjectIdHex(c.Param("_id"))})
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"status": 500,
	//		"msg":    err.Error(),
	//	})
	//	return
	//}

	err = db.C(models.CollectionLink).Update(bson.M{"_id": bson.ObjectIdHex(c.Param("_id"))},
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

// Update link
func UpdateLink(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	var link models.Link
	err := c.BindJSON(&link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	// 查找原来的文档
	query := bson.M{
		"_id": bson.ObjectIdHex(c.Param("_id")),
	}

	// 更新
	err = db.C(models.CollectionLink).Update(query, link)
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
		"data":   link,
	})
}

// Open link
func OpenLink(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	// 查找原来的文档
	query := bson.M{
		"_id": bson.ObjectIdHex(c.Param("_id")),
	}

	err := db.C(models.CollectionLink).Update(query,
		bson.M{"$inc": bson.M{"openCount": 1}})
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

type ResultLink struct {
	models.Link
	ImgUrl string `json:"imgUrl"`
}

func calcImgUrl(imgName string) string {
	if imgName == "" {
		return ""
	}
	return "http://static.d36.net/funnylink/links/image/png/" + imgName + "-small"
}
