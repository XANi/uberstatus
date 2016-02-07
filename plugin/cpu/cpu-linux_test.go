package cpu

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
//	"fmt"
)

func TestCpuLinuxTicks(t *testing.T) {
    c, err := GetCpuTicks()
	Convey("total",t,func() {
		So(err,ShouldBeNil)
		So(c.total,ShouldBeGreaterThan,0)
		So(c.user,ShouldBeGreaterThan,0)
		So(c.system,ShouldBeGreaterThan,0)
	})
}
