package config

import (
	"fmt"
	"github.com/op/go-logging"
	"gopkg.in/yaml.v1"
	"io/ioutil"
)

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


func LoadConfig() Config {
	var cfg Config
	cfgFile := "/home/xani/src/my/uberstatus/cfg/uberstatus.default.conf"
	log.Info(fmt.Sprintf("Loading config file: %s", cfgFile))
	raw_cfg, err := ioutil.ReadFile(cfgFile)
	err = yaml.Unmarshal([]byte(raw_cfg), &cfg)
	_ = err
	str, _ := yaml.Marshal(cfg)
	log.Warning(string(str))
	return cfg
}
