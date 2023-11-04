package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/go-findconfig/findconfig"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"gopkg.in/yaml.v3"
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

type Config struct {
	Providers []*Provider
}

type TerraformVersionOutput struct {
	// https://github.com/hashicorp/terraform/blob/05f877166dec78b571ca36ca4922ece8b83fd0f8/internal/command/version.go#L28-L33
	ProviderSelections map[string]string `json:"provider_selections"`
}

type Provider struct {
	Name               string
	VersionConstraints string `yaml:"version"`
}

func (c *Controller) Run(_ context.Context, _ *logrus.Entry, param *ParamRun, vout *TerraformVersionOutput) error {
	cfg := &Config{}
	if err := c.readConfig(cfg, param); err != nil {
		return err
	}

	providers := make(map[string]version.Constraints, len(cfg.Providers))
	for _, provider := range cfg.Providers {
		if provider.Name == "" {
			return errors.New("providers[].name is required")
		}
		if provider.VersionConstraints == "" {
			providers[provider.Name] = nil
			continue
		}
		constraints, err := version.NewConstraint(provider.VersionConstraints)
		if err != nil {
			return fmt.Errorf("parse version constraints: %w", err)
		}
		providers[provider.Name] = constraints
	}

	for providerName, providerVersion := range vout.ProviderSelections {
		constraints, ok := providers[providerName]
		if !ok {
			return logerr.WithFields(errors.New("this Terraform Provider is disallowed"), logrus.Fields{ //nolint:wrapcheck
				"provider_name": providerName,
			})
		}
		v, err := version.NewVersion(providerVersion)
		if err != nil {
			return fmt.Errorf("parse the provider version as semver: %w", logerr.WithFields(err, logrus.Fields{
				"provider_name":    providerName,
				"provider_version": providerVersion,
			}))
		}
		if !constraints.Check(v) {
			return logerr.WithFields(errors.New("this Terraform Provider version is disallowed"), logrus.Fields{ //nolint:wrapcheck
				"provider_name":                providerName,
				"provider_version":             providerVersion,
				"provider_version_constraints": constraints.String(),
			})
		}
	}
	return nil
}

func (c *Controller) readConfig(cfg *Config, param *ParamRun) error {
	if param.ConfigBody != "" {
		if err := yaml.Unmarshal([]byte(param.ConfigBody), cfg); err != nil {
			return fmt.Errorf("marshal configuration as YAML: %w", err)
		}
		return nil
	}
	if param.ConfigFilePath != "" {
		return c.readConfigFile(cfg, param.ConfigFilePath)
	}
	cfgFilePath := findconfig.Find(param.PWD, func(p string) bool {
		f, err := afero.Exists(c.fs, p)
		if err != nil {
			return false
		}
		return f
	}, ".tfprovidercheck.yaml")
	if cfgFilePath == "" {
		return errors.New("configuration file .tfprovidercheck.yaml isn't found")
	}
	return c.readConfigFile(cfg, cfgFilePath)
}

func (c *Controller) readConfigFile(cfg *Config, cfgFilePath string) error {
	f, err := c.fs.Open(cfgFilePath)
	if err != nil {
		return fmt.Errorf("open a configuration file: %w", err)
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return fmt.Errorf("read a configuration file as YAML: %w", err)
	}
	return nil
}
