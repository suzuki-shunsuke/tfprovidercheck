package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/tfprovidercheck/pkg/cli"
	"github.com/suzuki-shunsuke/tfprovidercheck/pkg/log"
	"golang.org/x/term"
)

var (
	version = ""
	commit  = "" //nolint:gochecknoglobals
	date    = "" //nolint:gochecknoglobals
)

type HasExitCode interface {
	ExitCode() int
}

func main() {
	if code := core(); code != 0 {
		os.Exit(code)
	}
}

func core() int {
	logger, logLevelVar := log.New(version)
	runner := cli.Runner{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		LDFlags: &cli.LDFlags{
			Version: version,
			Commit:  commit,
			Date:    date,
		},
		Logger:      logger,
		LogLevelVar: logLevelVar,
		Env: &cli.Env{
			Config:     os.Getenv("TFPROVIDERCHECK_CONFIG"),
			ConfigBody: os.Getenv("TFPROVIDERCHECK_CONFIG_BODY"),
		},
		IsTerminal: term.IsTerminal(0),
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := runner.Run(ctx, os.Args...); err != nil {
		slogerr.WithError(logger, err).Error("tfprovidercheck failed")
		return 1
	}
	return 0
}
