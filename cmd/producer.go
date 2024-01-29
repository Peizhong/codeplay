package cmd

import (
	"context"
	"os"
	"syscall"

	"github.com/oklog/run"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(producerCmd)
}

var producerCmd = &cobra.Command{
	Use:   "producer",
	Short: "run kafka event producer",
	Run: func(cmd *cobra.Command, args []string) {
		rootCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		g := &run.Group{}
		g.Add(run.SignalHandler(rootCtx, os.Interrupt, syscall.SIGTERM))

		g.Run()
	},
}
