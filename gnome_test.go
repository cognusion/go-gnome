package gnome

import (
	"crypto/rand"
	_ "embed"
	"testing"

	"github.com/fortytw2/leaktest"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_BufferToStreamer(t *testing.T) {
	defer leaktest.Check(t)()

	Convey("When FileToBuffer is called on a known WAV file", t, func() {
		b, err := FileToBuffer("metronome1.wav")
		So(b, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("and that Buffer of known-good WAV data is sent to BufferToStreamer, everything works as expected.", func() {
			s, f, e := BufferToStreamer(b)
			So(s, ShouldNotBeNil)
			So(s, ShouldNotBeZeroValue)
			So(f, ShouldNotBeZeroValue)
			So(e, ShouldBeNil)
		})
	})

	Convey("When BufferToStreamer is called on a known-bad []byte data, everything fails as expected.", t, func() {

		b := RPool.Get()
		defer b.Close()

		rb := make([]byte, 32)
		// Read random bytes into the slice
		rand.Read(rb)
		b.Reset(rb)

		s, f, e := BufferToStreamer(b)
		So(s, ShouldBeNil)
		So(f, ShouldBeZeroValue)
		So(e, ShouldNotBeNil)
	})

	Convey("When BufferToStreamer is called on a nil []byte data, everything fails as expected.", t, func() {

		b := RPool.Get()
		defer b.Close()
		b.Reset(make([]byte, 0))

		s, f, e := BufferToStreamer(b)
		So(s, ShouldBeNil)
		So(f, ShouldBeZeroValue)
		So(e, ShouldNotBeNil)
	})
}

func Test_BPMS(t *testing.T) {
	Convey("When sequences of known BPMs are converted to Durations and back, they all line up.", t, func() {
		for i := 0; i < 600; i += 10 {
			So(ToBPM(FromBPM(int32(i))), ShouldEqual, i)
		}
	})
}
