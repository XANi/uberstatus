package config

import (
	"github.com/op/go-logging"
	"gopkg.in/yaml.v3"
)


var CfgFiles =[]string{
	"$HOME/.config/uberstatus/uberstatus.conf",
	"./cfg/uberstatus.conf",
	"./cfg/uberstatus.default.conf",
	"/usr/share/doc/uberstatus/uberstatus.example.conf",
}
var log = logging.MustGetLogger("main")


type PluginConfig struct {
	Name string `yaml:"name"`
	Instance string `yaml:"instance"`
	Plugin string `yaml:"plugin"`
	Config yaml.Node `yaml:"config"`
}
// pass your config struct to this function, it will fill it
func (p *PluginConfig) GetConfig(i interface{}) error{
	return p.Config.Decode(i)

}


type Config struct {
	Plugins []PluginConfig
}
func (c *Config) SetConfig(s string) {
	log.Infof("Loaded config file from %s",s)
}
func (c *Config) GetDefaultConfig() string {
	return exampleConfig
}
