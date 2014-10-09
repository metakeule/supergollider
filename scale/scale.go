package scale

import (
	"github.com/metakeule/music"
	"github.com/metakeule/music/note"
)

type Scale interface {
	Degree(degree int) music.Parameter
	Base() note.Note
}

type Chromatic struct {
	BaseNote note.Note
}

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

func (s *Chromatic) Degree(scalePosition int) music.Parameter {
	return note.Note(float64(s.BaseNote) + float64(scalePosition))
	//return note.MidiCps(float64(s.BaseNote) + float64(scalePosition))
}

func (s *Chromatic) Base() note.Note {
	return s.BaseNote
}

type Periodic struct {
	Steps             []uint // each steps in factor of chromatic steps, begins with the second step (first is basetone)
	NumChromaticSteps uint   // the number of chromatic steps that correspond to one scale periodicy
	BaseNote          note.Note
}

func (s *Periodic) Base() note.Note {
	return s.BaseNote
}

// TODO: test it
func (s *Periodic) Degree(scalePosition int) music.Parameter {
	// we need to calculate the position in terms of the chromatic scale
	// and then we return the frequency via MidiCps
	num := len(s.Steps)

	posInScale := scalePosition % num
	cycle := scalePosition / num

	temp := int(s.BaseNote) + (cycle * int(s.NumChromaticSteps))

	if posInScale == 0 {
		return note.Note(float64(temp))
		//return note.MidiCps(float64(temp))
	}

	if posInScale < 0 {
		for i := posInScale; i < 0; i++ {
			temp -= int(s.Steps[num+i])
		}
		return note.Note(float64(temp))
		// return note.MidiCps(float64(temp))
	}

	for i := 0; i < posInScale; i++ {
		temp += int(s.Steps[i])
	}
	return note.Note(float64(temp)) //  note.MidiCps(float64(temp))
}

func Ionian(base note.Note) *Periodic {
	return &Periodic{
		Steps:             []uint{2, 2, 1, 2, 2, 2, 1},
		NumChromaticSteps: 12,
		BaseNote:          base,
	}
}

func Dorian(base note.Note) *Periodic {
	return &Periodic{
		Steps:             []uint{2, 1, 2, 2, 2, 1, 2},
		NumChromaticSteps: 12,
		BaseNote:          base,
	}
}

func Phrygian(base note.Note) *Periodic {
	return &Periodic{
		Steps:             []uint{1, 2, 2, 2, 1, 2, 2},
		NumChromaticSteps: 12,
		BaseNote:          base,
	}
}

func Lydian(base note.Note) *Periodic {
	return &Periodic{
		Steps:             []uint{2, 2, 2, 1, 2, 2, 1},
		NumChromaticSteps: 12,
		BaseNote:          base,
	}
}

func Mixolydian(base note.Note) *Periodic {
	return &Periodic{
		Steps:             []uint{2, 2, 1, 2, 2, 1, 2},
		NumChromaticSteps: 12,
		BaseNote:          base,
	}
}

func Aeolian(base note.Note) *Periodic {
	return &Periodic{
		Steps:             []uint{2, 1, 2, 2, 1, 2, 2},
		NumChromaticSteps: 12,
		BaseNote:          base,
	}
}

func Locrian(base note.Note) *Periodic {
	return &Periodic{
		Steps:             []uint{1, 2, 2, 1, 2, 2, 2},
		NumChromaticSteps: 12,
		BaseNote:          base,
	}
}

func Hypolydian(base note.Note) *Periodic {
	return Ionian(base)
}

func Hypomixolydian(base note.Note) *Periodic {
	return Dorian(base)
}

func Dur(base note.Note) *Periodic {
	return Ionian(base)
}

func Major(base note.Note) *Periodic {
	return Dur(base)
}

func Moll(base note.Note) *Periodic {
	return Aeolian(base)
}

func Minor(base note.Note) *Periodic {
	return Moll(base)
}

func Hypodorian(base note.Note) *Periodic {
	return Aeolian(base)
}

func Hypophrygian(base note.Note) *Periodic {
	return Locrian(base)
}

var Mood = map[string]func(base note.Note) *Periodic{
	"serious":  Dorian,
	"sad":      Hypodorian,
	"vehement": Phrygian,
	"tender":   Hypophrygian,
	"happy":    Lydian,
	"pious":    Hypolydian,
	"youthful": Mixolydian,
}

type stepper struct {
	Scale   Scale
	Step    int
	Degrees []int
}

func NewStepper(scale Scale, degrees ...int) *stepper {
	if len(degrees) < 1 {
		panic("must have at least 1 degree")
	}
	return &stepper{Scale: scale, Degrees: degrees}
}

func (s *Periodic) Stepper(degrees ...int) *stepper {
	return NewStepper(s, degrees...)
}

func (s *Chromatic) Stepper(degrees ...int) *stepper {
	return NewStepper(s, degrees...)
}

func (s *stepper) Params() map[string]float64 {
	degr := s.Degrees[s.Step]
	if s.Step < len(s.Degrees)-1 {
		s.Step++
	} else {
		s.Step = 0
	}
	return s.Scale.Degree(degr).Params()
}

type seq struct {
	Scales  []Scale
	current int
}

func (s *seq) Degree(degree int) music.Parameter {
	return s.Scales[s.current].Degree(degree)
}

func (s *seq) Base() note.Note {
	return s.Scales[s.current].Base()
}

func (s *seq) Next(pos string) music.Pattern {
	return music.Exec(pos, func() {
		if len(s.Scales)-1 <= s.current {
			s.current = 0
			return
		}
		s.current++
	})
}

func (s *seq) Set(pos string, index int) music.Pattern {
	if index < 0 {
		panic("index < 0 not allowed")
	}
	if len(s.Scales)-1 <= index {
		panic("index larger than len(Scales)-1")
	}
	return music.Exec(pos, func() { s.current = index })
}

func Seq(scales ...Scale) *seq {
	return &seq{Scales: scales}
}

/*
func (s *stepper) SetScale(scale Scale) *stepper {
	return &stepper{
		scale:   scale,
		degrees: s.degrees,
		step:    s.step,
	}
}
*/
/*
func (s *stepper) SetScale(scale Scale) {
	s.scale = scale
}
*/

/*
func (s *stepper) SetDegrees(degrees ...int) *stepper {
	if len(degrees) < 1 {
		panic("must have at least 1 degree")
	}
	return &stepper{
		scale:   s.scale,
		degrees: degrees,
		step:    0,
	}
}
*/

/*
func (s *stepper) SetDegrees(degrees ...int) {
	if len(degrees) < 1 {
		panic("must have at least 1 degree")
	}
	s.step = 0
	s.degrees = degrees
}
*/

/*
func (r *rhythm) pos() (pos_ string) {
	pos_ = r.positions[r.positionsIndex]
	if r.positionsIndex < len(r.positions)-1 {
		r.positionsIndex++
	} else {
		r.positionsIndex = 0
	}
	return pos_
}

*/
