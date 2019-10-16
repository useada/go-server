package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"serve/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func JWTPrepare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			return
		}

		log.Print("get token: ", token)

		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)
	}
}

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsI, ok := c.Get("claims")
		var claims *CustomClaims
		if ok {
			claims = claimsI.(*CustomClaims)
		}

		if claims == nil {
			c.JSON(http.StatusOK, gin.H{
				"status": 403,
				"msg":    "请重新登录",
			})
			c.Abort()
			return
		}
	}
}

// JWTAuth 中间件，检查token
func JWTAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("claims").(*CustomClaims)

		db := c.MustGet("db").(*mgo.Database)

		var user models.User
		query := bson.M{
			"_id": claims.ID,
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

		if user.Role != models.UserRoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"status": 403,
				"msg":    "没有权限",
			})

			c.Abort()
			return
		}
	}
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// 一些常量
var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "funnylink.net-sk-2019"
)

// 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	ID bson.ObjectId `json:"_id"`
	//Email string        `json:"email"`
	jwt.StandardClaims
}

// 新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// 获取signKey
func GetSignKey() string {
	return SignKey
}

// 这是SignKey
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

func ExpiresAt() int64 {
	return time.Now().Add(24 * time.Hour).Unix()
}

// CreateToken 生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析Token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		fmt.Println(err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		//claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		claims.StandardClaims.ExpiresAt = ExpiresAt()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
