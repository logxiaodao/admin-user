package middleware

import (
	"net/http"
)

const (
	seconds = 1
	total   = 100
	quota   = 5
)

type PeriodLimitMiddleware struct {
}

func NewPeriodLimitMiddleware() *PeriodLimitMiddleware {
	return &PeriodLimitMiddleware{}
}

func (p *PeriodLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO 这里可以自定义限流逻辑
		next(w, r)

	}
}
