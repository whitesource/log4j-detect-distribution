package cmd

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/whitesource/log4j-detect/cmd/clioptions"
	"github.com/whitesource/log4j-detect/cmd/scan"
	"go.uber.org/zap"
	"os"
)

// Nop returns an empty logger, that does not print anything anywhere
func Nop() logr.Logger {
	nop := zerolog.Nop()
	return zerologr.New(&nop)
}

// NewCmdRoot returns a new cobra.Command implementing the root command
func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use: "log4j-detect",
	}

	cmd.CompletionOptions.DisableDefaultCmd = true

	var logger logr.Logger
	if os.Getenv("SECRET_WHITESOURCE_DEBUG_MODE") == "true" {
		l, _ := zap.NewDevelopment()
		logger = zapr.NewLogger(l)
		logger.Info("SECRET WHITESOURCE DEBUG MODE ACTIVATED!!!")
	} else {
		logger = Nop()
	}

	streams := clioptions.StandardIOStreams()

	cmd.AddCommand(
		scan.NewCmdScan(logger, streams),
	)

	return cmd
}
