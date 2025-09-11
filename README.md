

# gnome
`import "github.com/cognusion/go-gnome"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>
Package gnome is a library for building a nearly-realtime metro...gnomes. Get it? GET IT?!
One would think there would already be such a thing, but given the complexities involved in getting
decent timing, there was not. Supports WAV, MP3, and Ogg Vorbis as passthoughs from `gopxl/beep`
(which stands on the shoulders of other giants).

This is not perfect, either. If the system is very busy, the rhythm will not be smooth.
If the tempo is exceptionally high, the rhythm will not be smooth.
On a normal system, doing nothing else, BPMs under 180 are almost always great.

Gnome was built from scratch expecting WASM as the target platform. Some decisions that may seem odd
were made because of that. Most odd decisions are simply odd. The consuming app was written for my
brother-in-law's music students: [MetroGnome](<a href="https://github.com/cognusion/metrognome">https://github.com/cognusion/metrognome</a>).




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func BufferToStreamer(buff Buffer) (beep.StreamSeekCloser, beep.Format, error)](#BufferToStreamer)
* [func FileToBuffer(filename string) (*recyclable.Buffer, error)](#FileToBuffer)
* [func FromBPM(bpm int32) time.Duration](#FromBPM)
* [func ToBPM(interval time.Duration) int32](#ToBPM)
* [type Buffer](#Buffer)
* [type Gnome](#Gnome)
  * [func NewGnomeFromBuffer(buff Buffer, ts *TimeSignature, tickFunc TickFunc) (*Gnome, error)](#NewGnomeFromBuffer)
  * [func NewGnomeFromFile(soundFile string, ts *TimeSignature, tickFunc TickFunc) (*Gnome, error)](#NewGnomeFromFile)
  * [func (g *Gnome) Change(tempo int32)](#Gnome.Change)
  * [func (g *Gnome) Close()](#Gnome.Close)
  * [func (g *Gnome) IsPaused() bool](#Gnome.IsPaused)
  * [func (g *Gnome) IsRunning() bool](#Gnome.IsRunning)
  * [func (g *Gnome) Mute()](#Gnome.Mute)
  * [func (g *Gnome) Pause()](#Gnome.Pause)
  * [func (g *Gnome) ReplaceStreamerFromBuffer(buff Buffer) error](#Gnome.ReplaceStreamerFromBuffer)
  * [func (g *Gnome) Restart() error](#Gnome.Restart)
  * [func (g *Gnome) SetTickFilter(tf TickFilter) error](#Gnome.SetTickFilter)
  * [func (g *Gnome) Start()](#Gnome.Start)
  * [func (g *Gnome) Stop()](#Gnome.Stop)
* [type TickFilter](#TickFilter)
* [type TickFunc](#TickFunc)
* [type TimeSignature](#TimeSignature)
  * [func NewTimeSignature(beats, noteValue, tempo int32) *TimeSignature](#NewTimeSignature)
  * [func (ts *TimeSignature) FromString(sig string) error](#TimeSignature.FromString)
  * [func (ts *TimeSignature) String() string](#TimeSignature.String)
  * [func (ts *TimeSignature) TempoFromDuration(interval time.Duration)](#TimeSignature.TempoFromDuration)
  * [func (ts *TimeSignature) TempoToDuration() time.Duration](#TimeSignature.TempoToDuration)
  * [func (ts *TimeSignature) ToString() string](#TimeSignature.ToString)


#### <a name="pkg-files">Package files</a>
[gnome.go](https://github.com/cognusion/go-gnome/tree/master/gnome.go) [ts.go](https://github.com/cognusion/go-gnome/tree/master/ts.go)



## <a name="pkg-variables">Variables</a>
``` go
var (
    // RPool is a recyclable.BufferPool
    RPool = recyclable.NewBufferPool()
)
```


## <a name="BufferToStreamer">func</a> [BufferToStreamer](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=8482:8560#L305)
``` go
func BufferToStreamer(buff Buffer) (beep.StreamSeekCloser, beep.Format, error)
```
BufferToStreamer checks the first 12 bytes of the Buffer to see if it's a WAV,
and tries to decode it as that or an MP3. Errors are returned if anything fails.



## <a name="FileToBuffer">func</a> [FileToBuffer](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=8054:8116#L291)
``` go
func FileToBuffer(filename string) (*recyclable.Buffer, error)
```
FileToBuffer opens and reads the filename into a Buffer, returning it or an error.



## <a name="FromBPM">func</a> [FromBPM](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=7508:7545#L272)
``` go
func FromBPM(bpm int32) time.Duration
```
FromBPM converts a beats-per-minute tempo to a Microsecond-precise time.Duration.



## <a name="ToBPM">func</a> [ToBPM](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=7769:7809#L282)
``` go
func ToBPM(interval time.Duration) int32
```
ToBPM converts a Microsecond-precise time.Duration to a beats-per-minute tempo.




## <a name="Buffer">type</a> [Buffer](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=1535:1566#L45)
``` go
type Buffer = io.ReadSeekCloser
```
Buffer is our local interface to encapsulate all the interfaces










## <a name="Gnome">type</a> [Gnome](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=2000:2378#L56)
``` go
type Gnome struct {
    // TS tracks and reports the time signature information for the 'gnome.
    TS *TimeSignature
    // contains filtered or unexported fields
}

