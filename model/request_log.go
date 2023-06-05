// Package model -> model
package model

import "time"

type RequestLog struct {
	Latency      time.Duration // time for accessing
	StatusCode   int
	ClientIP     string
	Method       string
	RelativePath string
	UserAgent    string
}
