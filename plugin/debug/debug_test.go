package debug

import (
	"fmt"
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	. "github.com/smartystreets/goconvey/convey"
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

	Convey("create", t, func() {
		So(err, ShouldBeNil)
		So(out.Init(), ShouldBeNil)
		So(out.GetUpdateInterval(), ShouldEqual, 1234)
	})
}
