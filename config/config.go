package config

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var CfgFiles = []string{
	"$HOME/.config/uberstatus/uberstatus.conf",
	"./cfg/uberstatus.conf",
	"./cfg/uberstatus.default.conf",
	"/usr/share/doc/uberstatus/uberstatus.example.conf",
}

type PluginConfig struct {
	Name     string             `yaml:"name"`
	Instance string             `yaml:"instance"`
	Plugin   string             `yaml:"plugin"`
	Config   yaml.Node          `yaml:"config"`
	Logger   *zap.SugaredLogger `yaml:"-"`
}

// pass your config struct to this function, it will fill it
func (p *PluginConfig) GetConfig(i interface{}) error {
	if p.Config.Kind != 0 {
		return p.Config.Decode(i)
	} else {
		return nil
	}
}

type Config struct {
	Plugins          []PluginConfig
	PanicOnBadPlugin bool               `yaml:"panic_on_bad_plugin"`
	Logger           *zap.SugaredLogger `yaml:"-"`
}

func (c *Config) SetConfig(s string) {
	c.Logger.Infof("Loaded config file from %s", s)
}
func (c *Config) GetDefaultConfig() string {
	return exampleConfig
}
