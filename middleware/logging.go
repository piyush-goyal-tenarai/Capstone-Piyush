package middleware

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type LogEntry struct {
	RequestID string `json:"request_id"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	Duration  string `json:"duration"`
	ClientIP  string `json:"client_ip"`
	Error     string `json:"error,omitempty"`
}

func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		entry := LogEntry{
			RequestID: c.GetString("request_id"),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Status:    c.Writer.Status(),
			Duration:  time.Since(start).Round(time.Microsecond).String(),
			ClientIP:  c.ClientIP(),
		}

		if len(c.Errors) > 0 {
			entry.Error = c.Errors.String()
		}

		logJSON, _ := json.Marshal(entry)
		log.Println(string(logJSON))
	}
}
