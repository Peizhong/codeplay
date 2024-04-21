package cmd

import (
	"context"
	"os"
	"syscall"

	"github.com/oklog/run"
	"github.com/peizhong/codeplay/pkg/logger"
	"github.com/peizhong/codeplay/service/app_engine"
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
		rootCtx, cancel := context.WithCancel(context.Background())
		defer cancel()
		g := &run.Group{}
		g.Add(run.SignalHandler(rootCtx, os.Interrupt, syscall.SIGTERM))
		g.Run()
	},
}
