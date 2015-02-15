package config

import (
	"fmt"
	"github.com/op/go-logging"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"os"
	"path/filepath"
)

var log = logging.MustGetLogger("main")

var cfgFiles =[]string{
	"$HOME/.config/uberstatus/uberstatus.conf",
	"./cfg/uberstatus.conf",
	"./cfg/uberstatus.default.conf",
	"/usr/share/doc/uberstatus/uberstatus.example.conf",
}



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
	var cfgFile string
	for _,element := range cfgFiles {
		filename, _ := filepath.Abs(os.ExpandEnv(element))
		log.Warning(filename)
		if _, err := os.Stat(filename); err == nil {
			cfgFile = filename
			break
		}
	}
	if cfgFile == "" {
		log.Panic("could not find config file: %v", cfgFiles)
	}
	log.Info(fmt.Sprintf("Loading config file: %s", cfgFile))
	raw_cfg, err := ioutil.ReadFile(cfgFile)
	err = yaml.Unmarshal([]byte(raw_cfg), &cfg)
	_ = err
	str, _ := yaml.Marshal(cfg)
	log.Warning(string(str))
	return cfg
}
