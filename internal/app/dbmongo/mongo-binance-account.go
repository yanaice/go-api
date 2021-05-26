package dbmongo

import "go-starter-project/internal/app/database"

type binanceAccount struct{}

func GetBinanceAccountDatabase() database.BinanceAccountDatabase {
	binanceAccount := &binanceAccount{}
	return binanceAccount
}
