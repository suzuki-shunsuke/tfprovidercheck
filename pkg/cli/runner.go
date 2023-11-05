package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/tfprovidercheck/pkg/controller"
	"github.com/urfave/cli/v2"
)

type Runner struct {
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	LDFlags    *LDFlags
	LogE       *logrus.Entry
	Env        *Env
	IsTerminal bool
}

type Env struct {
	Config     string
	ConfigBody string
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
		param.ConfigBody = r.Env.ConfigBody
		if param.ConfigBody == "" {
			param.ConfigFilePath = r.Env.Config
		}
	}

	if r.IsTerminal {
		return ErrNoStdin
	}

	vout := &controller.TerraformVersionOutput{}
	if err := json.NewDecoder(r.Stdin).Decode(vout); err != nil {
		return fmt.Errorf(`parse stdin as the output of "terraform version -json": %w`, err)
	}

	ctrl := controller.New(afero.NewOsFs())
	return ctrl.Run(c.Context, r.LogE, param, vout) //nolint:wrapcheck
}
