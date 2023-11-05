package controller

import (
	"github.com/spf13/afero"
)

type Controller struct {
	fs afero.Fs
}

func New(fs afero.Fs) *Controller {
	return &Controller{
		fs: fs,
	}
}

type ParamRun struct {
	ConfigFilePath string
	ConfigBody     string
	PWD            string
}

type TerraformVersionOutput struct {
	// https://github.com/hashicorp/terraform/blob/05f877166dec78b571ca36ca4922ece8b83fd0f8/internal/command/version.go#L28-L33
	ProviderSelections map[string]string `json:"provider_selections"`
}

type Provider struct {
	Name               string
	VersionConstraints string `yaml:"version"`
}
