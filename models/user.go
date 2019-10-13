package models

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// CollectionUser holds the name of the user collection
	CollectionUser = "users"

	UserRoleAdmin      = 9
	UserRoleRegistered = 1
	UserRoleChecked    = 2
)

type CaptchaData struct {
	CaptchaID   string `json:"captchaID" bson:"captchaID"`
	VerifyValue string `json:"verifyValue" bson:"verifyValue"`
}

// User model
type User struct {
	ID    bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Email string        `json:"email" bson:"email"`
	// 9: admin 0: common 1: organization 2: customer
	Type             int           `json:"type" bson:"type"`
	Role             int           `json:"role" bson:"role"`
	Password         string        `json:"-" binding:"required" bson:"password"`
	Phone            string        `json:"phone,omitempty" bson:"phone,omitempty"`
	NickName         string        `json:"nickName,omitempty" bson:"nickName,omitempty"`
	Profile          string        `json:"profile,omitempty" bson:"profile,omitempty"`
	Salt             string        `json:"-" bson:"salt,omitempty"`
	CheckCode        string        `json:"-" bson:"checkCode,omitempty"`
	CreatedBy        bson.ObjectId `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	CreatedAt        time.Time     `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	OrganizationID   bson.ObjectId `json:"organizationId,omitempty" bson:"organizationId,omitempty"`
	CustomerID       bson.ObjectId `json:"customerId,omitempty" bson:"customerId,omitempty"`
	OrganizationName string        `json:"organizationName" bson:"organizationName"`
	CreatedName      string        `json:"createdName" bson:"createdName"`
	CustomerName     string        `json:"customerName" bson:"customerName"`
	OrgCount         int           `json:"orgCount,omitempty" bson:"orgCount,omitempty"`
	CustomerCount    int           `json:"customerCount,omitempty" bson:"customerCount,omitempty"`
	ProductCount     int           `json:"productCount,omitempty" bson:"productCount,omitempty"`
	DeviceCount      int           `json:"deviceCount,omitempty" bson:"deviceCount,omitempty"`
}

// Login param
type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetUserByEmail(db *mgo.Database, email string) (*User, error) {
	var user User

	err := db.C(CollectionUser).
		Find(bson.M{"email": email}).
		One(&user)
	if err != nil {
		if !strings.Contains(err.Error(), `not found`) {
			return nil, err
		} else {
			return nil, nil
		}
	}

	return &user, nil
}

func GetUserByID(db *mgo.Database, id bson.ObjectId) (*User, error) {
	var user User

	err := db.C(CollectionUser).
		Find(bson.M{"_id": id}).
		One(&user)
	if err != nil {
		if !strings.Contains(err.Error(), `not found`) {
			return nil, err
		} else {
			return nil, nil
		}
	}

	return &user, nil
}
