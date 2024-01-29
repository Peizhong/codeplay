package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peizhong/codeplay/config"
	_ "github.com/peizhong/codeplay/gen/swagger/webapi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginswagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func registerSwaggerApi(r *gin.Engine) {
	// local: http://localhost:3000/-/swagger/webapi/index.html#
	// vm: http://10.10.10.1:30000/-/swagger/webapi/index.html#
	r.GET("/-/swagger/webapi/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
}

func registerBasicApi(r *gin.Engine) {
	r.GET("/-/healthy", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{}) })
	r.GET("/-/ready", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{}) })
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func registerLogicApi(r *gin.Engine) {
	g := r.Group("/api", mwMetric())
	g.POST("/echo", echo)
}

func RegisterRoutes(r *gin.Engine) {
	registerSwaggerApi(r)
	registerBasicApi(r)
	registerLogicApi(r)
}

type EchoRequest struct {
	Message string `json:"message" binding:"required"`
}

// @Summary Echo
// @Description response message
// @Tags cloud
// @Accept json
// @Produce json
// @Param req body EchoRequest true "Request"
// @Success 200 {object} gin.H
// @Router /api/echo [post]
func echo(c *gin.Context) {
	var req EchoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  req.Message,
		"hostname": config.HostName,
	})
}
