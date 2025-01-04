package middlewares

import (
	"net/http"

	"golang.org/x/time/rate"
	"pizzagoland/utils"
)

var limiter = rate.NewLimiter(1, 3) // 1 запрос в секунду, до 3 запросов в буфере

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			utils.Logger.Warn("Rate limit exceeded")
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
