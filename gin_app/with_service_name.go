package gin_app

import (
	"github.com/gin-gonic/gin"
)

func WithServiceName(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("serviceName", name)
		c.Next()
	}
}
