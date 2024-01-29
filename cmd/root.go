package cmd

import (
	"os"

	"github.com/peizhong/codeplay/config"
	"github.com/peizhong/codeplay/pkg/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "codeplay",
	Short: "play code",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Sugar().Infoln("have fun!!")
	},
}

func Execute() {
	config.Init()

	logger.InitLogger()
	defer logger.Flush()

	if err := rootCmd.Execute(); err != nil {
		logger.Sugar().Warnln("root cmd exited with err", err.Error())
		os.Exit(1)
	}
	logger.Sugar().Infoln("root cmd exited :)")
}
