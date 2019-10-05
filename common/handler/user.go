package handler

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	myjwt "serve/middleware"
	"serve/models"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

// List all users
func ListUser(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	var users []models.User
	err := db.C(models.CollectionUser).Find(nil).All(&users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Success",
		"data":   users,
	})
}

// Get a user
func GetUser(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	var user models.User

	err := db.C(models.CollectionUser).
		FindId(bson.ObjectIdHex(c.Param("_id"))).
		One(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Success",
		"data":   user,
	})
}

// Get a user from email
func GetEmail(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	var user models.User

	err := db.C(models.CollectionUser).
		Find(bson.M{"email": c.Param("email")}).
		One(&user)
	if err != nil {
		if !strings.Contains(err.Error(), `not found`) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": 500,
				"msg":    err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": 200,
				"msg":    "not found",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Success",
		"data":   user,
	})
}

// Create a user
func CreateUser(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	user.Role = 0
	user.CreatedAt = time.Now()

	err = db.C(models.CollectionUser).Insert(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	if user.OrganizationID != "" {
		err = db.C(models.CollectionOrg).Update(bson.M{"_id": user.OrganizationID},
			bson.M{"$inc": bson.M{"memberCount": 1}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": 500,
				"msg":    err.Error(),
			})
			return
		}
	}
	if user.CustomerID != "" {
		err = db.C(models.CollectionCustomer).Update(bson.M{"_id": user.CustomerID},
			bson.M{"$inc": bson.M{"memberCount": 1}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": 500,
				"msg":    err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Success",
	})
}

// Delete user
func DeleteUser(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	var user models.User

	err := db.C(models.CollectionUser).
		FindId(bson.ObjectIdHex(c.Param("_id"))).
		One(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	err = db.C(models.CollectionUser).Remove(bson.M{"_id": bson.ObjectIdHex(c.Param("_id"))})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	if user.OrganizationID != "" {
		err = db.C(models.CollectionOrg).Update(bson.M{"_id": user.OrganizationID},
			bson.M{"$inc": bson.M{"memberCount": -1}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": 500,
				"msg":    err.Error(),
			})
			return
		}
	}
	if user.CustomerID != "" {
		err = db.C(models.CollectionCustomer).Update(bson.M{"_id": user.CustomerID},
			bson.M{"$inc": bson.M{"memberCount": -1}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": 500,
				"msg":    err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Success",
	})
}

// Update user
func UpdateUser(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	user.Role = 0

	// 查找原来的文档
	query := bson.M{
		"_id": bson.ObjectIdHex(c.Param("_id")),
	}

	// 更新
	err = db.C(models.CollectionUser).Update(query, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Success",
		"data":   user,
	})
}

// List all organization users
func ListOrgUsers(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	var users []models.User
	query := bson.M{
		"organizationId": bson.ObjectIdHex(c.Param("_id")),
		"type":           1,
	}
	err := db.C(models.CollectionUser).Find(query).All(&users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Success",
		"data":   users,
	})
}

// List all customer users
func ListCustomerUsers(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	var users []models.User
	query := bson.M{
		"customerId": bson.ObjectIdHex(c.Param("_id")),
		"type":       2,
	}
	err := db.C(models.CollectionUser).Find(query).All(&users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Success",
		"data":   users,
	})
}

// LoginResult 登录结果结构
type LoginResult struct {
	Token string `json:"token"`
	models.User
}

// Login 登录
func Login(c *gin.Context) {
	var loginReq models.LoginReq
	if c.BindJSON(&loginReq) == nil {

		db := c.MustGet("db").(*mgo.Database)

		var user models.User
		query := bson.M{
			"email": loginReq.Email,
		}

		err := db.C(models.CollectionUser).
			Find(query).
			One(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": 500,
				"msg":    err.Error(),
			})
			return
		}

		if user.Password == "" || user.Password != loginReq.Password {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "Wrong Password!",
			})
		} else {
			generateToken(c, user)
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "Json parse error!",
		})
	}
}

// 生成令牌
func generateToken(c *gin.Context, user models.User) {
	j := &myjwt.JWT{
		SigningKey: []byte("FogDong"),
	}
	claims := myjwt.CustomClaims{
		ID: user.ID,
		//user.Email,
		StandardClaims: jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 过期时间 一小时
			Issuer:    "FogDong",                       //签名的发行者
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}

	data := LoginResult{
		User:  user,
		Token: token,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "登录成功！",
		"data":   data,
	})
	return
}

func AnonymousUser(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	index := int(rand.Int31n(30000))
	count := 1

	query := bson.M{}

	var idiom models.Idiom
	err := db.C(models.CollectionIdiom).Find(query).Skip(index).Limit(count).One(&idiom)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	type AnonymousUser struct {
		Nickname string `json:"nickname" bson:"nickname"`
	}
	au := AnonymousUser{Nickname: idiom.Word}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "Success",
		"data":   au,
	})
}
