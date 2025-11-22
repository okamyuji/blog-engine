package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter レートリミッター
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int
	burst    int
}

// visitor 訪問者の情報
type visitor struct {
	lastSeen time.Time
	count    int
	window   time.Time
}

// NewRateLimiter 新しいRateLimiterを作成
func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
	}

	// 定期的にクリーンアップ
	go rl.cleanupVisitors()

	return rl
}

// Limit レート制限を適用するミドルウェア
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		if !rl.allow(ip) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// allow リクエストを許可するかチェック
func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	v, exists := rl.visitors[ip]

	if !exists {
		rl.visitors[ip] = &visitor{
			lastSeen: now,
			count:    1,
			window:   now,
		}
		return true
	}

	// ウィンドウのリセット
	if now.Sub(v.window) > time.Minute {
		v.count = 1
		v.window = now
		v.lastSeen = now
		return true
	}

	// レート制限チェック
	if v.count >= rl.rate {
		v.lastSeen = now
		return false
	}

	v.count++
	v.lastSeen = now
	return true
}

// cleanupVisitors 古い訪問者情報を削除
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 10*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}
