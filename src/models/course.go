package models

type Course struct {
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
}
