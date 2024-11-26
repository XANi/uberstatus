package debug

import (
	"fmt"
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestNew(t *testing.T) {
	var cfg config.Config
	err := yaml.Unmarshal([]byte(`
---
plugins:
    - name: debug
      instance: inst1
      plugin: debug
      config:
        prefix: dbx
        interval: 1234
`), &cfg)
	if err != nil {
		t.Fatalf("bad input yaml: %s", err)
	}

	fmt.Printf("%+v", cfg)
	ucfg := uber.PluginConfig{
		Config: cfg.Plugins[0],
	}
	out, err := New(ucfg)
	assert.NoError(t, err)
	assert.NoError(t, out.Init())
	assert.Equal(t, out.GetUpdateInterval(), 1234)

}
