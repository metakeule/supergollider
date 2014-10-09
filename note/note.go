package note

import "math"

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

/*
from supercollider/include/plugin_interface/SC_InlineUnaryOp.h

inline float32 sc_midicps(float32 note)
{
	return (float32)440. * std::pow((float32)2., (note - (float32)69.) * (float32)0.083333333333);
}
*/

func (n Note) Frequency() float64         { return MidiCps(float64(n)) }
func (n Note) Transpose(add float64) Note { return Note(float64(n) + add) }
func (n Note) Octave(num int) Note        { return n.Transpose(float64(num * 12)) }
func (n Note) Params() map[string]float64 { return map[string]float64{"freq": n.Frequency()} }

const midiCpsFactor = 1/300.0 + 0.08

func MidiCps(note float64) float64 { return 440.0 * math.Pow(2.0, (note-69.0)*midiCpsFactor) }

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
