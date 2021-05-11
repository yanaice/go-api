package main

import (
	"github.com/gin-gonic/gin"
	"go-starter-project/internal/app/cacheredis"
	appConf "go-starter-project/internal/app/config"
	"go-starter-project/internal/app/dbmongo"
	"go-starter-project/internal/app/handler"
	"go-starter-project/internal/app/server"
	"go-starter-project/internal/app/service"
	"go-starter-project/pkg/config"
	"log"
	"os"
)

func main() {
	confFileList := os.Getenv("CONFIG_FILES")
	err := config.Init(confFileList, &appConf.Conf)
	if err != nil {
		log.Panic("Cannot read config files")
	}

	dbmongo.Init()
	tagDB := dbmongo.GetTagDatabsase()
	staffDB := dbmongo.GetStaffDatabsase()

	cacheredis.Init()
	sessions := cacheredis.GetUserSessionsCache()
	blacklist := cacheredis.GetTokenBlacklistCache()

	r := gin.New()
	h := handler.Init(r)
	h.TagSvc = service.TagServiceInit(tagDB)
	h.StaffSvc = service.StaffServiceInit(staffDB, sessions)
	h.TokenSvc = service.TokenServiceInit(blacklist, h.StaffSvc)

	svr := server.NewServer(r)
	svr.ListenAndServe()
}
