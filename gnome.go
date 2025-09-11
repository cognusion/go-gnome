// Package gnome is a library for building a nearly-realtime metro...gnomes. Get it? GET IT?!
// One would think there would already be such a thing, but given the complexities involved in getting
// decent timing, there was not. Supports WAV, MP3, and Ogg Vorbis as passthoughs from `gopxl/beep`
// (which stands on the shoulders of other giants).
//
// This is not perfect, either. If the system is very busy, the rhythm will not be smooth.
// If the tempo is exceptionally high, the rhythm will not be smooth.
// On a normal system, doing nothing else, BPMs under 180 are almost always great.
//
// Gnome was built from scratch expecting WASM as the target platform. Some decisions that may seem odd
// were made because of that. Most odd decisions are simply odd. The consuming app was written for my
// brother-in-law's music students: [MetroGnome](https://github.com/cognusion/metrognome).
package gnome

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"sync/atomic"
	"time"

	"github.com/cognusion/go-gnome/speaker"
	"github.com/cognusion/go-recyclable"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
	"github.com/h2non/filetype"
	uatomic "go.uber.org/atomic"
	//"github.com/gopxl/beep/v2/speaker" This stays commented out, to remind us we customized speaker
)

const (
	beatDivisions = 60000000
)

var (
	// RPool is a recyclable.BufferPool
	RPool = recyclable.NewBufferPool()
)

// Buffer is our local interface to encapsulate all the interfaces
type Buffer = io.ReadSeekCloser

// TickFunc is a function called during a tick event. The current beat number is passed to it.
// These functions must absolutely handle being called multiple times in quick succession, handling order,
// overlaps, and all of that jazz.
type TickFunc func(int)

// TickFilter is passed the beat count, and returns true if the tick sound should be played
type TickFilter func(int) bool

// Gnome is a metro...gnome. Get it? Get it?!
type Gnome struct {
	// TS tracks and reports the time signature information for the 'gnome.
	TS *TimeSignature

	player     beep.StreamSeekCloser
	interval   uatomic.Duration
	ctx        context.Context
	cancelFunc func()
	pauseChan  chan bool
	paused     atomic.Bool
	mute       atomic.Bool
	running    atomic.Bool
	tickFilter atomic.Pointer[TickFilter]
	tickFunc   TickFunc
}

// NewGnomeFromBuffer takes a Buffer, a TimeSignature and an optional tickFunc to call when
// the 'gnome fires, and gives you a Gnome or an error. :)
func NewGnomeFromBuffer(buff Buffer, ts *TimeSignature, tickFunc TickFunc) (*Gnome, error) {
	// Require a ts
	if ts == nil {
		return nil, fmt.Errorf("a valid TimeSignature is required")
	}

	// Check the buffer and open a streamer.
	streamer, format, err := BufferToStreamer(buff)
	if err != nil {
		return nil, fmt.Errorf("decoding file failed: %w", err)
	}

	// Prime the speaker
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	ctx, cancel := context.WithCancel(context.Background())
	g := &Gnome{
		player:     streamer,
		pauseChan:  make(chan bool, 1),
		TS:         ts,
		ctx:        ctx,
		tickFunc:   tickFunc,
		cancelFunc: cancel,
	}
	g.interval.Store(ts.TempoToDuration())
	tf := TickFilter(func(beat int) bool {
		return true
	})
	g.tickFilter.Store(&tf) // default filter returns true

	return g, nil
}

// NewGnomeFromFile takes a filename, a TimeSignature and an optional tickFunc to call when
// the 'gnome fires, and gives you a Gnome or an error. :)
func NewGnomeFromFile(soundFile string, ts *TimeSignature, tickFunc TickFunc) (*Gnome, error) {
	buff, err := FileToBuffer(soundFile)
	if err != nil {
		return nil, err
	}
	return NewGnomeFromBuffer(buff, ts, tickFunc)
}

