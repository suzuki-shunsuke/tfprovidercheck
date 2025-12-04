package cli_test

import (
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/suzuki-shunsuke/tfprovidercheck/pkg/cli"
)

func TestLDFlags_VersionString(t *testing.T) {
	t.Parallel()
	data := []struct {
		name    string
		exp     string
		ldflags *cli.LDFlags
	}{
		{
			name: "ldflags is nil",
			exp:  "unknown",
		},
		{
			name:    "version and date is empty",
			exp:     "unknown",
			ldflags: &cli.LDFlags{},
		},
		{
			name: "version is empty",
			exp:  "(2023-11-04T15:39:47Z)",
			ldflags: &cli.LDFlags{
				Date: "2023-11-04T15:39:47Z",
			},
		},
		{
			name: "date is empty",
			exp:  "1.0.0",
			ldflags: &cli.LDFlags{
				Version: "1.0.0",
			},
		},
		{
			name: "normal",
			exp:  "1.0.0 (2023-11-04T15:39:47Z)",
			ldflags: &cli.LDFlags{
				Version: "1.0.0",
				Date:    "2023-11-04T15:39:47Z",
			},
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			v := d.ldflags.VersionString()
			if v != d.exp {
				t.Fatalf("wanted %s, got %s", d.exp, v)
			}
		})
	}
}

func TestRunner_Run(t *testing.T) { //nolint:cyclop,tparallel,funlen
	t.Parallel()
	data := []struct {
		name   string
		runner *cli.Runner
		args   []string
		isErr  bool
		err    error
		errMsg string
	}{
		{
			name: "no stdin",
			err:  cli.ErrNoStdin,
			args: []string{
				"tfprovidercheck",
			},
			runner: &cli.Runner{
				Env:        &cli.Env{},
				IsTerminal: true,
			},
		},
		{
			name:   "invalid stdin",
			errMsg: `parse stdin as the output of "terraform version -json`,
			args: []string{
				"tfprovidercheck",
			},
			runner: &cli.Runner{
				Env:        &cli.Env{},
				IsTerminal: false,
				Stdin:      strings.NewReader("hello"),
			},
		},
		{
			name: "normal",
			args: []string{
				"tfprovidercheck",
			},
			runner: &cli.Runner{
				Env: &cli.Env{
					ConfigBody: `providers:
- name: registry.terraform.io/hashicorp/aws
`,
				},
				IsTerminal: false,
				Stdin: strings.NewReader(`{
  "terraform_version": "1.5.7",
  "platform": "darwin_arm64",
  "provider_selections": {
    "registry.terraform.io/hashicorp/aws": "3.76.1"
  },
  "terraform_outdated": true
}`),
			},
		},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	for _, d := range data { //nolint:paralleltest
		d.runner.Logger = logger
		t.Run(d.name, func(t *testing.T) {
			// app.Run isn't goroutine safe
			// t.Parallel()
			if err := d.runner.Run(t.Context(), d.args...); err != nil { //nolint:nestif
				if !d.isErr && d.err == nil && d.errMsg == "" {
					t.Fatal(err)
				}
				if d.err != nil {
					if !errors.Is(err, d.err) {
						t.Fatalf("wanted %v, got %v", d.err, err)
					}
				}
				if d.errMsg != "" {
					errMsg := err.Error()
					if !strings.Contains(errMsg, d.errMsg) {
						t.Fatalf(`error message doesn't include expected message. message=%s, expected_message=%s`, errMsg, d.errMsg)
					}
				}
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
		})
	}
}
