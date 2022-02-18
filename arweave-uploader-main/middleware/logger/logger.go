package logger

import (
	"time"

	"ccian.cc/really/arweave-api/pkg/logging"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const (
	gContextLoggerName = "really-arweave-api-logger"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set requestID
		reqLogger := logging.WithFields(logrus.Fields{
			"requestID": uuid.NewV4(),
		})
		c.Set(gContextLoggerName, reqLogger)
		reqLogger.WithField("path", c.Request.URL.Path).Debug("request begin")
		c.Next()
	}
}

func GetContextLogger(c *gin.Context) *logrus.Entry {
	cLogger, created := getLogger(c)
	if created {
		c.Set(gContextLoggerName, cLogger)
	}

	return cLogger
}

func getLogger(c *gin.Context) (logger *logrus.Entry, created bool) {
	log, exists := c.Get(gContextLoggerName)
	if !exists {
		return injectReqPathLogger(c), true
	}

	logger, exists = log.(*logrus.Entry)
	if !exists {
		return injectReqPathLogger(c), true
	}

	return logger, false
}

func Time() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)
		status := c.Writer.Status()

		reqLogger, _ := getLogger(c)

		if len(c.Errors) > 0 {
			reqLogger.WithFields(logrus.Fields{
				"latency": latency,
				"status":  status,
				"errors":  c.Errors,
			}).Info("request failed")
		} else {
			reqLogger.WithFields(logrus.Fields{
				"latency": latency,
				"status":  status,
			}).Debug("request done")
		}
	}
}

func injectReqPathLogger(c *gin.Context) *logrus.Entry {
	return logging.WithFields(logrus.Fields{"path": c.Request.URL.Path})
}
