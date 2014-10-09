package sampleplayer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/metakeule/music"

	"github.com/metakeule/music/dyn"

	"github.com/metakeule/music/note"

	"sort"
)

type Sampler interface {
	SampleAndStrech(dyn float64, freq float64) (path string, strech float64)
}

// dynMapper goes through the given map in sorted order, returning the string for the largest
// dyn value in m that is <= dyn
// It can be used to map dynlitudes to strings that are part of a sample name
func dynMapper(m map[float64]string, dyn float64) string {
	if len(m) == 0 {
		panic("empty map not allowed")
	}
	fl := []float64{}

	for f := range m {
		fl = append(fl, f)
	}

	sort.Float64s(fl)

	last := m[fl[0]]

	for _, f := range fl {
		if f > dyn {
			break
		}
		last = m[f]
	}
	return last
}

// freqMapper goes through the given map in sorted order, returning the string for the largest
// freq value in m that is <= freq
// It can be used to map frequencies to strings that are part of a sample name
// it returns the frequency that was found and the string
func freqMapper(m map[float64]string, freq float64) (string, float64) {
	if len(m) == 0 {
		panic("empty map not allowed")
	}
	fl := []float64{}

	for f := range m {
		fl = append(fl, f)
	}

	sort.Float64s(fl)

	lastFreq := fl[0]
	lastStr := m[lastFreq]

	for _, f := range fl {
		if f > freq {
			break
		}
		lastFreq = f
		lastStr = m[lastFreq]
	}
	return lastStr, lastFreq
}

func paramMapper(m map[float64]string, value float64) (string, float64) {
	if len(m) == 0 {
		panic("empty map not allowed")
	}
	fl := []float64{}

	for f := range m {
		fl = append(fl, f)
	}

	sort.Float64s(fl)

	lastVal := fl[0]
	lastStr := m[lastVal]

	for _, f := range fl {
		if f > value {
			break
		}
		lastVal = f
		lastStr = m[lastVal]
	}
	return lastStr, lastVal
}

type Parameter interface {
	Name() string // name of the placeholder
	// Value calculates the parameter value from a name
	Value(name string) float64
}

type Scanner interface {
	// maps parameter name to a map of parameter values to names
	Scan(string) map[string]map[float64]string
}

type Instrument struct {
	Dir         string
	Format      string // must be a string like `flute_{{freq}}_{{dur}}_{{amp}}_normal.wav`, e.g. flute_Gs4_15_pianissimo_normal.wav
	AllVariants map[string]map[float64]string
}

func CalculatePhilharmonicFreq(s string) float64 {
	l := len(s)
	if l != 2 && l != 3 {
		panic("freq string is too long, must be 2-3 characters, is " + s)
	}

	i, err := strconv.ParseInt(s[l-1:l], 10, 16)
	if err != nil {
		panic(err.Error())
	}

	nt := s[:l-1]

	var midiBase note.Note
	switch nt {
	case "C":
		midiBase = note.C4
	case "Cs":
		midiBase = note.Cis4
	case "D":
		midiBase = note.D4
	case "Ds":
		midiBase = note.Dis4
	case "E":
		midiBase = note.E4
	case "F":
		midiBase = note.F4
	case "Fis":
		midiBase = note.Fis4
	case "G":
		midiBase = note.G4
	case "Gs":
		midiBase = note.Gis4
	case "A":
		midiBase = note.A4
	case "As":
		midiBase = note.Ais4
	case "B":
		midiBase = note.B4
	}
	return note.MidiCps(float64((i-4)*12) + float64(midiBase))
}

var Long = music.Dur(1000)
var VeryLong = music.Dur(5000)
var Phrase = music.Dur(10000)

func CalculatePhilharmonicDur(s string) float64 {
	switch s {
	case "long":
		return 1000
	case "very-long":
		return 5000
	case "phrase":
		return 10000
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err.Error())
	}
	return f
}

type dynamic struct {
	name  string
	value float64
}

func (d dynamic) Name() string   { return d.name }
func (d dynamic) Value() float64 { return d.value }

func (d dynamic) Params() map[string]float64 {
	return music.Dyn(d.value).Params()
}

