

# gnome
`import "github.com/cognusion/go-gnome"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package gnome is a library for building a nearly-realtime metro...gnomes. Get it? GET IT?!
One would think there would already be such a thing, but given the complexities involved in getting
decent timing, there was not.

This is not perfect, either. If the system is very busy, the rhythm
will not be smooth. If the tempo is exceptionally high, the rhythm will not be smooth. On a normal
system, doing nothing else, BPMs under 180 are almost always great.

Gnome was built from scratch expecting WASM as the target platform. Some decisions that may seem odd
were made because of that. Most odd decisions are simply. The consuming app was written for my
brother-in-law's music students: [MetroGnome](<a href="https://github.com/cognusion/metrognome">https://github.com/cognusion/metrognome</a>)




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func BufferToStreamer(buff io.ReadSeekCloser) (beep.StreamSeekCloser, beep.Format, error)](#BufferToStreamer)
* [func FileToBuffer(filename string) (*recyclable.Buffer, error)](#FileToBuffer)
* [func FromBPM(bpm int32) time.Duration](#FromBPM)
* [func ToBPM(interval time.Duration) int32](#ToBPM)
* [type Gnome](#Gnome)
  * [func NewGnome(soundFile string, tempo int32) (*Gnome, error)](#NewGnome)
  * [func NewGnomeBuffer(buff io.ReadSeekCloser, tempo int32) (*Gnome, error)](#NewGnomeBuffer)
  * [func NewGnomeBufferTick(buff io.ReadSeekCloser, tempo int32, tickFunc func(int)) (*Gnome, error)](#NewGnomeBufferTick)
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


## <a name="BufferToStreamer">func</a> [BufferToStreamer](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=6811:6900#L269)
``` go
func BufferToStreamer(buff io.ReadSeekCloser) (beep.StreamSeekCloser, beep.Format, error)
```
BufferToStreamer checks the first 12 bytes of the Buffer to see if it's a WAV,
and tries to decode it as that or an MP3. Errors are returned if anything fails.



## <a name="FileToBuffer">func</a> [FileToBuffer](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=6383:6445#L255)
``` go
func FileToBuffer(filename string) (*recyclable.Buffer, error)
```
FileToBuffer opens and reads the filename into a Buffer, returning it or an error.



## <a name="FromBPM">func</a> [FromBPM](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=5837:5874#L236)
``` go
func FromBPM(bpm int32) time.Duration
```
FromBPM converts a beats-per-minute tempo to a Microsecond-precise time.Duration.



## <a name="ToBPM">func</a> [ToBPM](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=6098:6138#L246)
``` go
func ToBPM(interval time.Duration) int32
```
ToBPM converts a Microsecond-precise time.Duration to a beats-per-minute tempo.




## <a name="Gnome">type</a> [Gnome](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=1233:1549#L42)
``` go
type Gnome struct {
    // TS tracks and reports the time signature information for the 'gnome.
    TS *TimeSignature
    // contains filtered or unexported fields
}

