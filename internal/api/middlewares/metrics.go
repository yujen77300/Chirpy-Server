package middlewares

import (
	"net/http"
	"sync/atomic"
)

type MetricsMiddleware struct {
	hits *atomic.Int32
}

func NewMetricsMiddleware(hits *atomic.Int32) *MetricsMiddleware {
	return &MetricsMiddleware{
		hits: hits,
	}
}

func (m *MetricsMiddleware) MetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.hits.Add(1)
		next.ServeHTTP(w, r)
	})
}
