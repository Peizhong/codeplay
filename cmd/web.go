package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/google/gops/agent"
	"github.com/oklog/run"
	"github.com/peizhong/codeplay/config"
	"github.com/peizhong/codeplay/pkg/cache"
	"github.com/peizhong/codeplay/pkg/logger"
	"github.com/peizhong/codeplay/service/naming"
	"github.com/peizhong/codeplay/service/web"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(webCmd)
}

func initModules() error {
	if config.C.RedisAddr != "" {
		if err := cache.InitShareCache(config.C.RedisAddr, config.C.RedisPassword); err != nil {
			return err
		}
	}
	return nil
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "run web service",
	Run: func(cmd *cobra.Command, args []string) {
		rootCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := initModules(); err != nil {
			logger.Sugar().Fatalf("init modules error: %v", err)
		}

		g := &run.Group{}
		g.Add(run.SignalHandler(rootCtx, os.Interrupt, syscall.SIGTERM))

		r := gin.Default()
		web.RegisterRoutes(r)
		srv := &http.Server{
			Handler: r,
			Addr:    fmt.Sprintf(":%d", config.C.HttpPort),
		}
		g.Add(func() error {
			logger.Sugar().Infoln("start web service on", srv.Addr)
			return srv.ListenAndServe()
		}, func(err error) {
			srv.Shutdown(context.TODO())
		})
		if config.C.GetFeature("enable_gops").Bool() {
			g.Add(func() error {
				if err := agent.Listen(agent.Options{}); err != nil {
					return err
				}
				logger.Sugar().Infoln("gops service started")
				<-rootCtx.Done()
				return nil
			}, func(err error) {
				cancel()
			})
		}
		if config.C.RedisAddr != "" {
			naming.RegisterService(g, config.C.Hostname)
		}
		g.Run()
	},
}
