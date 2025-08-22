package gnome

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_TimeSignature(t *testing.T) {

	Convey("When a NewTimeSignature is created, the values are correct", t, func() {
		ts := NewTimeSignature(11, 2, 240)
		So(ts.Beats.Load(), ShouldEqual, 11)
		So(ts.NoteValue.Load(), ShouldEqual, 2)
		So(ts.Tempo.Load(), ShouldEqual, 240)
		So(ts.TempoToDuration(), ShouldEqual, 250*time.Millisecond)
		Convey("And when a signature is provided, it is parsed properly and the values reflect them", func() {
			err := ts.FromString("6/8")
			So(err, ShouldBeNil)
			So(ts.Beats.Load(), ShouldEqual, 6)
			So(ts.NoteValue.Load(), ShouldEqual, 8)
			So(ts.ToString(), ShouldEqual, "6/8")
		})
		Convey("And when the Tempo is set from a Duration, it is parsed properly and the values reflect that", func() {
			ts.TempoFromDuration(500 * time.Millisecond)
			So(ts.Tempo.Load(), ShouldEqual, 120)
			So(ts.TempoToDuration(), ShouldEqual, 500*time.Millisecond)
			So(ts.String(), ShouldEndWith, "t=120")
		})
	})
	Convey("When a zero-value TimeSignature is created, the values are zero-value too", t, func() {
		ts := TimeSignature{}
		So(ts.Beats.Load(), ShouldBeZeroValue)
		So(ts.NoteValue.Load(), ShouldBeZeroValue)
		So(ts.Tempo.Load(), ShouldBeZeroValue)
		So(ts.TempoToDuration(), ShouldEqual, 0)
		Convey("And when bad signatures are provided, they fail properly and do not update the values", func() {
			err := ts.FromString("0/4") // zero not allowed for beats
			So(err, ShouldNotBeNil)
			err = ts.FromString("4/0") // zero not allowed for note value
			So(err, ShouldNotBeNil)
			err = ts.FromString("68") // missing separator
			So(err, ShouldNotBeNil)
			err = ts.FromString("4/four") // text
			So(err, ShouldNotBeNil)
			err = ts.FromString("six/four") // so much text
			So(err, ShouldNotBeNil)
			err = ts.FromString("4/6/8") // too many separators
			So(err, ShouldNotBeNil)

			So(ts.Beats.Load(), ShouldBeZeroValue)
			So(ts.NoteValue.Load(), ShouldBeZeroValue)
		})
		Convey("And when the Tempo is set wrong, it fails and does not update the value", func() {
			ts.Tempo.Store(60)
			ts.TempoFromDuration(0)
			So(ts.Tempo.Load(), ShouldEqual, 60)
		})
	})

}
