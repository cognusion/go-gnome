package gnome

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/spf13/cast"
)

// TimeSignature is a mechaism to make tracking time signatures,
// and their changes, easier and safe.
type TimeSignature struct {
	Beats     atomic.Int32
	NoteValue atomic.Int32
	Tempo     atomic.Int32
}

// NewTimeSignature returns an initialized TimeSignature.
func NewTimeSignature(beats, noteValue, tempo int32) *TimeSignature {
	ts := TimeSignature{}
	ts.Beats.Store(beats)
	ts.NoteValue.Store(noteValue)
	ts.Tempo.Store(tempo)
	return &ts
}

// TempoToDuration returns the tempo as a time.Duration.
func (ts *TimeSignature) TempoToDuration() time.Duration {
	if ts.Tempo.Load() <= 0 {
		// Safety
		return 0
	}
	microsPerBeat := time.Duration(beatDivisions / ts.Tempo.Load())
	return microsPerBeat * time.Microsecond
}

// TempoFromDuration sets the Tempo based on the interval provided.
func (ts *TimeSignature) TempoFromDuration(interval time.Duration) {
	if interval == 0 {
		return
	}
	ts.Tempo.Store(cast.ToInt32(60 / interval.Seconds()))
}

// ToString returns the signature string (e.g. "4/4" or "6/8"), without the Tempo.
func (ts *TimeSignature) ToString() string {
	return fmt.Sprintf("%d/%d", ts.Beats.Load(), ts.NoteValue.Load())
}

// FromString takes a signature string (e.g. "4/4" or "6/8") and sets TimeSignature accordingly.
func (ts *TimeSignature) FromString(sig string) error {
	parts := strings.Split(sig, "/")
	if len(parts) != 2 {
		return fmt.Errorf("signature is unparsable")
	}
	beats, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("signature is unparsable: %w", err)
	} else if beats == 0 {
		return fmt.Errorf("beats cannot be zero")
	}
	nv, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("signature is unparsable: %w", err)
	} else if nv == 0 {
		return fmt.Errorf("note value cannot be zero")
	}
	// Post: everything parsed ok
	ts.Beats.Store(cast.ToInt32(beats))
	ts.NoteValue.Store(cast.ToInt32(nv))

	return nil
}

// String pretty-prints the values of the TimeSignature.
func (ts *TimeSignature) String() string {
	return fmt.Sprintf("%d/%d t=%d", ts.Beats.Load(), ts.NoteValue.Load(), ts.Tempo.Load())
}