```
Gnome is a metro...gnome. Get it? Get it?!







### <a name="NewGnomeFromBuffer">func</a> [NewGnomeFromBuffer](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=2531:2621#L74)
``` go
func NewGnomeFromBuffer(buff Buffer, ts *TimeSignature, tickFunc TickFunc) (*Gnome, error)
```
NewGnomeFromBuffer takes a Buffer, a TimeSignature and an optional tickFunc to call when
the 'gnome fires, and gives you a Gnome or an error. :)


### <a name="NewGnomeFromFile">func</a> [NewGnomeFromFile](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=3525:3618#L109)
``` go
func NewGnomeFromFile(soundFile string, ts *TimeSignature, tickFunc TickFunc) (*Gnome, error)
```
NewGnomeFromFile takes a filename, a TimeSignature and an optional tickFunc to call when
the 'gnome fires, and gives you a Gnome or an error. :)





### <a name="Gnome.Change">func</a> (\*Gnome) [Change](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=5348:5383#L181)
``` go
func (g *Gnome) Change(tempo int32)
```
Change sets a new tempo.




### <a name="Gnome.Close">func</a> (\*Gnome) [Close](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=5063:5086#L165)
``` go
func (g *Gnome) Close()
```
Close terminally cleans up all the things.




### <a name="Gnome.IsPaused">func</a> (\*Gnome) [IsPaused](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=5658:5689#L194)
``` go
func (g *Gnome) IsPaused() bool
```
IsPaused returns if the 'gnome is paused.
This should not be used for blocking or timing decisions




### <a name="Gnome.IsRunning">func</a> (\*Gnome) [IsRunning](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=5490:5522#L188)
``` go
func (g *Gnome) IsRunning() bool
```
IsRunning returns if the 'gnome is running.




### <a name="Gnome.Mute">func</a> (\*Gnome) [Mute](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=4834:4856#L153)
``` go
func (g *Gnome) Mute()
```
Mute toggles whether or not audio will play.




### <a name="Gnome.Pause">func</a> (\*Gnome) [Pause](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=5175:5198#L172)
``` go
func (g *Gnome) Pause()
```
Pause toggles whether the 'gnome is paused.




### <a name="Gnome.ReplaceStreamerFromBuffer">func</a> (\*Gnome) [ReplaceStreamerFromBuffer](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=3857:3917#L118)
``` go
func (g *Gnome) ReplaceStreamerFromBuffer(buff Buffer) error
```
ReplaceStreamerFromBuffer attempts to replace the current steamer with a new one from the provided buffer.




### <a name="Gnome.Restart">func</a> (\*Gnome) [Restart](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=4301:4332#L134)
``` go
func (g *Gnome) Restart() error
```
Restart will re-initialize some stopped components so the 'gnome can carry on.




### <a name="Gnome.SetTickFilter">func</a> (\*Gnome) [SetTickFilter](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=5865:5915#L200)
``` go
func (g *Gnome) SetTickFilter(tf TickFilter) error
```
SetTickFilter installs a new tickFilter. A tickFilter is passed the beat count and returns true if
the sound should be played for that beat




### <a name="Gnome.Start">func</a> (\*Gnome) [Start](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=4652:4675#L146)
``` go
func (g *Gnome) Start()
```
Start -s the Gnome. Only works the first time you call it.




### <a name="Gnome.Stop">func</a> (\*Gnome) [Stop](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=4926:4948#L158)
``` go
func (g *Gnome) Stop()
```
Stop stops the Gnome.




## <a name="TickFilter">type</a> [TickFilter](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=1922:1952#L53)
``` go
type TickFilter func(int) bool
```
TickFilter is passed the beat count, and returns true if the tick sound should be played










## <a name="TickFunc">type</a> [TickFunc](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=1805:1828#L50)
``` go
type TickFunc func(int)
```
TickFunc is a function called during a tick event. The current beat number is passed to it.
These functions must absolutely handle being called multiple times in quick succession, handling order,
overlaps, and all of that jazz.










## <a name="TimeSignature">type</a> [TimeSignature](https://github.com/cognusion/go-gnome/tree/master/ts.go?s=209:310#L15)
``` go
type TimeSignature struct {
    Beats     atomic.Int32
    NoteValue atomic.Int32
    Tempo     atomic.Int32
}

