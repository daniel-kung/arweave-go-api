package method

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SupportedMethodOfAny
//
// use gin.Any to accept multiple methods, but would pass only with the valid ones.
func SupportedMethodOfAny() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet:
			// Options would be used in CORS
			c.Next()
		case http.MethodOptions:
			// return 200 for option method
			c.AbortWithStatus(http.StatusOK)
		default:
			// return 404 for any other methods
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
}
