// Package gnome is a library for building a nearly-realtime metro...gnomes. Get it? GET IT?!
// One would think there would already be such a thing, but given the complexities involved in getting
// decent timing, there was not.
//
// This is not perfect, either. If the system is very busy, the rhythm
// will not be smooth. If the tempo is exceptionally high, the rhythm will not be smooth. On a normal
// system, doing nothing else, BPMs under 180 are almost always great.
//
// Gnome was built from scratch expecting WASM as the target platform. Some decisions that may seem odd
// were made because of that. Most odd decisions are simply. The consuming app was written for my
// brother-in-law's music students: [MetroGnome](https://github.com/cognusion/metrognome)
package gnome

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cognusion/go-recyclable"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
	uatomic "go.uber.org/atomic"
)

const (
	beatDivisions = 60000000
)

var (
	// RPool is a recyclable.BufferPool
	RPool = recyclable.NewBufferPool()
)

// Gnome is a metro...gnome. Get it? Get it?!
type Gnome struct {
	// TS tracks and reports the time signature information for the 'gnome.
	TS *TimeSignature

	player     beep.StreamSeekCloser
	interval   uatomic.Duration
	ctx        context.Context
	cancelFunc func()
	pauseChan  chan bool
	mute       atomic.Bool
	running    atomic.Bool
	tickFunc   func(int)
}

// NewGnomeBufferTick takes an io.ReadSeekCloser, tempo, and a tickFunc to call when the 'gnome fires.
func NewGnomeBufferTick(buff io.ReadSeekCloser, tempo int32, tickFunc func(int)) (*Gnome, error) {
	g, err := NewGnomeBuffer(buff, tempo)
	if err != nil {
		return nil, err
	}
	// Add the tickFunc
	g.tickFunc = tickFunc

	return g, nil
}

