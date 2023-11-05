package controller_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/tfprovidercheck/pkg/controller"
)

const (
	dirPermission  os.FileMode = 0o775
	filePermission os.FileMode = 0o644
)

func newFs(files map[string]string, dirs ...string) (afero.Fs, error) {
	fs := afero.NewMemMapFs()
	for name, body := range files {
		if err := afero.WriteFile(fs, name, []byte(body), filePermission); err != nil {
			return nil, err //nolint:wrapcheck
		}
	}
	for _, dir := range dirs {
		if err := fs.MkdirAll(dir, dirPermission); err != nil {
			return nil, err //nolint:wrapcheck
		}
	}
	return fs, nil
}

func TestController_Run(t *testing.T) { //nolint:gocognit,cyclop,funlen
	t.Parallel()
	data := []struct {
		name   string
		files  map[string]string
		dirs   []string
		param  *controller.ParamRun
		vout   *controller.TerraformVersionOutput
		isErr  bool
		err    error
		errMsg string
	}{
		{
			name: "invalid config body",
			param: &controller.ParamRun{
				ConfigBody: "{",
			},
			errMsg: `marshal configuration as YAML: `,
		},
		{
			name:   "config isn't found",
			param:  &controller.ParamRun{},
			errMsg: `configuration file .tfprovidercheck.yaml isn't found`,
		},
		{
			name: "config path",
			param: &controller.ParamRun{
				ConfigFilePath: "/workspace/foo.yaml",
				PWD:            "/workspace",
			},
			files: map[string]string{
				"/workspace/foo.yaml": `providers:
- name: registry.terraform.io/hashicorp/aws
`,
			},
			vout: &controller.TerraformVersionOutput{
				ProviderSelections: map[string]string{
					"registry.terraform.io/hashicorp/aws": "3.0.0",
				},
			},
		},
		{
			name: "invalid config file",
			param: &controller.ParamRun{
				ConfigFilePath: "/workspace/foo.yaml",
				PWD:            "/workspace",
			},
			files: map[string]string{
				"/workspace/foo.yaml": `{`,
			},
			vout: &controller.TerraformVersionOutput{
				ProviderSelections: map[string]string{
					"registry.terraform.io/hashicorp/aws": "3.0.0",
				},
			},
			errMsg: `read a configuration file as YAML: `,
		},
		{
			name: "invalid config file path",
			param: &controller.ParamRun{
				ConfigFilePath: "/workspace/foo.yaml",
				PWD:            "/workspace",
			},
			vout: &controller.TerraformVersionOutput{
				ProviderSelections: map[string]string{
					"registry.terraform.io/hashicorp/aws": "3.0.0",
				},
			},
			errMsg: `open a configuration file: `,
		},
		{
			name: "config search",
			param: &controller.ParamRun{
				PWD: "/workspace/foo",
			},
			files: map[string]string{
				"/workspace/.tfprovidercheck.yaml": `providers:
- name: registry.terraform.io/hashicorp/aws
`,
			},
			vout: &controller.TerraformVersionOutput{
				ProviderSelections: map[string]string{
					"registry.terraform.io/hashicorp/aws": "3.0.0",
				},
			},
		},
		{
			name: "config body",
			param: &controller.ParamRun{
				ConfigBody: `providers:
- name: registry.terraform.io/hashicorp/aws
`,
			},
			vout: &controller.TerraformVersionOutput{
				ProviderSelections: map[string]string{
					"registry.terraform.io/hashicorp/aws": "3.0.0",
				},
			},
		},
		{
			name: "provider name is required",
			param: &controller.ParamRun{
				ConfigBody: `providers:
- version: ">= 3.0.0"
`,
			},
			vout: &controller.TerraformVersionOutput{
				ProviderSelections: map[string]string{
					"registry.terraform.io/hashicorp/aws": "3.0.0",
				},
			},
			err: controller.ErrProviderNameIsRequired,
		},
		{
			name: "invalid version constraint",
			param: &controller.ParamRun{
				ConfigBody: `providers:
- name: registry.terraform.io/hashicorp/aws
  version: "=> 3.0.0"
`,
			},
			vout: &controller.TerraformVersionOutput{
				ProviderSelections: map[string]string{
					"registry.terraform.io/hashicorp/aws": "3.0.0",
				},
			},
			errMsg: "parse version constraints",
		},
		{
			name: "disallowed provider",
			param: &controller.ParamRun{
				ConfigBody: `providers:
- name: foo
`,
			},
			vout: &controller.TerraformVersionOutput{
				ProviderSelections: map[string]string{
					"registry.terraform.io/hashicorp/aws": "3.0.0",
				},
			},
			err: controller.ErrDisallowedProvider,
		},
		{
			name: "disallowed provider version",
			param: &controller.ParamRun{
				ConfigBody: `providers:
- name: registry.terraform.io/hashicorp/aws
  version: ">= 4.0.0"
`,
			},
			vout: &controller.TerraformVersionOutput{
				ProviderSelections: map[string]string{
					"registry.terraform.io/hashicorp/aws": "3.0.0",
				},
			},
			err: controller.ErrDisallowedProviderVersion,
		},
		{
			name: "invalid version",
			param: &controller.ParamRun{
				ConfigBody: `providers:
- name: registry.terraform.io/hashicorp/aws
`,
			},
			vout: &controller.TerraformVersionOutput{
				ProviderSelections: map[string]string{
					"registry.terraform.io/hashicorp/aws": "...",
				},
			},
			errMsg: `parse the provider version as semver: `,
		},
	}
	for _, d := range data {
		d := d
		logE := logrus.NewEntry(logrus.New())
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			fs, err := newFs(d.files, d.dirs...)
			if err != nil {
				t.Fatal(err)
			}
			ctrl := controller.New(fs)
			if err := ctrl.Run(context.Background(), logE, d.param, d.vout); err != nil { //nolint:nestif
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
