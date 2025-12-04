package controller

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hashicorp/go-version"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

func (c *Controller) Run(_ context.Context, _ *slog.Logger, param *ParamRun, vout *TerraformVersionOutput) error {
	cfg := &Config{}
	if err := c.readConfig(cfg, param); err != nil {
		return err
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
			return fmt.Errorf("parse version constraints: %w", slogerr.With(err,
				"provider_name", provider.Name,
				"provider_version_constraints", provider.VersionConstraints,
			))
		}
		providers[provider.Name] = constraints
	}
	return nil
}

func validate(vout *TerraformVersionOutput, providers map[string]version.Constraints) error {
	for providerName, providerVersion := range vout.ProviderSelections {
		constraints, ok := providers[providerName]
		if !ok {
			return slogerr.With(ErrDisallowedProvider, "provider_name", providerName) //nolint:wrapcheck
		}
		v, err := version.NewVersion(providerVersion)
		if err != nil {
			return fmt.Errorf("parse the provider version as semver: %w", slogerr.With(err,
				"provider_name", providerName,
				"provider_version", providerVersion,
			))
		}
		if !constraints.Check(v) {
			return slogerr.With(ErrDisallowedProviderVersion, //nolint:wrapcheck
				"provider_name", providerName,
				"provider_version", providerVersion,
				"provider_version_constraints", constraints.String(),
			)
		}
	}
	return nil
}
