package handler

import (
	"github.com/gin-gonic/gin"
	"go-starter-project/internal/app/model"
	"go-starter-project/pkg/auth"
	"go-starter-project/pkg/derror"
	commonhdl "go-starter-project/pkg/handler"
	"net/http"
)

func (h *Handler) registerTokenSvc(r *gin.RouterGroup) {
	for _, route := range h.TokenRoutes() {
		r.Handle(route.Method, route.Pattern, route.Handler)
	}
}

func (h *Handler) TokenRoutes() []Route {
	return []Route{
		{
			Name:        "CheckToken",
			Description: "Check token",
			Method:      http.MethodGet,
			Pattern:     "/check",
			Handler:     h.checkToken,
		},
		{
			Name:        "RevokeToken",
			Description: "Revoke token",
			Method:      http.MethodPost,
			Pattern:     "/revoke",
			Handler:     h.revokeToken,
		},
		{
			Name:        "RefreshToken",
			Description: "Refresh token",
			Method:      http.MethodPost,
			Pattern:     "/refresh",
			Handler:     h.refreshTokenPost,
		},
	}
}

func (h *Handler) checkToken(c *gin.Context) {
	token, err := auth.GetTokenFromHeader(c)
	if err != nil {
		commonhdl.HandlerReturnError(c, derror.E(err).SetCode(derror.ErrCodeUnauthorized).SetDebug("check token failed: invalid token header").Log())
		return
	}

	if err := h.TokenSvc.CheckToken(token); err != nil {
		commonhdl.HandlerReturnError(c, err)
		return
	}

	commonhdl.HandlerReturnData(c, model.Response{Status: model.ResponseStatusSuccess})
}

func (h *Handler) revokeToken(c *gin.Context) {
	var pairRequest model.TokensPair
	if err := c.ShouldBindJSON(&pairRequest); err != nil {
		commonhdl.HandlerReturnError(c, derror.E(err).SetCode(derror.ErrCodeInputError).
			SetDebug("revoke token: bind body error").Log())
		return
	}

	if pairRequest.RefreshToken == "" {
		// TODO: getting refresh token from cookies
	}

	if err := h.TokenSvc.RevokeToken(pairRequest); err != nil {
		commonhdl.HandlerReturnError(c, err)
	} else {
		// TODO: unset refresh token in cookie
		commonhdl.HandlerReturnData(c, model.Response{Status: model.ResponseStatusSuccess})
	}
}

func (h *Handler) refreshTokenPost(c *gin.Context) {
	var refreshRequest model.TokensPair
	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		commonhdl.HandlerReturnError(c, derror.E(err).SetCode(derror.ErrCodeInputError).
			SetDebug("refresh token failed: bind body error").Log())
		return
	}

	if refreshRequest.RefreshToken == "" {
		// TODO: getting refresh token from cookies
	}

	if token, err := h.TokenSvc.RefreshToken(refreshRequest); err != nil {
		commonhdl.HandlerReturnError(c, err)
	} else {
		// TODO: set cookies
		commonhdl.HandlerReturnData(c, token)
	}
}