var (
	CrescendoDecrescendo = dynamic{"cresc-decresc", dyn.CrescendoDecrescendo.Value()}
	Crescendo            = dynamic{"crescendo", dyn.Crescendo.Value()}
	MoltoPianissimo      = dynamic{"molto-pianissimo", dyn.PianissimoPiano.Value()}
	Pianissimo           = dynamic{"pianissimo", dyn.Pianissimo.Value()}
	Piano                = dynamic{"piano", dyn.Piano.Value()}
	MezzoPiano           = dynamic{"mezzo-piano", dyn.Mezzopiano.Value()}
	MezzoForte           = dynamic{"mezzo-forte", dyn.Mezzoforte.Value()}
	Forte                = dynamic{"forte", dyn.Forte.Value()}
	Fortissimo           = dynamic{"fortissimo", dyn.Fortissimo.Value()}
	Decrescendo          = dynamic{"decrescendo", dyn.Decrescendo.Value()}
)

func CalculatePhilharmonicAmp(s string) float64 {
	switch s {
	case CrescendoDecrescendo.Name():
		return CrescendoDecrescendo.Value()
	case Crescendo.Name():
		return Crescendo.Value()
	case MoltoPianissimo.Name():
		return MoltoPianissimo.Value()
	case Pianissimo.Name():
		return Pianissimo.Value()
	case Piano.Name():
		return Piano.Value()
	case MezzoPiano.Name():
		return MezzoPiano.Value()
	case MezzoForte.Name():
		return MezzoForte.Value()
	case Forte.Name():
		return Forte.Value()
	case Fortissimo.Name():
		return Fortissimo.Value()
	case Decrescendo.Name():
		return Decrescendo.Value()
	}
	panic("unknown dyn string " + s)
}

type variant struct {
	name  string
	value float64
}

func (v variant) Name() string   { return v.name }
func (v variant) Value() float64 { return v.value }
func (v variant) Params() map[string]float64 {
	return music.Param("variant", v.value).Params()
}

func newVariant(name string) variant {
	i++
	return variant{name, i}
}

var (
	i      float64
	Normal = variant{"normal", i}

	Nonlegato  = newVariant("nonlegato")
	Legato     = newVariant("legato")
	ArcoLegato = newVariant("arco-legato")

	NonVibrato   = newVariant("non-vibrato")
	Vibrato      = newVariant("vibrato")
	MoltoVibrato = newVariant("molto-vibrato")

	ArcoTenuto = newVariant("arco-tenuto")
	Tenuto     = newVariant("tenuto")

	Tremolo     = newVariant("tremolo")
	ArcoTremolo = newVariant("arco-tremolo")
	PizzTremolo = newVariant("pizz-tremolo")

	Staccato      = newVariant("staccato")
	ArcoStaccato  = newVariant("arco-staccato")
	Staccatissimo = newVariant("staccatissimo")

	MajorTrill     = newVariant("major-trill")
	MinorTrill     = newVariant("minor-trill")
	ArcoMajorTrill = newVariant("arco-major-trill")
	ArcoMinorTrill = newVariant("arco-minor-trill")

	Harmonics          = newVariant("harmonics")
	Harmonic           = newVariant("harmonic")
	ArcoHarmonic       = newVariant("arco-harmonic")
	ArtificialHarmonic = newVariant("artificial-harmonic")
	NaturalHarmonic    = newVariant("natural-harmonic")

	Glissando         = newVariant("glissando")
	ArcoGlissando     = newVariant("arco-glissando")
	PizzGlissando     = newVariant("pizz-glissando")
	HarmonicGlissando = newVariant("harmonic-glissando")

	ArcoNormal          = newVariant("arco-normal")
	ArcoDetache         = newVariant("arco-detache")
	ArcoSpiccato        = newVariant("arco-spiccato")
	ArcoColLegnoBattuto = newVariant("arco-col-legno-battuto")
	ArcoColLegnoTratto  = newVariant("arco-col-legno-tratto")
	ArcoAuTalon         = newVariant("arco-au-talon")
	ArcoMartele         = newVariant("arco-martele")
	ArcoSulPonticello   = newVariant("arco-sul-ponticello")
	ArcoPuntaDarco      = newVariant("arco-punta-d'arco")
	ArcoSulTasto        = newVariant("arco-sul-tasto")
	ArcoPortato         = newVariant("arco-portato")

	PizzNormal      = newVariant("pizz-normal")
	SnapPizz        = newVariant("snap-pizz")
	PizzQuasiGuitar = newVariant("pizz-quasi-guitar")

	TonguedSlur     = newVariant("tongued-slur")
	TripleTonguing  = newVariant("triple-tonguing")
	DoubleTonguing  = newVariant("double-tonguing")
	FlutterTonguing = newVariant("fluttertonguing")

	ConSord = newVariant("con-sord")
	Mute    = newVariant("mute")
	Subtone = newVariant("subtone")
	String  = newVariant("string")
)

