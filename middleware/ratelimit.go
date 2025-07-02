// middleware/ratelimit.go

package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type VisitLimiter struct {
	sync.Mutex
	lastVisit map[string]time.Time
	count     map[string]int
}

func NewVisitLimiter() *VisitLimiter {
	return &VisitLimiter{
		lastVisit: make(map[string]time.Time),
		count:     make(map[string]int),
	}
}

func RateLimitMiddleware(limiter *VisitLimiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()

			limiter.Lock()
			defer limiter.Unlock()

			now := time.Now()
			lastTime, exists := limiter.lastVisit[ip]

			// اگر اولین درخواست است یا بیش از 1 دقیقه از درخواست قبلی گذشته
			if !exists || now.Sub(lastTime) > time.Minute {
				limiter.lastVisit[ip] = now
				limiter.count[ip] = 1
				return next(c)
			}

			// محدودیت 10 درخواست در دقیقه
			if limiter.count[ip] >= 10 {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "درخواست‌های بیش از حد مجاز. لطفا کمی صبر کنید.",
				})
			}

			limiter.lastVisit[ip] = now
			limiter.count[ip]++

			return next(c)
		}
	}
}
