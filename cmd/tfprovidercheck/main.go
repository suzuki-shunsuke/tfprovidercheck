package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
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
	logE := log.New(version)
	if err := core(logE); err != nil {
		logerr.WithError(logE, err).Fatal("tfprovidercheck failed")
	}
}

func core(logE *logrus.Entry) error {
	runner := cli.Runner{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		LDFlags: &cli.LDFlags{
			Version: version,
			Commit:  commit,
			Date:    date,
		},
		LogE: logE,
		Env: &cli.Env{
			Config:     os.Getenv("TFPROVIDERCHECK_CONFIG"),
			ConfigBody: os.Getenv("TFPROVIDERCHECK_CONFIG_BODY"),
		},
		IsTerminal: term.IsTerminal(0),
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return runner.Run(ctx, os.Args...) //nolint:wrapcheck

}