func CalculatePhilharmonicVariant(s string) float64 {
	switch s {
	case Normal.Name():
		return Normal.Value()
	case ArcoNormal.Name():
		return ArcoNormal.Value()
	case PizzNormal.Name():
		return PizzNormal.Value()
	case MoltoVibrato.Name():
		return MoltoVibrato.Value()
	case NonVibrato.Name():
		return NonVibrato.Value()
	case ArcoGlissando.Name():
		return ArcoGlissando.Value()
	case SnapPizz.Name():
		return SnapPizz.Value()
	case ArcoMajorTrill.Name():
		return ArcoMajorTrill.Value()
	case ArcoMinorTrill.Name():
		return ArcoMinorTrill.Value()
	case ArcoDetache.Name():
		return ArcoDetache.Value()
	case ArcoLegato.Name():
		return ArcoLegato.Value()
	case ArtificialHarmonic.Name():
		return ArtificialHarmonic.Value()
	case NaturalHarmonic.Name():
		return NaturalHarmonic.Value()
	case ArcoSpiccato.Name():
		return ArcoSpiccato.Value()
	case ArcoStaccato.Name():
		return ArcoStaccato.Value()
	case ArcoTremolo.Name():
		return ArcoTremolo.Value()
	case PizzGlissando.Name():
		return PizzGlissando.Value()
	case ArcoAuTalon.Name():
		return ArcoAuTalon.Value()
	case ArcoColLegnoBattuto.Name():
		return ArcoColLegnoBattuto.Value()
	case ArcoColLegnoTratto.Name():
		return ArcoColLegnoTratto.Value()
	case ArcoSulPonticello.Name():
		return ArcoSulPonticello.Value()
	case ArcoSulTasto.Name():
		return ArcoSulTasto.Value()
	case ConSord.Name():
		return ConSord.Value()
	case ArcoMartele.Name():
		return ArcoMartele.Value()
	case ArcoTenuto.Name():
		return ArcoTenuto.Value()
	case HarmonicGlissando.Name():
		return HarmonicGlissando.Value()
	case ArcoPuntaDarco.Name():
		return ArcoPuntaDarco.Value()
	case PizzQuasiGuitar.Name():
		return PizzQuasiGuitar.Value()
	case PizzTremolo.Name():
		return PizzTremolo.Value()
	case Mute.Name():
		return Mute.Value()
	case Vibrato.Name():
		return Vibrato.Value()
	case Glissando.Name():
		return Glissando.Value()
	case ArcoHarmonic.Name():
		return ArcoHarmonic.Value()
	case Harmonics.Name():
		return Harmonics.Value()
	case ArcoPortato.Name():
		return ArcoPortato.Value()
	case Harmonic.Name():
		return Harmonic.Value()
	case Subtone.Name():
		return Subtone.Value()
	case String.Name():
		return String.Value()
	case Staccatissimo.Name():
		return Staccatissimo.Value()
	case Staccato.Name():
		return Staccato.Value()
	case Tenuto.Name():
		return Tenuto.Value()
	case TonguedSlur.Name():
		return TonguedSlur.Value()
	case TripleTonguing.Name():
		return TripleTonguing.Value()
	case Tremolo.Name():
		return Tremolo.Value()
	case Legato.Name():
		return Legato.Value()
	case Nonlegato.Name():
		return Nonlegato.Value()
	case DoubleTonguing.Name():
		return DoubleTonguing.Value()
	case FlutterTonguing.Name():
		return FlutterTonguing.Value()
	case MajorTrill.Name():
		return MajorTrill.Value()
	case MinorTrill.Name():
		return MinorTrill.Value()
	}
	panic("unknown variant string " + s)
}

