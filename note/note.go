// Package amp provides shortcuts for midi notes and their frequencies
package note

import (
	"fmt"
	"math"
	"strconv"
)

// stolen from https://groups.google.com/forum/?fromgroups=#!topic/golang-nuts/ITZV08gAugI
// return rounded version of x with prec precision.
func roundViaFloat(x float64, prec int) float64 {
	frep := strconv.FormatFloat(x, 'f', prec, 64)
	f, _ := strconv.ParseFloat(frep, 64)
	return f
}

// roundDec rounds to n places after the decimal point
// examples
//  roundDec( 200.234, 2) =>  200.23
//  roundDec(4200.234, 2) => 4200.23
func roundDec(x float64, n int) float64 {
	nonDecs := len(fmt.Sprintf("%.0f", x))
	return roundViaFloat(x, nonDecs+n)
}

type Note float64

const (
	C0 Note = (12 + iota)
	Cis0
	D0
	Dis0
	E0
	F0
	Fis0
	G0
	Gis0
	A0
	Ais0
	B0
	C1
	Cis1
	D1
	Dis1
	E1
	F1
	Fis1
	G1
	Gis1
	A1
	Ais1
	B1
	C2
	Cis2
	D2
	Dis2
	E2
	F2
	Fis2
	G2
	Gis2
	A2
	Ais2
	B2
	C3
	Cis3
	D3
	Dis3
	E3
	F3
	Fis3
	G3
	Gis3
	A3
	Ais3
	B3
	C4
	Cis4
	D4
	Dis4
	E4
	F4
	Fis4
	G4
	Gis4
	A4
	Ais4
	B4
	C5
	Cis5
	D5
	Dis5
	E5
	F5
	Fis5
	G5
	Gis5
	A5
	Ais5
	B5
	C6
	Cis6
	D6
	Dis6
	E6
	F6
	Fis6
	G6
	Gis6
	A6
	Ais6
	B6
	C7
	Cis7
	D7
	Dis7
	E7
	F7
	Fis7
	G7
	Gis7
	A7
	Ais7
	B7
	C8
	Cis8
	D8
	Dis8
	E8
	F8
	Fis8
	G8
	Gis8
	A8
	Ais8
	B8
	C9
	Cis9
	D9
	Dis9
	E9
	F9
	Fis9
	G9
	Gis9
	A9
	Ais9
	B9
)

var noteNames = map[Note]string{
	C0:   "C0",
	Cis0: "Cis0",
	D0:   "D0",
	Dis0: "Dis0",
	E0:   "E0",
	F0:   "F0",
	Fis0: "Fis0",
	G0:   "G0",
	Gis0: "Gis0",
	A0:   "A0",
	Ais0: "Ais0",
	B0:   "B0",
	C1:   "C1",
	Cis1: "Cis1",
	D1:   "D1",
	Dis1: "Dis1",
	E1:   "E1",
	F1:   "F1",
	Fis1: "Fis1",
	G1:   "G1",
	Gis1: "Gis1",
	A1:   "A1",
	Ais1: "Ais1",
	B1:   "B1",
	C2:   "C2",
	Cis2: "Cis2",
	D2:   "D2",
	Dis2: "Dis2",
	E2:   "E2",
	F2:   "F2",
	Fis2: "Fis2",
	G2:   "G2",
	Gis2: "Gis2",
	A2:   "A2",
	Ais2: "Ais2",
	B2:   "B2",
	C3:   "C3",
	Cis3: "Cis3",
	D3:   "D3",
	Dis3: "Dis3",
	E3:   "E3",
	F3:   "F3",
	Fis3: "Fis3",
	G3:   "G3",
	Gis3: "Gis3",
	A3:   "A3",
	Ais3: "Ais3",
	B3:   "B3",
	C4:   "C4",
	Cis4: "Cis4",
	D4:   "D4",
	Dis4: "Dis4",
	E4:   "E4",
	F4:   "F4",
	Fis4: "Fis4",
	G4:   "G4",
	Gis4: "Gis4",
	A4:   "A4",
	Ais4: "Ais4",
	B4:   "B4",
	C5:   "C5",
	Cis5: "Cis5",
	D5:   "D5",
	Dis5: "Dis5",
	E5:   "E5",
	F5:   "F5",
	Fis5: "Fis5",
	G5:   "G5",
	Gis5: "Gis5",
	A5:   "A5",
	Ais5: "Ais5",
	B5:   "B5",
	C6:   "C6",
	Cis6: "Cis6",
	D6:   "D6",
	Dis6: "Dis6",
	E6:   "E6",
	F6:   "F6",
	Fis6: "Fis6",
	G6:   "G6",
	Gis6: "Gis6",
	A6:   "A6",
	Ais6: "Ais6",
	B6:   "B6",
	C7:   "C7",
	Cis7: "Cis7",
	D7:   "D7",
	Dis7: "Dis7",
	E7:   "E7",
	F7:   "F7",
	Fis7: "Fis7",
	G7:   "G7",
	Gis7: "Gis7",
	A7:   "A7",
	Ais7: "Ais7",
	B7:   "B7",
	C8:   "C8",
	Cis8: "Cis8",
	D8:   "D8",
	Dis8: "Dis8",
	E8:   "E8",
	F8:   "F8",
	Fis8: "Fis8",
	G8:   "G8",
	Gis8: "Gis8",
	A8:   "A8",
	Ais8: "Ais8",
	B8:   "B8",
	C9:   "C9",
	Cis9: "Cis9",
	D9:   "D9",
	Dis9: "Dis9",
	E9:   "E9",
	F9:   "F9",
	Fis9: "Fis9",
	G9:   "G9",
	Gis9: "Gis9",
	A9:   "A9",
	Ais9: "Ais9",
	B9:   "B9",
}

