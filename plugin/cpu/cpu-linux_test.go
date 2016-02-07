package cpu

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	//	"fmt"
)

func TestCpuLinuxTicks(t *testing.T) {
	c, err := GetCpuTicks()
	Convey("Total", t, func() {
		So(err, ShouldBeNil)
		So(c[0].total, ShouldBeGreaterThan, 0)
		So(c[0].user, ShouldBeGreaterThan, 0)
		So(c[0].system, ShouldBeGreaterThan, 0)
	})
	Convey("First CPU", t, func() {
		So(err, ShouldBeNil)
		So(c[1].total, ShouldBeGreaterThan, 0)
		So(c[1].user, ShouldBeGreaterThan, 0)
	 	So(c[1].system, ShouldBeGreaterThan, 0)
	})
}
