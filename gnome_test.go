package gnome

import (
	"crypto/rand"
	_ "embed"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewGnomeBufferTick(t *testing.T) {
	// beep leaks a goro from ebitengine: https://github.com/gopxl/beep/issues/107
	//defer leaktest.Check(t)()

	Convey("When...", t, func(c C) {
		var i int64

		tf := func(tick int) {
			i++
			c.So(tick, ShouldNotBeZeroValue)
		}
		buff, err := FileToBuffer("metronome1.wav")
		So(buff, ShouldNotBeNil)
		So(err, ShouldBeNil)

		g, e := NewGnomeBufferTick(buff, 240, tf)
		So(e, ShouldBeNil)
		So(g, ShouldNotBeNil)
		defer g.Close()

		g.Mute() // let's not metronome during a test
		g.Start()
		<-time.After(time.Second)
		So(g.IsRunning(), ShouldBeTrue)
		So(g.Restart(), ShouldBeError)

		g.Pause() // Pause
		oldi := i // cache i
		<-time.After(time.Second)
		So(i, ShouldBeBetweenOrEqual, oldi, oldi+1) // Pause means Pause, with possible slip of 1
		g.Pause()                                   // resume
		g.Stop()

		So(g.IsRunning(), ShouldBeFalse)
		So(i, ShouldBeBetweenOrEqual, 2, 5)

	})
}

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
