package dbmongo

import "go-starter-project/internal/app/database"

type binanceSymbol struct{}

func GetBinanceSymbolDatabase() database.BinanceSymbolDatabase {
	binanceSymbol := &binanceSymbol{}
	return binanceSymbol
}
