package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionIdiom = "idioms"
)

// Idiom model
type Idiom struct {
	ID           bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Derivation   string        `json:"derivation" bson:"derivation"`
	Example      string        `json:"example" bson:"example"`
	Explanation  string        `json:"explanation" bson:"explanation"`
	Pinyin       string        `json:"pinyin" bson:"pinyin"`
	Word         string        `json:"word" bson:"word"`
	Abbreviation string        `json:"abbreviation" bson:"abbreviation"`
}
