package log

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-starter-project/pkg/threadlocal"
)

func GinCorrelationIdHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Request.Header.Get("X-Correlation-Id")
		if id == "" {
			id = uuid.New().String()
		}
		threadlocal.SetCorrelationID(id, func() {
			c.Next()
		})
	}
}
