package model

type Tag struct {
	ID   string `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name" json:"name"`
}
