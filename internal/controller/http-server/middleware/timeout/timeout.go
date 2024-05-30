package timeoutMiddleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func New(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			if ctx.Err() == context.DeadlineExceeded {

				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}

			cancel()
		}()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