```
Gnome is a metro...gnome. Get it? Get it?!







### <a name="NewGnome">func</a> [NewGnome](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=3059:3119#L110)
``` go
func NewGnome(soundFile string, tempo int32) (*Gnome, error)
```
NewGnome loads a WAV or MP3, and plays it every interval, returning an error if there is a problem
loading the file.


### <a name="NewGnomeBuffer">func</a> [NewGnomeBuffer](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=1951:2023#L69)
``` go
func NewGnomeBuffer(buff io.ReadSeekCloser, tempo int32) (*Gnome, error)
```
NewGnomeBuffer takes an io.ReadSeekCloser and a tempo.


### <a name="NewGnomeBufferTick">func</a> [NewGnomeBufferTick](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=1654:1750#L57)
``` go
func NewGnomeBufferTick(buff io.ReadSeekCloser, tempo int32, tickFunc func(int)) (*Gnome, error)
```
NewGnomeBufferTick takes an io.ReadSeekCloser, tempo, and a tickFunc to call when the 'gnome fires.


### <a name="NewGnomeWithTickFunc">func</a> [NewGnomeWithTickFunc](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=2729:2821#L98)
``` go
func NewGnomeWithTickFunc(soundFile string, tempo int32, tickFunc func(int)) (*Gnome, error)
```
NewGnomeWithTickFunc takes a file string, tempo, and a tickFunc to call when the 'gnome fires.





### <a name="Gnome.Change">func</a> (\*Gnome) [Change](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=4246:4281#L163)
``` go
func (g *Gnome) Change(tempo int32)
```
Change sets a new tempo.




### <a name="Gnome.Close">func</a> (\*Gnome) [Close](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=4022:4045#L149)
``` go
func (g *Gnome) Close()
```
Close terminally cleans up all the things.




### <a name="Gnome.IsRunning">func</a> (\*Gnome) [IsRunning](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=4388:4420#L170)
``` go
func (g *Gnome) IsRunning() bool
```
IsRunning returns if the 'gnome is running.




### <a name="Gnome.Mute">func</a> (\*Gnome) [Mute](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=3816:3838#L138)
``` go
func (g *Gnome) Mute()
```
Mute toggles whether or not audio will play.




### <a name="Gnome.Pause">func</a> (\*Gnome) [Pause](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=4117:4140#L155)
``` go
func (g *Gnome) Pause()
```
Pause toggles whether the 'gnome is paused.




### <a name="Gnome.Restart">func</a> (\*Gnome) [Restart](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=3319:3350#L119)
``` go
func (g *Gnome) Restart() error
```
Restart will re-initialize some stopped components so the 'gnome can carry on.




### <a name="Gnome.Start">func</a> (\*Gnome) [Start](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=3652:3675#L131)
``` go
func (g *Gnome) Start()
```
Start -s the Gnome. Only works the first time you call it.




### <a name="Gnome.Stop">func</a> (\*Gnome) [Stop](https://github.com/cognusion/go-gnome/tree/master/gnome.go?s=3908:3930#L143)
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










### <a name="TimeSignature.FromString">func</a> (\*TimeSignature) [FromString](https://github.com/cognusion/go-gnome/tree/master/ts.go?s=1112:1165#L45)
``` go
func (ts *TimeSignature) FromString(sig string) error
```
FromString takes a signature string (e.g. "4/4" or "6/8") and sets TimeSignature accordingly.




### <a name="TimeSignature.String">func</a> (\*TimeSignature) [String](https://github.com/cognusion/go-gnome/tree/master/ts.go?s=1816:1856#L70)
``` go
func (ts *TimeSignature) String() string
```
String pretty-prints the values of the TimeSignature.




### <a name="TimeSignature.TempoFromDuration">func</a> (\*TimeSignature) [TempoFromDuration](https://github.com/cognusion/go-gnome/tree/master/ts.go?s=658:724#L32)
``` go
func (ts *TimeSignature) TempoFromDuration(interval time.Duration)
```
TempoFromDuration sets the Tempo based on the interval provided.




### <a name="TimeSignature.TempoToDuration">func</a> (\*TimeSignature) [TempoToDuration](https://github.com/cognusion/go-gnome/tree/master/ts.go?s=369:425#L22)
``` go
func (ts *TimeSignature) TempoToDuration() time.Duration
```
TempoToDuration returns the tempo as a time.Duration.




### <a name="TimeSignature.ToString">func</a> (\*TimeSignature) [ToString](https://github.com/cognusion/go-gnome/tree/master/ts.go?s=900:942#L40)
``` go
func (ts *TimeSignature) ToString() string
```
ToString returns the signature string (e.g. "4/4" or "6/8"), without the Tempo.








- - -
Generated by [godoc2md](http://github.com/cognusion/godoc2md)