// ReplaceStreamerFromBuffer attempts to replace the current steamer with a new one from the provided buffer.
func (g *Gnome) ReplaceStreamerFromBuffer(buff Buffer) error {
	if g.player == nil {
		// Should never happen unless they didn't use a New func.
		return fmt.Errorf("current streamer not initialized")
	}
	streamer, _, err := BufferToStreamer(buff)
	if err != nil {
		return err
	}
	g.player.Close() // close the old streamer

	g.player = streamer
	return nil
}

// Restart will re-initialize some stopped components so the 'gnome can carry on.
func (g *Gnome) Restart() error {
	if g.running.CompareAndSwap(false, true) {
		g.ctx, g.cancelFunc = context.WithCancel(context.Background())
		go g.ticker(g.tick, g.tickFunc)
	} else {
		// Already running
		return fmt.Errorf("already running. Stop before restarting")
	}
	return nil
}

// Start -s the Gnome. Only works the first time you call it.
func (g *Gnome) Start() {
	if g.running.CompareAndSwap(false, true) {
		go g.ticker(g.tick, g.tickFunc)
	} // else already running
}

// Mute toggles whether or not audio will play.
func (g *Gnome) Mute() {
	g.mute.Swap(!g.mute.Load()) // toggle
}

// Stop stops the Gnome.
func (g *Gnome) Stop() {
	g.cancelFunc()
	g.paused.Store(false)
	g.running.Store(false)
}

// Close terminally cleans up all the things.
func (g *Gnome) Close() {
	g.player.Close()
	speaker.Close()

}

// Pause toggles whether the 'gnome is paused.
func (g *Gnome) Pause() {
	if g.running.Load() {
		g.pauseChan <- true
		g.paused.Swap(!g.paused.Load()) // toggle
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

// IsPaused returns if the 'gnome is paused.
// This should not be used for blocking or timing decisions
func (g *Gnome) IsPaused() bool {
	return g.paused.Load()
}

// SetTickFilter installs a new tickFilter. A tickFilter is passed the beat count and returns true if
// the sound should be played for that beat
func (g *Gnome) SetTickFilter(tf TickFilter) error {
	if tf == nil {
		return fmt.Errorf("passed TickFilter is nil")
	}
	g.tickFilter.Store(&tf)
	return nil
}

// tick() is the default TickFunc what happens every time the timer fires.
func (g *Gnome) tick(beat int) {
	// Sound!
	tf := *g.tickFilter.Load()
	if tf(beat) && !g.mute.Load() {
		speaker.Clear()
		g.player.Seek(0)
		speaker.Play(g.player)
	}
}

// ticker is our internal loop for handling time. The channel-based time functions like
// time.Ticker and time.After are very very variable. Doing a tight loop with a
// controlled variable-duration time.Sleep has the least drift I have found. Wow, does Go
// need a high-precision timer package.
func (g *Gnome) ticker(ticks ...TickFunc) {
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
				}
			default:
				// Sleep to Next Target
				t = t.Add(g.interval.Load())
				// Compute the desired sleep time to reach the target
				// Sleep
				time.Sleep(time.Until(t))

				// iterate over the ticks.
				// This is more reliably performant than having
				// g.tick call tickFunc, oddly enough.
				for _, tick := range ticks {
					if tick != nil {
						tick(c)
					}
				}

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
func BufferToStreamer(buff Buffer) (beep.StreamSeekCloser, beep.Format, error) {
	var (
		err      error
		streamer beep.StreamSeekCloser
		format   beep.Format
	)

	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	_, err = buff.Read(head)
	buff.Seek(0, 0) // reset
	if err != nil {
		return streamer, format, err
	}

	kind, err := filetype.Match(head)
	if err != nil {
		return streamer, format, err
	}
	switch kind.MIME.Value {
	case "audio/x-wav":
		// WAV
		streamer, format, err = wav.Decode(buff)
	case "audio/mpeg":
		// MP3
		streamer, format, err = mp3.Decode(buff)
	case "audio/ogg":
		// Ogg Vorbis
		streamer, format, err = vorbis.Decode(buff)
	default:
		err = fmt.Errorf("buffer does not contain a supported format: %s (MIME: %s)", kind.Extension, kind.MIME.Value)
	}

	return streamer, format, err
}
