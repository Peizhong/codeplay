package cmd

import (
	"os"

	"github.com/peizhong/codeplay/pkg/logger"
	"github.com/peizhong/codeplay/pkg/stat"
	"github.com/spf13/cobra"
)

var (
	BuildDate string
	GitBranch string
	GitCommit string
)

var rootCmd = &cobra.Command{
	Use:   "codeplay",
	Short: "play code",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Sugar().Infoln("have fun!!")
	},
}

func Execute() {
	logger.InitLogger()
	stat.RegisterSystemMetrics(GitBranch, GitCommit, BuildDate)

	err := rootCmd.Execute()
	if err != nil {
		logger.Sugar().Warnln("root cmd exited with err", err.Error())
	}
	logger.Flush()
	if err != nil {
		os.Exit(1)
	}
}
