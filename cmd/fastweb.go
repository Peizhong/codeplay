package cmd

import (
	"context"
	"os"
	"syscall"

	"github.com/oklog/run"
	"github.com/peizhong/codeplay/service/fastweb"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fastwebCmd)
}

var fastwebCmd = &cobra.Command{
	Use:   "fastweb",
	Short: "run web service using fasthttp",
	Run: func(cmd *cobra.Command, args []string) {
		rootCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		g := &run.Group{}
		g.Add(run.SignalHandler(rootCtx, os.Interrupt, syscall.SIGTERM))
		fastweb.RegisterRuntime(g)
		g.Run()
	},
}