```
TimeSignature is a mechaism to make tracking time signatures,
and their changes, easier and safe.







### <a name="NewTimeSignature">func</a> [NewTimeSignature](https://github.com/cognusion/go-gnome/tree/master/ts.go?s=370:437#L22)
``` go
func NewTimeSignature(beats, noteValue, tempo int32) *TimeSignature
```
NewTimeSignature returns an initialized TimeSignature.





### <a name="TimeSignature.FromString">func</a> (\*TimeSignature) [FromString](https://github.com/cognusion/go-gnome/tree/master/ts.go?s=1355:1408#L54)
``` go
func (ts *TimeSignature) FromString(sig string) error
```
FromString takes a signature string (e.g. "4/4" or "6/8") and sets TimeSignature accordingly.




### <a name="TimeSignature.String">func</a> (\*TimeSignature) [String](https://github.com/cognusion/go-gnome/tree/master/ts.go?s=2059:2099#L79)
``` go
func (ts *TimeSignature) String() string
```
String pretty-prints the values of the TimeSignature.




### <a name="TimeSignature.TempoFromDuration">func</a> (\*TimeSignature) [TempoFromDuration](https://github.com/cognusion/go-gnome/tree/master/ts.go?s=901:967#L41)
``` go
func (ts *TimeSignature) TempoFromDuration(interval time.Duration)
```
TempoFromDuration sets the Tempo based on the interval provided.




### <a name="TimeSignature.TempoToDuration">func</a> (\*TimeSignature) [TempoToDuration](https://github.com/cognusion/go-gnome/tree/master/ts.go?s=612:668#L31)
``` go
func (ts *TimeSignature) TempoToDuration() time.Duration
```
TempoToDuration returns the tempo as a time.Duration.




### <a name="TimeSignature.ToString">func</a> (\*TimeSignature) [ToString](https://github.com/cognusion/go-gnome/tree/master/ts.go?s=1143:1185#L49)
``` go
func (ts *TimeSignature) ToString() string
```
ToString returns the signature string (e.g. "4/4" or "6/8"), without the Tempo.








- - -
Generated by [godoc2md](http://github.com/cognusion/godoc2md)
