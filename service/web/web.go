package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/peizhong/codeplay/config"
	_ "github.com/peizhong/codeplay/gen/openapiv2/evaluator"
	_ "github.com/peizhong/codeplay/gen/swagger/webapi"
	rpc_evaluator "github.com/peizhong/codeplay/rpc/evaluator"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginswagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func registerEvaluatorGateway(r *gin.Engine, rpcAddr string) {
	// http://localhost:3000/swagger/evaluator/index.html
	r.GET("/swagger/evaluator/*any", ginswagger.CustomWrapHandler(&ginswagger.Config{
		URL:          "doc.json",
		InstanceName: "evaluator",
	}, swaggerfiles.Handler))

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := rpc_evaluator.RegisterEvaluatorHandlerFromEndpoint(context.Background(), gwmux, rpcAddr, opts); err != nil {
		panic(fmt.Errorf("failed to register gRPC gateway: %v", err))
	}
	r.Any("/rpc/evaluator/*any", gin.WrapH(gwmux))
}

func registerSwaggerApi(r *gin.Engine) {
	// http://localhost:3000/swagger/webapi/index.html#
	r.GET("/swagger/webapi/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
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
	registerEvaluatorGateway(r, ":9090")
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
		"hostname": config.C.Hostname,
	})
}
