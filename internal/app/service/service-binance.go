package service

import (
	"go-starter-project/internal/app/database"
)

type BinanceService interface {
	GetBinanceAccountAndSymbol(username string) (string, error)
}

type serviceBinanceImpl struct {
	dba database.BinanceAccountDatabase
	dbs database.BinanceSymbolDatabase
}

func BinanceServiceInit(dba database.BinanceAccountDatabase, dbs database.BinanceSymbolDatabase) BinanceService {
	return &serviceBinanceImpl{dba: dba, dbs: dbs}
}

func (s *serviceBinanceImpl) GetBinanceAccountAndSymbol(username string) (string, error) {

	return "HELLO WORLD " + username, nil
}
