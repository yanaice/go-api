package main

import (
	"fmt"
	appConf "go-starter-project/internal/app/config"
	"go-starter-project/internal/app/dbmongo"
	"go-starter-project/internal/app/handler"
	"go-starter-project/internal/app/server"
	"go-starter-project/internal/app/service"
	"go-starter-project/pkg/config"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	confFileList := os.Getenv("CONFIG_FILES")
	err := config.Init(confFileList, &appConf.Conf)
	if err != nil {
		log.Panic("Cannot read config files")
	}

	dbmongo.Init()
	binanceAccountDB := dbmongo.GetBinanceAccountDatabase()
	binanceSymbolDB := dbmongo.GetBinanceSymbolDatabase()
	tagDB := dbmongo.GetTagDatabase()
	// staffDB := dbmongo.GetStaffDatabase()

	// cacheredis.Init()
	// sessions := cacheredis.GetUserSessionsCache()
	// blacklist := cacheredis.GetTokenBlacklistCache()

	r := gin.New()
	h := handler.Init(r)
	h.BinanceService = service.BinanceServiceInit(binanceAccountDB, binanceSymbolDB)
	h.TagSvc = service.TagServiceInit(tagDB)
	// h.StaffSvc = service.StaffServiceInit(staffDB, sessions)
	// h.TokenSvc = service.TokenServiceInit(blacklist, h.StaffSvc)

	svr := server.NewServer(r)
	svr.ListenAndServe()
	fmt.Println("Start server listen port:")
}
