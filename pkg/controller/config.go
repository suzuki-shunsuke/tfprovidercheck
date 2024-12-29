package controller

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Providers []*Provider `json:"providers"`
}

func (c *Controller) readConfig(cfg *Config, param *ParamRun) error {
	if param.ConfigBody != "" {
		if err := yaml.Unmarshal([]byte(param.ConfigBody), cfg); err != nil {
			return fmt.Errorf("marshal configuration as YAML: %w", err)
		}
		return nil
	}
	if param.ConfigFilePath == "" {
		// Don't search a configuration file from the current directory to the root directory to prevent being replaced with a fake configuration file
		param.ConfigFilePath = ".tfprovidercheck.yaml"
	}
	return c.readConfigFile(cfg, param.ConfigFilePath)
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