type PhilharmonicScanner struct {
	registry map[string]map[float64]string

	//           instrument  variant      dur        dyn          freq    fileextension
	registry2     map[string]map[float64]map[float64]map[float64]map[float64]string
	baseDir       string
	loadedSamples map[string]struct{}

	//registry3 map[string]map[[3]float64]float64
}

func (p *PhilharmonicScanner) AddToRegistry(instrument string, variant, dur, dyn, freq float64, ext string) {
	/*
			fmt.Printf(`adding to registry:
		instrument: %s
		variant: %f
		dur: %f
		dyn: %f
		freq: %f
		ext: %s

		`, instrument, variant, dur, dyn, freq, ext)
	*/

	instrM, instrOk := p.registry2[instrument]
	if !instrOk {
		instrM = map[float64]map[float64]map[float64]map[float64]string{}
		p.registry2[instrument] = instrM
	}

	varM, varOk := instrM[variant]
	if !varOk {
		varM = map[float64]map[float64]map[float64]string{}
		instrM[variant] = varM
	}

	durM, durOk := varM[dur]
	if !durOk {
		durM = map[float64]map[float64]string{}
		varM[dur] = durM
	}

	dynM, dynOk := durM[dyn]
	if !dynOk {
		dynM = map[float64]string{}
		durM[dyn] = dynM
	}

	dynM[freq] = ext
}

func NewPhilharmonicScanner(baseDir string) *PhilharmonicScanner {
	return &PhilharmonicScanner{
		registry:      map[string]map[float64]string{},
		registry2:     map[string]map[float64]map[float64]map[float64]map[float64]string{},
		baseDir:       baseDir,
		loadedSamples: map[string]struct{}{},
	}
}

var _ music.SampleLibrary = &PhilharmonicScanner{}

func (p *PhilharmonicScanner) Channels() []int {
	return []int{1, 2}
}

/*
type SampleLibrary interface {
	SamplePath(instrument string, params map[string]float64) string
	Channels() []int // channel variants
}
*/

// instrument => variant => dur => dyn => frequency
func (p *PhilharmonicScanner) Scan(filename string) {
	// fmt.Printf("scanning %s\n", filename)
	ext := filepath.Ext(filename)
	filename = strings.TrimSuffix(filename, ext)
	//  flute_Gs4_15_pianissimo_normal.wav
	arr := strings.Split(filename, "_")
	// res := map[string]map[float64]string{}
	instrument := arr[0]
	// res["instrument"] = map[float64]string{0: arr[0]}
	freq := CalculatePhilharmonicFreq(arr[1])
	// res["freq"] = map[float64]string{CalculatePhilharmonicFreq(arr[1]): arr[1]}
	// res["dur"] = map[float64]string{CalculatePhilharmonicDur(arr[2]): arr[2]}
	dur := CalculatePhilharmonicDur(arr[2])
	// res["dyn"] = map[float64]string{CalculatePhilharmonicAmp(arr[3]): arr[3]}
	dyn := CalculatePhilharmonicAmp(arr[3])
	// res["variant"] = map[float64]string{CalculatePhilharmonicVariant(arr[4]): arr[4]}
	variant := CalculatePhilharmonicVariant(arr[4])
	// res["ext"] = map[float64]string{0: ext} // file extension

	p.AddToRegistry(instrument, variant, dur, dyn, freq, ext)

	res := map[string]map[float64]string{}
	res["instrument"] = map[float64]string{0: arr[0]}
	res["freq"] = map[float64]string{CalculatePhilharmonicFreq(arr[1]): arr[1]}
	res["dur"] = map[float64]string{CalculatePhilharmonicDur(arr[2]): arr[2]}
	res["dyn"] = map[float64]string{CalculatePhilharmonicAmp(arr[3]): arr[3]}
	res["variant"] = map[float64]string{CalculatePhilharmonicVariant(arr[4]): arr[4]}
	res["ext"] = map[float64]string{0: ext} // file extension

	// fmt.Printf("%#v\n", res)

	for param, vals := range res {
		rVals, ok := p.registry[param]
		if ok {
			for k, v := range vals {
				rVals[k] = v
			}
		} else {
			p.registry[param] = vals
		}
	}

}

