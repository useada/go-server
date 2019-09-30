package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	// CollectionCustomer holds the name of the customer collection
	CollectionCustomer = "customers"
)

// Customer model
type Customer struct {
	ID               bson.ObjectId   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name             string          `json:"name" bson:"name"`
	Description      string          `json:"description,omitempty" bson:"description,omitempty"`
	Address          string          `json:"address,omitempty" bson:"address,omitempty"`
	Contact          string          `json:"contact,omitempty" bson:"contact,omitempty"`
	Phone            string          `json:"phone,omitempty" bson:"phone,omitempty"`
	OrganizationID   bson.ObjectId   `json:"organizationId" bson:"organizationId"`
	ProductID        []bson.ObjectId `json:"productId,omitempty" bson:"productId,omitempty"`
	CreatedBy        bson.ObjectId   `json:"createdBy" bson:"createdBy"`
	OrganizationName string          `json:"organizationName" bson:"organizationName"`
	ProductName      string          `json:"productName" bson:"productName"`
	CreatedName      string          `json:"createdName" bson:"createdName"`
	CreatedAt        time.Time       `json:"createdAt" bson:"createdAt"`
	MemberCount      int             `json:"memberCount,omitempty" bson:"memberCount,omitempty"`
	ProductCount     int             `json:"productCount,omitempty" bson:"productCount,omitempty"`
	DeviceCount      int             `json:"deviceCount,omitempty" bson:"deviceCount,omitempty"`
}
