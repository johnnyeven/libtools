package gin_app

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/johnnyeven/libtools/courier/status_error"
	"github.com/johnnyeven/libtools/duration"
	"github.com/johnnyeven/libtools/log/context"
)

var (
	REQUEST_ID_NAME = "x-request-id"
	ProjectRef      = os.Getenv("PROJECT_REF")
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		d := duration.NewDuration()

		reqID := c.Request.Header.Get(REQUEST_ID_NAME)

		if reqID == "" {
			reqID = uuid.New().String()
		}

		c.Set(REQUEST_ID_NAME, reqID)
		c.Header("X-Reversion", ProjectRef)

		context.SetLogID(reqID)

		fields := logrus.Fields{
			"tag":       "access",
			"log_id":    reqID,
			"remote_ip": c.ClientIP(),
			"method":    c.Request.Method,
			"pathname":  c.Request.URL.Path,
		}

		c.Next()

		fields["status"] = c.Writer.Status()
		fields["request_time"] = d.Get()

		logger := logrus.WithFields(fields)

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				statusErr := status_error.FromError(err.Err)
				if statusErr.Status() >= 500 {
					logger.Errorf(statusErr.Error())
				} else {
					logger.Warnf(statusErr.Error())
				}
			}
		} else {
			logger.Info("")
		}
	}
}
