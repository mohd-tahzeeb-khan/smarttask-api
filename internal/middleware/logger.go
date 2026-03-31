package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)

		if query != "" {
			path = path + "?" + query
		}

		log.Printf("%s%s%s | %s%-7s%s | %13v | %-15s | %s%d%s | %s",
			"\033[1m", time.Now().Format("2006/01/02 - 15:04:05"), "\033[0m",
			methodColor, method, "\033[0m",
			latency,
			clientIP,
			statusColor, statusCode, "\033[0m",
			path,
		)
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 500:
		return "\033[31m"
	case code >= 400:
		return "\033[33m"
	case code >= 300:
		return "\033[36m"
	default:
		return "\033[32m"
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return "\033[34m"
	case "POST":
		return "\033[32m"
	case "PUT", "PATCH":
		return "\033[33m"
	case "DELETE":
		return "\033[31m"
	default:
		return "\033[37m"
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v", err)
				c.AbortWithStatusJSON(500, gin.H{
					"success": false,
					"error":   fmt.Sprintf("internal server error: %v", err),
				})
			}
		}()
		c.Next()
	}
}
