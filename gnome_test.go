package gnome

import (
	"crypto/rand"
	_ "embed"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewGnomeBufferTick(t *testing.T) {
	// beep leaks a goro from ebitengine: https://github.com/gopxl/beep/issues/107
	//defer leaktest.Check(t)()

	Convey("When a Buffer from FileToBuffer on a known-good wav file is passed to NewGnomeBufferTick, everything works as expected.", t, func(c C) {
		var i atomic.Int64

		tf := func(tick int) {
			i.Add(1)
			c.So(tick, ShouldNotBeZeroValue)
		}

		g, e := NewGnomeFromFile("testfiles/metronome1.wav", NewTimeSignature(4, 4, 240), tf)
		So(e, ShouldBeNil)
		So(g, ShouldNotBeNil)
		defer g.Close()

		g.Mute() // let's not metronome during a test
		g.Start()
		<-time.After(time.Second)
		So(g.IsRunning(), ShouldBeTrue)
		So(g.Restart(), ShouldBeError)

		g.Pause()        // Pause
		oldi := i.Load() // cache i
		<-time.After(time.Second)
		SoMsg("Too many post-Pause ticks!", i.Load(), ShouldBeBetweenOrEqual, oldi, oldi+1) // Pause means Pause, with possible slip of 1
		g.Pause()                                                                           // resume
		g.Stop()

		So(g.IsRunning(), ShouldBeFalse)
		SoMsg("Too many total ticks, given pauses and tempo!", i.Load(), ShouldBeBetweenOrEqual, 2, 5)

	})
}

func Test_BufferToStreamer(t *testing.T) {
	defer leaktest.Check(t)()

	Convey("When FileToBuffer is called on a known WAV file", t, func() {
		b, err := FileToBuffer("testfiles/metronome1.wav")
		So(err, ShouldBeNil)
		So(b, ShouldNotBeNil)

		Convey("and that Buffer of known-good WAV data is sent to BufferToStreamer, everything works as expected.", func() {
			s, f, e := BufferToStreamer(b)
			So(e, ShouldBeNil)
			So(s, ShouldNotBeNil)
			So(s, ShouldNotBeZeroValue)
			So(f, ShouldNotBeZeroValue)
		})
	})

	Convey("When FileToBuffer is called on a known MP3 file", t, func() {
		b, err := FileToBuffer("testfiles/metronome1.mp3")
		So(err, ShouldBeNil)
		So(b, ShouldNotBeNil)

		Convey("and that Buffer of known-good MP3 data is sent to BufferToStreamer, everything works as expected.", func() {
			s, f, e := BufferToStreamer(b)
			So(e, ShouldBeNil)
			So(s, ShouldNotBeNil)
			So(s, ShouldNotBeZeroValue)
			So(f, ShouldNotBeZeroValue)
		})
	})

	Convey("When FileToBuffer is called on a known OGG file", t, func() {
		b, err := FileToBuffer("testfiles/metronome1.ogg")
		So(err, ShouldBeNil)
		So(b, ShouldNotBeNil)

		Convey("and that Buffer of known-good OGG data is sent to BufferToStreamer, everything works as expected.", func() {
			s, f, e := BufferToStreamer(b)
			So(e, ShouldBeNil)
			So(s, ShouldNotBeNil)
			So(s, ShouldNotBeZeroValue)
			So(f, ShouldNotBeZeroValue)
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
