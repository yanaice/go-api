package model

type BinanceSymbol struct {
	ID               string  `bson:"-" json:"id"`
	BinanceAccountID string  `bson:"binance_account_id" json:"binance_account_id"`
	Symbol           string  `bson:"symbol" json:"symbol"`
	Stoploss         float64 `bson:"stoploss" json:"stoploss"`
}
