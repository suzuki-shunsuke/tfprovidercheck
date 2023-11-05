package controller

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/go-findconfig/findconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Providers []*Provider
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
		return ErrConfigNotFound
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
