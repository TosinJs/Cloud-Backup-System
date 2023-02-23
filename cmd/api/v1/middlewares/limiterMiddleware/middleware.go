package limiterMiddleware

import (
	"net/http"
	"time"
	"tosinjs/cloud-backup/internal/entity/responseEntity"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type limiterMiddleware struct{}

func New() limiterMiddleware {
	return limiterMiddleware{}
}

func (l limiterMiddleware) LimitFileSize(limit int64) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var w http.ResponseWriter = c.Writer
		c.Request.Body = http.MaxBytesReader(w, c.Request.Body, limit)
		if err := c.Request.ParseMultipartForm(limit); err != nil {
			c.AbortWithStatusJSON(
				http.StatusRequestEntityTooLarge,
				responseEntity.BuildErrorResponseObject(
					http.StatusRequestEntityTooLarge,
					"200 MB Size Limit. Reduce the filesize and try again",
					c.FullPath(),
				))
			return
		}
		c.Next()
	}
	return fn
}

func (l limiterMiddleware) LimitRequests(allowedReq int) gin.HandlerFunc {
	rateLimiter := rate.NewLimiter(rate.Every(time.Minute), allowedReq)
	fn := func(c *gin.Context) {
		if !rateLimiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, responseEntity.BuildErrorResponseObject(
				http.StatusTooManyRequests, "Too Many Requests, Please wait and try again", c.FullPath(),
			))
			return
		}
		c.Next()
	}
	return fn
}
