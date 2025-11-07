package middleware

import (
	"net/http"

	"github.com/half-nothing/simple-fsd/src/utils"
	"github.com/labstack/echo/v4"
)

func CombinedKeyFunc(c echo.Context) string {
	return c.RealIP() + "|" + c.Path()
}

// RateLimitMiddleware Echo 限流中间件
func RateLimitMiddleware(limiter *utils.SlidingWindowLimiter, keyFunc func(c echo.Context) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := keyFunc(c)

			if !limiter.Allow(key) {
				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "请求次数过多, 请稍后再试",
					"data":    nil,
				})
			}

			return next(c)
		}
	}
}
