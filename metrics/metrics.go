package metrics

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

type AppMetrics struct {
	RequestCount  int64
	ErrorCount    int64
	DBErrorCount  int64
	mu            sync.Mutex
	totalLatency  time.Duration
	requestSample int64
}

var Instance = &AppMetrics{}

func (m *AppMetrics) IncRequest()  { atomic.AddInt64(&m.RequestCount, 1) }
func (m *AppMetrics) IncError()    { atomic.AddInt64(&m.ErrorCount, 1) }
func (m *AppMetrics) IncDBError()  { atomic.AddInt64(&m.DBErrorCount, 1) }

func (m *AppMetrics) RecordLatency(d time.Duration) {
	m.mu.Lock()
	m.totalLatency += d
	m.requestSample++
	m.mu.Unlock()
}

func (m *AppMetrics) AvgLatency() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.requestSample == 0 {
		return "0ms"
	}
	avg := m.totalLatency / time.Duration(m.requestSample)
	return avg.Round(time.Microsecond).String()
}

func (m *AppMetrics) Snapshot() gin.H {
	return gin.H{
		"request_count":    atomic.LoadInt64(&m.RequestCount),
		"error_count":      atomic.LoadInt64(&m.ErrorCount),
		"db_error_count":   atomic.LoadInt64(&m.DBErrorCount),
		"avg_latency":      m.AvgLatency(),
	}
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		Instance.IncRequest()

		c.Next()

		Instance.RecordLatency(time.Since(start))
		if c.Writer.Status() >= 400 {
			Instance.IncError()
		}
	}
}
