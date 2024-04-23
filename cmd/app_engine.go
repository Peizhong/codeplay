package cmd

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
	"github.com/peizhong/codeplay/pkg/logger"
	"github.com/peizhong/codeplay/service/app_engine"
	"github.com/peizhong/codeplay/service/web"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(appEngineCmd)
}

var appEngineCmd = &cobra.Command{
	Use:   "app_engine",
	Short: "run app in k8s",
	Run: func(cmd *cobra.Command, args []string) {
		if err := app_engine.InitInClusterK8sClient(); err != nil {
			logger.Sugar().Fatalln("k8s init error")
		}
		g := &run.Group{}
		g.Add(run.SignalHandler(context.Background(), os.Interrupt, syscall.SIGTERM))
		r := gin.Default()
		r.GET("/-/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{})
		})
		web.RegisterBasicApi(r)
		app_engine.RegisterHttpApiRoutes(r)
		httpSrv := &http.Server{
			Handler: r,
			Addr:    ":8080",
		}
		g.Add(func() error {
			return httpSrv.ListenAndServe()
		}, func(err error) {
			timeout, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			httpSrv.Shutdown(timeout)
		})
		g.Run()
	},
}
