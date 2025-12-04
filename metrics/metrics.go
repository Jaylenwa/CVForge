package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	reqCounter  = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "http_requests_total"}, []string{"method", "path", "status"})
	reqDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "http_request_duration_seconds"}, []string{"method", "path"})
)

func init() { prometheus.MustRegister(reqCounter, reqDuration) }

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		status := c.Writer.Status()
		reqCounter.WithLabelValues(c.Request.Method, c.FullPath(), toStatus(status)).Inc()
		reqDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(time.Since(start).Seconds())
	}
}

func Handler() gin.HandlerFunc { return gin.WrapH(promhttp.Handler()) }

func toStatus(s int) string { return strconv.Itoa(s) }
