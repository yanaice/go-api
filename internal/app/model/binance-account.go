package model

type BinanceAccount struct {
	ID       string `bson:"-" json:"id"`
	Username string `json:"username"`
	Type     string `json:"type"`
}
