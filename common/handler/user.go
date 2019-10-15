package handler

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"net/http"
	"serve/util"
	"time"

	"github.com/mojocn/base64Captcha"
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
func GetThisUser(c *gin.Context) {
	claims := c.MustGet("claims").(*myjwt.CustomClaims)
	userID := claims.ID

	db := c.MustGet("db").(*mgo.Database)

	var user models.User
	err := db.C(models.CollectionUser).FindId(userID).One(&user)
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
		"data":   user,
	})
}

// Create a user
func CreateUser(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	var userData struct {
		Email    string `json:"email" binding:"required" bson:"email"`
		Password string `json:"password" binding:"required" bson:"password"`
		Phone    string `json:"phone,omitempty" bson:"phone,omitempty"`
		NickName string `json:"nickName,omitempty" binding:"required" bson:"nickName,omitempty"`
		Profile  string `json:"profile,omitempty" bson:"profile,omitempty"`
		models.CaptchaData
	}
	err := c.BindJSON(&userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	verifyResult := base64Captcha.VerifyCaptcha(userData.CaptchaID, userData.VerifyValue)
	if !verifyResult {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "验证码错误",
		})
		return
	}

	oldUser, err := models.GetUserByEmail(db, userData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg":    err.Error(),
		})
		return
	}

	if oldUser != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 500,
			"msg":    "邮箱已存在，你可以直接登录",
		})
		return
	}

	var user models.User

	user.Email = userData.Email
	user.Password = userData.Password
	user.Phone = userData.Phone
	user.NickName = userData.NickName
	user.Profile = userData.Profile

	user.Role = models.UserRoleRegistered
	user.CreatedAt = time.Now()

	h := md5.New()
	h.Write([]byte(bson.NewObjectId().Hex()))
	cipherStr := h.Sum(nil)
	user.Salt = hex.EncodeToString(cipherStr)

	h.Reset()
	h.Write([]byte(bson.NewObjectId().Hex()))
	cipherStr = h.Sum(nil)
	user.CheckCode = hex.EncodeToString(cipherStr)

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

	util.SendRegisterMail(user.Email, user.CheckCode)

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
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
			c.JSON(http.StatusOK, gin.H{
				"status": 403,
				"msg":    "账号不存在",
			})
			return
		}

		if user.Password == "" || user.Password != loginReq.Password {
			c.JSON(http.StatusOK, gin.H{
				"status": 10001,
				"msg":    "邮箱和密码不匹配，请重新输入",
			})
			return
		}

		if user.Role < models.UserRoleChecked {
			c.JSON(http.StatusOK, gin.H{
				"status": 10002,
				"msg":    "请验证邮箱",
			})
			return
		}

		generateToken(c, user)
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

	status := 0
	data := LoginResult{
		User:  user,
		Token: token,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": status,
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
		NickName string `json:"nickName" bson:"nickName"`
	}
	au := AnonymousUser{NickName: idiom.Word}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "Success",
		"data":   au,
	})
}

func CheckUser(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	var req struct {
		SN string `json:"sn" bson:"sn"`
	}

	if c.BindJSON(&req) != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "无效验证码",
		})
		return
	}
	checkCode := req.SN

	var user models.User
	query := bson.M{
		"checkCode": checkCode,
	}

	err := db.C(models.CollectionUser).
		Find(query).
		One(&user)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "无效验证码",
		})
		return
	}

	err = db.C(models.CollectionUser).Update(bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"role": models.UserRoleChecked}})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "无效验证码",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "验证成功，<a href=\"/user-profile\">点击登录</a>",
	})
}

// Login 登录
func SendCheckCode(c *gin.Context) {

	var req struct {
		Email string `json:"email" bson:"email"`
	}

	if c.BindJSON(&req) != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "Json parse error!",
		})
	}

	db := c.MustGet("db").(*mgo.Database)

	var user models.User
	query := bson.M{
		"email": req.Email,
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

	util.SendRegisterMail(user.Email, user.CheckCode)

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "success!",
	})
}
