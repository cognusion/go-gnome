

# gnome
`import "github.com/cognusion/go-gnome"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package gnome is a library for building a nearly-realtime metro...gnomes. Get it? GET IT?!
One would think there would already be such a thing, but given the complexities involved in getting
decent timing, there was not. Supports WAV, MP3, FLAC, and Vorbis as passthoughs from `gopxl/beep`
(which stands on the shoulders of other giants).

This is not perfect, either. If the system is very busy, the rhythm will not be smooth.
If the tempo is exceptionally high, the rhythm will not be smooth.
On a normal system, doing nothing else, BPMs under 180 are almost always great.

Gnome was built from scratch expecting WASM as the target platform. Some decisions that may seem odd
were made because of that. Most odd decisions are simply. The consuming app was written for my
brother-in-law's music students: [MetroGnome](<a href="https://github.com/cognusion/metrognome">https://github.com/cognusion/metrognome</a>).




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func BufferToStreamer(buff Buffer) (beep.StreamSeekCloser, beep.Format, error)](#BufferToStreamer)
* [func FileToBuffer(filename string) (*recyclable.Buffer, error)](#FileToBuffer)
* [func FromBPM(bpm int32) time.Duration](#FromBPM)
* [func ToBPM(interval time.Duration) int32](#ToBPM)
* [type Buffer](#Buffer)
* [type Gnome](#Gnome)
  * [func NewGnome(soundFile string, tempo int32) (*Gnome, error)](#NewGnome)
  * [func NewGnomeBuffer(buff Buffer, tempo int32) (*Gnome, error)](#NewGnomeBuffer)
  * [func NewGnomeBufferTick(buff Buffer, tempo int32, tickFunc func(int)) (*Gnome, error)](#NewGnomeBufferTick)
  * [func NewGnomeFromBuffer(buff Buffer, ts *TimeSignature, tickFunc func(int)) (*Gnome, error)](#NewGnomeFromBuffer)
  * [func NewGnomeFromFile(soundFile string, ts *TimeSignature, tickFunc func(int)) (*Gnome, error)](#NewGnomeFromFile)
  * [func NewGnomeWithTickFunc(soundFile string, tempo int32, tickFunc func(int)) (*Gnome, error)](#NewGnomeWithTickFunc)
  * [func (g *Gnome) Change(tempo int32)](#Gnome.Change)
  * [func (g *Gnome) Close()](#Gnome.Close)
  * [func (g *Gnome) IsRunning() bool](#Gnome.IsRunning)
  * [func (g *Gnome) Mute()](#Gnome.Mute)
  * [func (g *Gnome) Pause()](#Gnome.Pause)
  * [func (g *Gnome) Restart() error](#Gnome.Restart)
  * [func (g *Gnome) Start()](#Gnome.Start)
  * [func (g *Gnome) Stop()](#Gnome.Stop)
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


## <a name="BufferToStreamer">func</a> [BufferToStreamer](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=7877:7955#L284)
``` go
func BufferToStreamer(buff Buffer) (beep.StreamSeekCloser, beep.Format, error)
```
BufferToStreamer checks the first 12 bytes of the Buffer to see if it's a WAV,
and tries to decode it as that or an MP3. Errors are returned if anything fails.



## <a name="FileToBuffer">func</a> [FileToBuffer](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=7449:7511#L270)
``` go
func FileToBuffer(filename string) (*recyclable.Buffer, error)
```
FileToBuffer opens and reads the filename into a Buffer, returning it or an error.



## <a name="FromBPM">func</a> [FromBPM](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=6903:6940#L251)
``` go
func FromBPM(bpm int32) time.Duration
```
FromBPM converts a beats-per-minute tempo to a Microsecond-precise time.Duration.



## <a name="ToBPM">func</a> [ToBPM](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=7164:7204#L261)
``` go
func ToBPM(interval time.Duration) int32
```
ToBPM converts a Microsecond-precise time.Duration to a beats-per-minute tempo.




## <a name="Buffer">type</a> [Buffer](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=1462:1493#L45)
``` go
type Buffer = io.ReadSeekCloser
```
Buffer is our local interface to encapsulate all the interfaces










## <a name="Gnome">type</a> [Gnome](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=1541:1857#L48)
``` go
type Gnome struct {
    // TS tracks and reports the time signature information for the 'gnome.
    TS *TimeSignature
    // contains filtered or unexported fields
}

