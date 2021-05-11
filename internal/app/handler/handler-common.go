package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) registerCommonSvc(r *gin.RouterGroup) {
	for _, route := range h.CommonRoutes() {
		r.Handle(route.Method, route.Pattern, route.Handler)
	}
}

func (h *Handler) CommonRoutes() []Route {
	return []Route{
		{
			Name:        "Ping",
			Description: "",
			Method:      http.MethodGet,
			Pattern:     "/ping",
			Handler:     h.Ping,
		},
		{
			Name:        "HealthCheck",
			Description: "",
			Method:      http.MethodGet,
			Pattern:     "/healthcheck",
			Handler:     h.HealthCheck,
		},
	}
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
