package controller

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

func (c *Controller) Run(_ context.Context, logE *logrus.Entry, param *ParamRun, vout *TerraformVersionOutput) error {
	cfg := &Config{}
	if err := c.readConfig(cfg, param); err != nil {
		return err
	}

	if len(vout.ProviderSelections) == 0 {
		logE.Warn(`provider_selections is empty. Maybe the input is wrong. Please run "terraform init" before running "terraform version -json" if you haven't run`)
		return nil
	}

	providers := make(map[string]version.Constraints, len(cfg.Providers))
	if err := parseConfig(cfg, providers); err != nil {
		return err
	}

	return validate(vout, providers)
}

func parseConfig(cfg *Config, providers map[string]version.Constraints) error {
	for _, provider := range cfg.Providers {
		if provider.Name == "" {
			return ErrProviderNameIsRequired
		}
		if provider.VersionConstraints == "" {
			providers[provider.Name] = nil
			continue
		}
		constraints, err := version.NewConstraint(provider.VersionConstraints)
		if err != nil {
			return fmt.Errorf("parse version constraints: %w", logerr.WithFields(err, logrus.Fields{
				"provider_name":                provider.Name,
				"provider_version_constraints": provider.VersionConstraints,
			}))
		}
		providers[provider.Name] = constraints
	}
	return nil
}

func validate(vout *TerraformVersionOutput, providers map[string]version.Constraints) error {
	for providerName, providerVersion := range vout.ProviderSelections {
		constraints, ok := providers[providerName]
		if !ok {
			return logerr.WithFields(ErrDisallowedProvider, logrus.Fields{ //nolint:wrapcheck
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
			return logerr.WithFields(ErrDisallowedProviderVersion, logrus.Fields{ //nolint:wrapcheck
				"provider_name":                providerName,
				"provider_version":             providerVersion,
				"provider_version_constraints": constraints.String(),
			})
		}
	}
	return nil
}
