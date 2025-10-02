// Package config
package config

import (
	"errors"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"time"
)

type HttpServerLimit struct {
	RateLimit         int           `json:"rate_limit"`
	RateLimitWindow   string        `json:"rate_limit_window"`
	RateLimitDuration time.Duration `json:"-"`
}

func defaultHttpServerLimit() *HttpServerLimit {
	return &HttpServerLimit{
		RateLimit:       15,
		RateLimitWindow: "1m",
	}
}

func (config *HttpServerLimit) checkValid(_ log.LoggerInterface) *ValidResult {
	if duration, err := time.ParseDuration(config.RateLimitWindow); err != nil {
		return ValidFailWith(errors.New("invalid json field http_server.rate_limit_window, %v"), err)
	} else {
		config.RateLimitDuration = duration
	}

	return ValidPass()
}