func (p *PhilharmonicScanner) Scan2(filename string) map[string]map[float64]string {
	fmt.Printf("scanning %s\n", filename)
	ext := filepath.Ext(filename)
	filename = strings.TrimSuffix(filename, ext)
	//  flute_Gs4_15_pianissimo_normal.wav
	arr := strings.Split(filename, "_")
	res := map[string]map[float64]string{}
	res["instrument"] = map[float64]string{0: arr[0]}
	res["freq"] = map[float64]string{CalculatePhilharmonicFreq(arr[1]): arr[1]}
	res["dur"] = map[float64]string{CalculatePhilharmonicDur(arr[2]): arr[2]}
	res["dyn"] = map[float64]string{CalculatePhilharmonicAmp(arr[3]): arr[3]}
	res["variant"] = map[float64]string{CalculatePhilharmonicVariant(arr[4]): arr[4]}
	res["ext"] = map[float64]string{0: ext} // file extension

	for param, vals := range res {
		rVals, ok := p.registry[param]
		if ok {
			for k, v := range vals {
				rVals[k] = v
			}
		} else {
			p.registry[param] = vals
		}
	}

	return res
}

// TODO lookup the nearest variant that exists in instrM
func (p *PhilharmonicScanner) FallbackVariant(variant float64, instrM map[float64]map[float64]map[float64]map[float64]string) (string, map[float64]map[float64]map[float64]string) {
	possible := map[float64]string{}

	// fmt.Printf("searching for variant %f in %#v\n", variant, instrM)

	for v := range instrM {
		if v == variant {
			return p.registry["variant"][v], instrM[v]
		}
		possible[v] = p.registry["variant"][v]
	}

	str, val := paramMapper(possible, variant)
	return str, instrM[val]
	// return Normal.Name(), instrM[Normal.Value()]
}

// TODO lookup the nearest dur that exists in variantM
func (p *PhilharmonicScanner) FallbackDur(dur float64, variantM map[float64]map[float64]map[float64]string) (string, map[float64]map[float64]string) {
	possible := map[float64]string{}

	for v := range variantM {
		if v == dur {
			return p.registry["dur"][v], variantM[v]
		}
		possible[v] = p.registry["dur"][v]
	}

	str, val := paramMapper(possible, dur)
	return str, variantM[val]

	// return "1", variantM[1]
}

// TODO lookup the nearest dyn that exists in durM
func (p *PhilharmonicScanner) FallbackDyn(dyn float64, durM map[float64]map[float64]string) (string, map[float64]string) {
	possible := map[float64]string{}

	for v := range durM {
		if v == dyn {
			return p.registry["dyn"][v], durM[v]
		}
		possible[v] = p.registry["dyn"][v]
	}

	str, val := paramMapper(possible, dyn)
	return str, durM[val]

	// return MezzoForte.Name(), durM[MezzoForte.Value()]
}

// TODO lookup the nearest freq that exists in dynM
func (p *PhilharmonicScanner) FallbackFreq(freq float64, dynM map[float64]string) (freqStr string, realFreq float64, ext string) {
	possible := map[float64]string{}

	for v := range dynM {
		if v == freq {
			return p.registry["freq"][v], freq, dynM[v]
		}
		possible[v] = p.registry["freq"][v]
	}

	str, val := paramMapper(possible, freq)
	return str, val, dynM[val]

	//return "A1", note.A.Frequency(), ".wav"
}

