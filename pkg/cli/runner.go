package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/tfprovidercheck/pkg/controller"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

type Runner struct {
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	LDFlags *LDFlags
	LogE    *logrus.Entry
}

type LDFlags struct {
	Version string
	Commit  string
	Date    string
}

func (l *LDFlags) VersionString() string {
	if l == nil {
		return "unknown"
	}
	if l.Version == "" {
		if l.Date == "" {
			return "unknown"
		}
		return fmt.Sprintf("(%s)", l.Date)
	}
	if l.Date == "" {
		return l.Version
	}
	return fmt.Sprintf("%s (%s)", l.Version, l.Date)
}

func (r *Runner) Run(ctx context.Context, args ...string) error {
	app := &cli.App{
		Name:  "tfprovidercheck",
		Usage: "Censor Terraform Providers",
		CustomAppHelpTemplate: `tfprovidercheck - Censor Terraform Providers

https://github.com/suzuki-shunsuke/tfprovidercheck

Usage:
  tfprovidercheck [<options>]

Options:
  -help, -h     Show help
  -version, -v  Show version
	-config, -c   Configuration file path
`,
		Version: r.LDFlags.VersionString(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Configuration file path",
			},
		},
		Action: r.run,
	}

	return app.RunContext(ctx, args) //nolint:wrapcheck
}

func (r *Runner) run(c *cli.Context) error {
	param := &controller.ParamRun{
		ConfigFilePath: c.String("config"),
	}

	if param.ConfigFilePath == "" {
		param.ConfigBody = os.Getenv("TFPROVIDERCHECK_CONFIG_BODY")
		if param.ConfigBody == "" {
			param.ConfigFilePath = os.Getenv("TFPROVIDERCHECK_CONFIG")
		}
	}

	if term.IsTerminal(0) {
		return errors.New(`stdin is missing. Please pass the result of "terraform version -json" to stdin`)
	}

	vout := &controller.TerraformVersionOutput{}
	if err := json.NewDecoder(r.Stdin).Decode(vout); err != nil {
		return fmt.Errorf(`parse stdin as the output of "terraform version -json": %w`, err)
	}

	ctrl := controller.New(afero.NewOsFs())
	return ctrl.Run(c.Context, r.LogE, param, vout) //nolint:wrapcheck
}
