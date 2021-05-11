package handler

import (
	"github.com/gin-gonic/gin"
	"go-starter-project/internal/app/service"
	api "go-starter-project/pkg/api/auth"
	"go-starter-project/pkg/auth"
	"go-starter-project/pkg/log"
)

type Route struct {
	Name        string
	Description string
	Method      string
	Pattern     string
	Handler     gin.HandlerFunc
}

type Handler struct {
	router   *gin.Engine
	TagSvc   service.TagService
	StaffSvc service.StaffService
	TokenSvc service.TokenService
}

func Init(r *gin.Engine) *Handler {

	h := &Handler{}
	h.router = r
	r.Use(log.GinCorrelationIdHandler(), log.GinLogger(), gin.Recovery())

	apiRG := h.router.Group("/")
	h.registerCommonSvc(apiRG)

	tagRG := h.router.Group("/tag")
	tagRG.Use(auth.GetMiddlewareAuthToken(api.AuthTokenAPIInit()))
	h.registerTagSvc(tagRG)

	staffRG := h.router.Group("/staff")
	staffRG.Use(auth.GetMiddlewareAuthToken(api.AuthTokenAPIInit()))
	h.registerStaffSvc(staffRG)

	tokenRG := h.router.Group("/token")
	h.registerTokenSvc(tokenRG)

	return h
}