func (p *PhilharmonicScanner) SamplePath(instrument string, params map[string]float64) string {
	instrM, instrOk := p.registry2[instrument]
	if !instrOk {
		p.ScanDir(filepath.Join(p.baseDir, instrument))
		instrM, instrOk = p.registry2[instrument]
		if !instrOk {
			panic("instrument " + instrument + " not found")
		}
	}

	variant, variantDefined := params["variant"]
	if !variantDefined {
		variant = Normal.Value()
	}

	variantStr, variantM := p.FallbackVariant(variant, instrM)

	dur, durDefined := params["dur"]
	if !durDefined {
		dur = 1.0
	}

	durStr, durM := p.FallbackDur(dur, variantM)

	dyn, dynDefined := params["dyn"]
	if !dynDefined {
		dyn = MezzoForte.Value()
	}

	dynStr, dynM := p.FallbackDyn(dyn, durM)

	freq, freqDefined := params["freq"]
	if !freqDefined {
		freq = note.A4.Frequency()
	}

	freqStr, realFreq, ext := p.FallbackFreq(freq, dynM)

	params["samplefreq"] = realFreq

	if realFreq == freq {
		params["rate"] = float64(1)
	} else {
		params["rate"] = freq / realFreq // "scale" the sample according to the frequency
	}

	sPath := filepath.Join(p.baseDir, instrument, fmt.Sprintf(`%s_%s_%s_%s_%s%s`, instrument, freqStr, durStr, dynStr, variantStr, ext))
	_, err := os.Stat(sPath)
	if err != nil {
		panic("can't find sample " + sPath)
	}

	if _, has := p.loadedSamples[sPath]; !has {
		// fmt.Printf("using %s\n", sPath)
		p.loadedSamples[sPath] = struct{}{}
	}

	return sPath
}

// SampleForParams returns the filename for the given params and changes the params if necessary
func (p *PhilharmonicScanner) SampleForParams2(instrument string, params map[string]float64) string {
	c := map[string]string{}

	for name, val := range params {

		possible, has := p.registry[name]

		if !has {
			panic("unknown parameter " + name)
		}

		cKey, res := paramMapper(possible, val)
		if name == "freq" {
			if res == val {
				params["rate"] = float64(1)
			} else {
				params["rate"] = val / res // "scale" the sample according to the frequency
			}
		}
		c[name] = cKey
	}

	c["instrument"] = instrument

	// flute_Gs4_15_pianissimo_normal.wav

	return fmt.Sprintf(`%s_%s_%s_%s_%s.%s`, c["instrument"], c["freq"], c["dur"], c["dyn"], c["variant"], c["ext"])
}

func (p *PhilharmonicScanner) ScanDir(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !f.IsDir() {
			ext := strings.ToLower(filepath.Ext(f.Name()))
			if ext == ".wav" || ext == ".aiff" {
				p.Scan(f.Name())
			}
		}
	}
	return nil
}

/*
probably a better way would be with NRT as meantioned in
http://new-supercollider-mailing-lists-forums-use-these.2681727.n2.nabble.com/max-amp-time-location-in-a-sound-file-td7581043.html#a7581256
(where this code is taken from)
here is the example that they are talking about
http://doc.sccode.org/Guides/Non-Realtime-Synthesis.html

echo '"hi".postln;' | sclang -r -s
or (for std input)
sclang -r -s -
to execute the code with sclang
*/
/*
func scCodeForOffset(file string) string {
	return fmt.Sprintf(
		//		`f = SoundFile.new; f.openRead(%#v); a = FloatArray.newClear(f.numFrames * f.numChannels); f.readData(a); ~time = a.abs.maxIndex / a.size * f.duration; ~time.postln;`,
		`p = %#v.pathMatch;  p.do({ |y|  try {  f = SoundFile.new; f.openRead(y); a = FloatArray.newClear(f.numFrames * f.numChannels); f.readData(a); ~time = a.abs.maxIndex / a.size * f.duration; f.close(); o = ();  o.offset = ~time; g=File(y+".offset", "w+"); g.write(JSON.stringify(o)); g.close; } { |error| 	}; });`,
		file,
	)
}
*/

/*
  p = "/media/usb0/philharmonic-orchestra/flute/*.wav".pathMatch;
  p.do({ |y|
	try {
    f = SoundFile.new;
		f.openRead(y);
		a = FloatArray.newClear(f.numFrames * f.numChannels);
		f.readData(a);
		~time = a.abs.maxIndex / a.size * f.duration;
		f.close();
		o = ();
		o.offset = ~time;
		g=File(y+".offset", "w+");
		g.write(JSON.stringify(o));
		g.close;
	} { |error|
		error.postln
	};
  });

*/
