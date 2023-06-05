// Package middleware -> delivery/middleware
package middleware

import (
	"log"
	"os"
	"time"

	"be-b-impact.com/csr/config"
	"be-b-impact.com/csr/model"
	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
)

func LogRequestMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(cfg.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	// create a new instance logger

	// set output file
	logger.SetOutput(file)
	return func(c *gin.Context) {
		c.Next()
		latency := time.Since(time.Now())
		requestLog := model.RequestLog{
			Latency:      latency,
			StatusCode:   c.Writer.Status(), // statusCode
			ClientIP:     c.ClientIP(),
			Method:       c.Request.Method,
			RelativePath: c.Request.URL.Path,
			UserAgent:    c.Request.UserAgent(),
		}
		switch {
		case c.Writer.Status() >= 500:
			logger.Error(requestLog)
		case c.Writer.Status() >= 400:
			logger.Warn(requestLog)
		default:
			logger.Info(requestLog) // >= 100 .. 200
		}
	}
}
