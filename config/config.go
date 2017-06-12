package config

import (
	"github.com/op/go-logging"
)


var CfgFiles =[]string{
	"$HOME/.config/uberstatus/uberstatus.conf",
	"./cfg/uberstatus.conf",
	"./cfg/uberstatus.default.conf",
	"/usr/share/doc/uberstatus/uberstatus.example.conf",
}
var log = logging.MustGetLogger("main")


type PluginConfig struct {
	Name string
	Instance string
	Plugin string
	Config map[string]interface{}
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
