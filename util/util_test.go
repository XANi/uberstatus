package util

import (
	"time"
	"testing"
	"unicode/utf8"
	. "github.com/smartystreets/goconvey/convey"
)
// TODO H/min should probably be in 30:40 format
var formatList = map[string]time.Duration {
	"112h32m": time.Duration(time.Hour * 112 + time.Minute * 32),
	"11h32m": time.Duration(time.Hour * 11 + time.Minute * 32),
	"4h12m": time.Duration(time.Hour * 4 + time.Minute * 12),
	"4h01m": time.Duration(time.Hour * 4 + time.Minute),
	"4h00m": time.Duration(time.Hour * 4),
	"44.9m": time.Duration(time.Minute * 44 +  time.Second * 54),
	"39.0m": time.Duration(time.Minute * 39),
	"4.21s": time.Duration(time.Second * 4 + time.Millisecond * 210),
	"3763ms": time.Duration(time.Millisecond * 3763),
	"763.0ms": time.Duration(time.Millisecond * 763),
	"453.0ms": time.Duration(time.Millisecond * 453),
	"35.11ms": time.Duration(time.Millisecond * 35 + time.Nanosecond * 111000),
	"5.00ms": time.Duration(time.Millisecond * 5),
	"700.6µs": time.Duration(time.Microsecond * 700 + 555),
	"321.6µs": time.Duration(time.Microsecond * 321 + 555),
	"150.6µs":time.Duration(time.Microsecond * 150 + 555),
	"34.05µs": time.Duration(time.Microsecond * 34 + 50),
	"3.56µs": time.Duration(time.Microsecond * 3 + 555),
	"4ns": time.Duration(4),
}

func TestDurationFormat(t *testing.T) {
	for format, dur := range formatList {
		formattedString := FormatDuration(dur)
		Convey("time: " + format, t, func() {
			So(formattedString, ShouldEndWith, format)
			So(utf8.RuneCountInString(formattedString), ShouldBeLessThanOrEqualTo,7)
		})
	}

}

func TestDurationPaddedFormat(t *testing.T) {
	for format, dur := range formatList {
		formattedString := FormatDurationPadded(dur)
		Convey("time: " + format, t, func() {
			So(formattedString, ShouldEndWith, format)
			So(utf8.RuneCountInString(formattedString), ShouldEqual,7)
		})
}

}
