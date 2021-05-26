package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) registerBinanceService(r *gin.RouterGroup) {
	for _, route := range h.binanceServiceRoutes() {
		r.Handle(route.Method, route.Pattern, route.Handler)
	}
}

func (h *Handler) binanceServiceRoutes() []Route {
	return []Route{
		{
			Name:        "Get binanceAccountDetail",
			Description: "Get binanceAccountDetail",
			Method:      http.MethodGet,
			Pattern:     "/binance-accounts",
			Handler:     h.Helloworld,
		},
	}
}

func (h *Handler) Helloworld(c *gin.Context) {
	result, err := h.BinanceService.GetBinanceAccountAndSymbol("YANA")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
}