```
Gnome is a metro...gnome. Get it? Get it?!







### <a name="NewGnome">func</a> [NewGnome](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=3774:3834#L118)
``` go
func NewGnome(soundFile string, tempo int32) (*Gnome, error)
```
NewGnome loads a WAV or MP3, and plays it every interval, returning an error if there is a problem
loading the file.

Deprecated: Migrate to NewGnomeFromFile


### <a name="NewGnomeBuffer">func</a> [NewGnomeBuffer](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=2269:2330#L72)
``` go
func NewGnomeBuffer(buff Buffer, tempo int32) (*Gnome, error)
```
NewGnomeBuffer takes an GnomeBuffer and a tempo.

Deprecated: Migrate to NewGnomeFromBuffer


### <a name="NewGnomeBufferTick">func</a> [NewGnomeBufferTick](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=2004:2089#L65)
``` go
func NewGnomeBufferTick(buff Buffer, tempo int32, tickFunc func(int)) (*Gnome, error)
```
NewGnomeBufferTick takes an GnomeBuffer, tempo, and a tickFunc to call when the 'gnome fires.

Deprecated: Migrate to NewGnomeFromBuffer


### <a name="NewGnomeFromBuffer">func</a> [NewGnomeFromBuffer](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=2556:2647#L78)
``` go
func NewGnomeFromBuffer(buff Buffer, ts *TimeSignature, tickFunc func(int)) (*Gnome, error)
```
NewGnomeFromBuffer takes a Buffer, a TimeSignature and an optional tickFunc to call when
the 'gnome fires, and gives you a Gnome or an error. :)


### <a name="NewGnomeFromFile">func</a> [NewGnomeFromFile](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=4063:4157#L124)
``` go
func NewGnomeFromFile(soundFile string, ts *TimeSignature, tickFunc func(int)) (*Gnome, error)
```
NewGnomeFromFile takes a filename, a TimeSignature and an optional tickFunc to call when
the 'gnome fires, and gives you a Gnome or an error. :)


### <a name="NewGnomeWithTickFunc">func</a> [NewGnomeWithTickFunc](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=3430:3522#L110)
``` go
func NewGnomeWithTickFunc(soundFile string, tempo int32, tickFunc func(int)) (*Gnome, error)
```
NewGnomeWithTickFunc takes a file string, tempo, and a tickFunc to call when the 'gnome fires.

Deprecated: Migrate to NewGnomeFromFile





### <a name="Gnome.Change">func</a> (\*Gnome) [Change](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=5312:5347#L178)
``` go
func (g *Gnome) Change(tempo int32)
```
Change sets a new tempo.




### <a name="Gnome.Close">func</a> (\*Gnome) [Close](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=5071:5094#L163)
``` go
func (g *Gnome) Close()
```
Close terminally cleans up all the things.




### <a name="Gnome.IsRunning">func</a> (\*Gnome) [IsRunning](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=5454:5486#L185)
``` go
func (g *Gnome) IsRunning() bool
```
IsRunning returns if the 'gnome is running.




### <a name="Gnome.Mute">func</a> (\*Gnome) [Mute](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=4865:4887#L152)
``` go
func (g *Gnome) Mute()
```
Mute toggles whether or not audio will play.




### <a name="Gnome.Pause">func</a> (\*Gnome) [Pause](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=5183:5206#L170)
``` go
func (g *Gnome) Pause()
```
Pause toggles whether the 'gnome is paused.




### <a name="Gnome.Restart">func</a> (\*Gnome) [Restart](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=4368:4399#L133)
``` go
func (g *Gnome) Restart() error
```
Restart will re-initialize some stopped components so the 'gnome can carry on.




### <a name="Gnome.Start">func</a> (\*Gnome) [Start](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=4701:4724#L145)
``` go
func (g *Gnome) Start()
```
Start -s the Gnome. Only works the first time you call it.




### <a name="Gnome.Stop">func</a> (\*Gnome) [Stop](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=4957:4979#L157)
``` go
func (g *Gnome) Stop()
```
Stop stops the Gnome.




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
