package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
)

var ratelimitSql = redis.NewScript(`
    local current = redis.call("INCR", KEYS[1])
    if current == 1 then
        redis.call("PEXPIRE", KEYS[1], ARGV[1])
    end
    return current
`)

var (
	ratelimitEvents = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_ratelimit_exceeded_total",
		Help: "Total number of requests blocked by rate limiter",
	}, []string{"method", "path"})
)

func RateLimiter(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:%s", ip)

		val, err := ratelimitSql.Run(ctx, rdb, []string{key}, window.Milliseconds()).Int64()
		if err != nil {
			slog.Error("Redis Lua script error", "error", err)
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", max(0, int64(limit)-val)))

		if val > int64(limit) {
			path := c.FullPath()
			if path == "" {
				path = "unmatched"
			}
			ratelimitEvents.WithLabelValues(c.Request.Method, path).Inc()
			slog.Warn("Rate limit exceeded", "ip", ip, "count", val)

			c.String(http.StatusTooManyRequests, "Too Many Requests\n")
			c.Abort()
			return
		}

		c.Next()
	}
}
