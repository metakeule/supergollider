// Package dyn provides shortcuts for dynamic measures as a matter of how
// something is played contrasting to amp that affects the amplitude
package dyn

import "github.com/metakeule/music"

type dyn struct {
	name  string
	value float64
}

func (d dyn) String() string             { return d.name }
func (d dyn) Value() float64             { return d.value }
func (d dyn) Params() map[string]float64 { return music.Dyn(d.value).Params() }

var (
	FortissimoForte = dyn{"FortissimoForte", 3}
	Fortissimo      = dyn{"Fortissimo", 2}
	Forte           = dyn{"Forte", 1}
	Mezzoforte      = dyn{"Mezzoforte", 0}
	Mezzopiano      = dyn{"Mezzopiano", -1}
	Piano           = dyn{"Piano", -2}
	Pianissimo      = dyn{"Pianissimo", -3}
	PianissimoPiano = dyn{"PianissimoPiano", -4}

	Martellato = dyn{"Martellato", 10}
	Marcato    = dyn{"Accent", 11}
	Accent     = dyn{"Accent", 11} // same as Marcato

	Sforzando      = dyn{"Sforzando", 20}
	SforzandoPiano = dyn{"SforzandoPiano", 21}
	Sforzatissimo  = dyn{"Sforzatissimo", 22}
	SempreSforzato = dyn{"SempreSforzato", 23}

	FortePiano      = dyn{"FortePiano", -70}
	FortissimoPiano = dyn{"FortissimoPiano", -71}
	MezzofortePiano = dyn{"MezzofortePiano", -72}
	PianoForte      = dyn{"PianoForte", -73}

	Crescendo            = dyn{"Crescendo", -100}
	Decrescendo          = dyn{"Decrescendo", -101}
	CrescendoDecrescendo = dyn{"CrescendoDecrescendo", -102}
	DecrescendoCrescendo = dyn{"DecrescendoCrescendo", -103}

// Staccato
//	staccatissimo
// sustained
)