/*
from supercollider/include/plugin_interface/SC_InlineUnaryOp.h

inline float32 sc_midicps(float32 note)
{
	return (float32)440. * std::pow((float32)2., (note - (float32)69.) * (float32)0.083333333333);
}
*/

func (n Note) Frequency() float64         { return MidiToFreq(float64(n)) }
func (n Note) Transpose(add float64) Note { return Note(float64(n) + add) }
func (n Note) Octave(num int) Note        { return n.Transpose(float64(num * 12)) }
func (n Note) Params() map[string]float64 { return map[string]float64{"freq": n.Frequency()} }
func (n Note) String() string             { return noteNames[n] }

const midiCpsFactor = 1/300.0 + 0.08

var midiToCps = map[float64]float64{}
var cpsToMidi = map[float64]float64{}
var midiToNote = map[float64]Note{}

func init() {
	/*
		prefill musical frequencies
	*/

	for note, _ := range noteNames {
		freq := note.Frequency()
		midiToCps[float64(note)] = freq
		cpsToMidi[freq] = float64(note)
		midiToNote[float64(note)] = note
	}
}

func FreqToNote(freq float64) Note {
	midi := FreqToMidi(freq)
	if midi < 12.0 || midi > 131.0 {
		panic(fmt.Sprintf("can't convert freq %.4f to Note", freq))
	}
	// fmt.Printf("freq: %.4f\n", freq)
	// midi = float64(int(midi))
	// return midiToNote[midi]
	return MidiToNote(midi)
}

func MidiToNote(midi float64) Note {
	if midi < 12.0 || midi > 131.0 {
		panic(fmt.Sprintf("can't convert midinote: %.4f to Note", midi))
	}
	// fmt.Printf("> midi: %.4f >== ", midi)
	midi = roundDec(midi, 0)
	// fmt.Printf("%.4f\n", midi)
	return midiToNote[midi]
}

// FreqToMidi returns the midi note for the given frequency
// It uses a cache for the most common midinotes 27.0 to 125.0
func FreqToMidi(freq float64) float64 {
	n, has := cpsToMidi[freq]
	if has {
		return n
	}
	return cpsMidi(freq)
}

// MidiToFreq returns the frequency for the given midi note
// It uses a cache for the most common midinotes 27.0 to 125.0
func MidiToFreq(note float64) float64 {
	f, has := midiToCps[note]
	if has {
		return f
	}
	return midiCps(note)
}

/*
inline float64 sc_log2(float64 x)
{
#ifdef HAVE_C99
return ::log2(std::abs(x));
#else
return std::log(std::abs(x)) * rlog2;
#endif
}

inline float32 sc_cpsmidi(float32 freq)
{
return sc_log2(freq * (float32)0.0022727272727) * (float32)12. + (float32)69.;
}
*/

/*
 */
func cpsMidi(freq float64) float64 {
	return math.Log2(math.Abs(freq/440.0))*12.0 + 69.0
}

/*
from supercollider/include/plugin_interface/SC_InlineUnaryOp.h

inline float32 sc_midicps(float32 note)
{
	return (float32)440. * std::pow((float32)2., (note - (float32)69.) * (float32)0.083333333333);
}
*/
func midiCps(note float64) float64 { return 440.0 * math.Pow(2.0, (note-69.0)*midiCpsFactor) }

// a1 = 440hz = 69 (midi note)
/*
  international Name     Midinote
	[C-1]                  0
	[A-1]                  9
	C0                     12
	A0                     21
	C1                     24
	A1                     33
	C2                     36
	A2                     45
	C3                     48
	A3                     57
	C4 (c')                60
	A4 (a')                69
	C5 (c'')               72
	A5 (a'')               81
	C6 (c''')              84
	A6 (a''')              93
	C7 (c'''')             96
	A7 (a'''')            105
	C8 (c''''')           108
	A8 (a''''')           117
*/
