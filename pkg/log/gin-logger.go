package log

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var timeFormat = "02/Jan/2006:15:04:05 -0700"

func isWebsocket(c *gin.Context) bool {
	return c.Request.Header.Get("Upgrade") == "websocket"
}

// Logger is the logrus logger handler
func GinLogger() gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		xForwardedFor := c.Request.Header.Get("x-forwarded-for")
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		handler := strings.Split(c.HandlerName(), ".")
		shortHandlerName := handler[len(handler)-1]
		if isWebsocket(c) {
			WithFields(Fields{
				"hostname":      hostname,
				"clientIP":      clientIP,
				"xForwardedFor": xForwardedFor,
				"method":        c.Request.Method,
				"path":          path,
				"referer":       referer,
				"userAgent":     clientUserAgent,
				"webSocket":     "yes",
			}).Infof(
				"%s - %s [%s] \"%s %s\" \"%s\" \"%s\" (WebSocket connection initiation)",
				clientIP,
				hostname,
				time.Now().Format(timeFormat),
				c.Request.Method,
				path,
				referer,
				clientUserAgent,
			)
		}
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		entry := WithFields(Fields{
			"hostname":      hostname,
			"statusCode":    statusCode,
			"latency":       latency, // time to process
			"clientIP":      clientIP,
			"xForwardedFor": xForwardedFor,
			"method":        c.Request.Method,
			"path":          path,
			"referer":       referer,
			"dataLength":    dataLength,
			"userAgent":     clientUserAgent,
			"handler":		 shortHandlerName,
		})

		if isWebsocket(c) {
			entry = entry.WithField("webSocket", "yes")
		}

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("%s - %s [%s] \"%s %s\" %d %d \"%s\" \"%s\" (%dms)", clientIP, hostname, time.Now().Format(timeFormat), c.Request.Method, path, statusCode, dataLength, referer, clientUserAgent, latency)
			switch {
			case statusCode >= 500:
				entry.Error(msg)
			case statusCode >= 400:
				entry.Warn(msg)
			default:
				entry.Info(msg)
			}
		}
	}
}
