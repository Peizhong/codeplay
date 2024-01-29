package web

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/peizhong/codeplay/pkg/stat"
)

func mwMetric() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		stat.RecordRequest(c.FullPath(), c.Request.Method, c.Writer.Status())
		stat.RecordRequestDuration(c.FullPath(), c.Request.Method, time.Since(start).Seconds())
	}
}