// NewGnomeBuffer takes an io.ReadSeekCloser and a tempo.
func NewGnomeBuffer(buff io.ReadSeekCloser, tempo int32) (*Gnome, error) {
	// Check the buffer and open a streamer.
	streamer, format, err := BufferToStreamer(buff)
	if err != nil {
		return nil, fmt.Errorf("decoding file failed: %w", err)
	}

	// Prime the speaker
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// Set the TS
	ts := &TimeSignature{}
	ts.FromString("4/4")
	ts.Tempo.Store(int32(tempo))

	ctx, cancel := context.WithCancel(context.Background())
	g := &Gnome{
		player:     streamer,
		pauseChan:  make(chan bool, 1),
		TS:         ts,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	g.interval.Store(ts.TempoToDuration())

	return g, nil
}

// NewGnomeWithTickFunc takes a file string, tempo, and a tickFunc to call when the 'gnome fires.
func NewGnomeWithTickFunc(soundFile string, tempo int32, tickFunc func(int)) (*Gnome, error) {
	g, e := NewGnome(soundFile, tempo)
	if e != nil {
		return nil, e
	}
	g.tickFunc = tickFunc

	return g, nil
}

// NewGnome loads a WAV or MP3, and plays it every interval, returning an error if there is a problem
// loading the file.
func NewGnome(soundFile string, tempo int32) (*Gnome, error) {
	buff, err := FileToBuffer(soundFile)
	if err != nil {
		return nil, err
	}
	return NewGnomeBuffer(buff, tempo)
}

// Restart will re-initialize some stopped components so the 'gnome can carry on.
func (g *Gnome) Restart() error {
	if g.running.CompareAndSwap(false, true) {
		g.ctx, g.cancelFunc = context.WithCancel(context.Background())
		go g.ticker()
	} else {
		// Already running
		return fmt.Errorf("already running. Stop before restarting")
	}
	return nil
}

// Start -s the Gnome. Only works the first time you call it.
func (g *Gnome) Start() {
	if g.running.CompareAndSwap(false, true) {
		go g.ticker()
	} // else already running
}

// Mute toggles whether or not audio will play.
func (g *Gnome) Mute() {
	g.mute.Swap(!g.mute.Load()) // toggle
}

// Stop stops the Gnome.
func (g *Gnome) Stop() {
	g.cancelFunc()
	g.running.Store(false)
}

// Close terminally cleans up all the things.
func (g *Gnome) Close() {
	g.player.Close()

}

// Pause toggles whether the 'gnome is paused.
func (g *Gnome) Pause() {
	if g.running.Load() {
		g.pauseChan <- true
	}
	// else we would wedge
}

// Change sets a new tempo.
func (g *Gnome) Change(tempo int32) {
	if tempo > 0 {
		g.interval.Store(FromBPM(tempo))
	}
}

// IsRunning returns if the 'gnome is running.
func (g *Gnome) IsRunning() bool {
	return g.running.Load()
}

// tick() is what happens every time the timer fires, in a separate goro.
// This must absolutely be able to be called multiple times, handling order,
// overlaps, and all of that jazz.
func (g *Gnome) tick(beat int) {
	// Sound!
	if !g.mute.Load() {
		g.player.Seek(0)
		speaker.Play(g.player)
	}

	if g.tickFunc != nil {
		g.tickFunc(beat)
	}
}

// ticker is our internal loop for handling time. The channel-based time functions like
// time.Ticker and time.After are very very variable. Doing a tight loop with a
// controlled variable-duration time.Sleep has the least drift I have found. Wow, does Go
// need a high-precision timer package.
func (g *Gnome) ticker() {
	var (
		t = time.Now()
		c int
	)

	go func() {
		c = 1 // we start counting at 1 in music, I have been told.
		for {
			select {
			case <-g.ctx.Done():
				// Bail!
				return
			case <-g.pauseChan:
				// Pause!
				select {
				case <-g.ctx.Done():
					// Bail!
					return
				case <-g.pauseChan:
					// Resume!
					t = time.Now()
					c = 1
				}
			default:
				// Sleep to Next Target
				t = t.Add(g.interval.Load())
				// Compute the desired sleep time to reach the target
				// Sleep
				time.Sleep(time.Until(t))

				go g.tick(c)

				c++
				if c > int(g.TS.Beats.Load()) {
					c = 1
				}
			}
		}
	}()
}

// FromBPM converts a beats-per-minute tempo to a Microsecond-precise time.Duration.
func FromBPM(bpm int32) time.Duration {
	if bpm <= 0 {
		// Safety
		return 0
	}
	microsPerBeat := time.Duration(beatDivisions / bpm)
	return microsPerBeat * time.Microsecond
}

// ToBPM converts a Microsecond-precise time.Duration to a beats-per-minute tempo.
func ToBPM(interval time.Duration) int32 {
	if interval == 0 {
		return 0
	}
	//#nosec G115 -- this cast is edgily unsafe, but we are dividing TF out of it.
	return int32(60 / interval.Seconds())
}

// FileToBuffer opens and reads the filename into a Buffer, returning it or an error.
func FileToBuffer(filename string) (*recyclable.Buffer, error) {
	f, err := os.Open(path.Clean(filename))
	if err != nil {
		return nil, fmt.Errorf("reading file failed: %w", err)
	}

	buff := RPool.Get()
	buff.ResetFromReader(f)
	f.Close()
	return buff, nil
}

// BufferToStreamer checks the first 12 bytes of the Buffer to see if it's a WAV,
// and tries to decode it as that or an MP3. Errors are returned if anything fails.
func BufferToStreamer(buff io.ReadSeekCloser) (beep.StreamSeekCloser, beep.Format, error) {
	var (
		b        = make([]byte, 12)
		n        int
		err      error
		streamer beep.StreamSeekCloser
		format   beep.Format
	)

	// determine type of byte data
	n, err = buff.Read(b)
	buff.Seek(0, 0) //reset the seek pointer

	if err != nil || n != 12 {
		err = fmt.Errorf("reading first 12 (%d) of the buffer failed: %w", n, err)
	} else if strings.Contains(string(b), "RIFF") && strings.Contains(string(b), "WAVE") {
		// WAV
		streamer, format, err = wav.Decode(buff)
	} else {
		// MP3 we hope
		streamer, format, err = mp3.Decode(buff)
	}

	return streamer, format, err
}
