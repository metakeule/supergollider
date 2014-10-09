package music

import (
	"math/big"
	"sort"
	"strings"
)

type Measure int

const (
	_ = iota
	_ = 1 << iota
	_
	_
	SemiFusa   Measure = 1 << iota // element has default flags
	Fusa                           // element should not have an id attribute
	SemiMinima                     // element should not have a class attribute
	Minima
	SemiBrevis
	Brevis
	Longa
	Maxima
)

type measureRationale struct {
	*big.Rat
	Name string
}

var measureNames = map[Measure]*measureRationale{
	Measure(1): &measureRationale{Rat: big.NewRat(1, 256)},
	Measure(2): &measureRationale{Rat: big.NewRat(1, 128)},
	Measure(4): &measureRationale{Rat: big.NewRat(1, 64)},
	Measure(8): &measureRationale{Rat: big.NewRat(1, 32)},
	SemiFusa:   &measureRationale{Rat: big.NewRat(1, 16), Name: "SemiFusa"},
	Fusa:       &measureRationale{Rat: big.NewRat(1, 8), Name: "Fusa"},
	SemiMinima: &measureRationale{Rat: big.NewRat(1, 4), Name: "SemiMinima"},
	Minima:     &measureRationale{Rat: big.NewRat(1, 2), Name: "Minima"},
	SemiBrevis: &measureRationale{Rat: big.NewRat(1, 1), Name: "SemiBrevis"},
	Brevis:     &measureRationale{Rat: big.NewRat(2, 1), Name: "Brevis"},
	Longa:      &measureRationale{Rat: big.NewRat(4, 1), Name: "Longa"},
	Maxima:     &measureRationale{Rat: big.NewRat(8, 1), Name: "Maxima"},
}

func convertMeasure(in Measure) string {
	if in == Measure(0) {
		return "0"
	}
	measureNamesKeys := []int{}

	for m := range measureNames {
		measureNamesKeys = append(measureNamesKeys, int(m))
	}

	sort.Ints(measureNamesKeys)
	inInt := int(in)
	_ = inInt
	for i := len(measureNamesKeys) - 1; i >= 0; i-- {
		key := measureNamesKeys[i]

		if inInt > key {
			a := big.NewRat(int64(inInt), 256)
			b := measureNames[Measure(key)].Rat
			r := new(big.Rat)
			s := r.Sub(a, b)
			return b.RatString() + " + " + s.RatString()
		}
	}
	return "not found"
}

func (m Measure) Name() string {
	s, found := measureNames[m]

	if !found {
		return convertMeasure(m)
	}

	if s.Name == "" {
		return s.RatString()
	}

	return s.Name
}

func (m Measure) String() string {
	s, found := measureNames[m]

	if found {
		return s.RatString()
	}
	return convertMeasure(m)
}

// measure via string
//func m_(s string) Measure {
func m_(s string) float64 {
	r := new(big.Rat)
	_, ok := r.SetString(s)
	if !ok {
		panic("can't convert to measure")
	}

	x := new(big.Rat)

	// d := big.NewRat(1, 256)
	d := big.NewRat(256, 1)

	x.Mul(r, d)

	f, _ := x.Float64()

	//return Measure(int(f))
	return f

}

func _M(s string) float64 {
	summands := strings.Split(s, "+")

	var sum float64

	for _, s := range summands {
		sum += m_(strings.TrimSpace(s))
	}

	return sum
}

func M(s string) Measure {
	return Measure(int(_M(s)))
}

/*
func init() {
	fmt.Println(
		"SemiFusa/2", SemiFusa/2, "\n",
		"SemiFusa", SemiFusa.Name(), "\n",
		"SemiFusa * 2", SemiFusa*2, "\n",
		"SemiFusa * 3", SemiFusa*3, "\n",
		"SemiFusa * 4", SemiFusa*4, "\n",
		"SemiFusa * 5", SemiFusa*5, "\n",
		"SemiFusa * 6", SemiFusa*6, "\n",
		"SemiFusa * 7", SemiFusa*7, "\n",
		"SemiFusa * 8", SemiFusa*8, "\n",
		"SemiFusa * 9", SemiFusa*9, "\n",
		"SemiFusa * 10", SemiFusa*10, "\n",
		"SemiFusa * 11", SemiFusa*11, "\n",
		"SemiFusa * 12", SemiFusa*12, "\n",
		"Fusa", Fusa.Name(), Fusa, "\n",
		"SemiMinima", SemiMinima.Name(), SemiMinima, "\n",
		"Minima", Minima.Name(), Minima, "\n",
		"SemiBrevis", SemiBrevis.Name(), SemiBrevis, "\n",
		"Brevis", Brevis.Name(), Brevis, "\n",
		"Longa", Longa.Name(), Longa, "\n",
		"Maxima", Maxima.Name(), Maxima, "\n",
		`M("1/20")`, M("1/20"), "\n",
		`M("1/19")`, M("1/19"), "\n",
		`M("1/18")`, M("1/18"), "\n",
		`M("1/17")`, M("1/17"), "\n",
		`M("1/16")`, M("1/16"), "\n",
		`M("1/15")`, M("1/15"), "\n",
		`M("1/14")`, M("1/14"), "\n",
		`M("1/13")`, M("1/13"), "\n",
		`M("1/12")`, M("1/12"), "\n",
		`M("1/11")`, M("1/11"), "\n",
		`M("1/10")`, M("1/10"), "\n",
		`M("1/9")`, M("1/9"), "\n",
		`M("1/8")`, M("1/8"), "\n",
		`M("1/7")`, M("1/7"), "\n",
		`M("1/6")`, M("1/6"), "\n",
		`M("1/5")`, M("1/5"), "\n",
		`M("1/4")`, M("1/4"), "\n",
		`M("1/2")`, M("1/2"), "\n",
		`M("1/3")`, M("1/3"), "\n",
		`M("2/3")`, M("2/3"), "\n",
		`M("1")`, M("1"), "\n",
		`M("2.5")`, M("2.5"), "\n",
		`M("2")`, M("2"), "\n",
		`M("3")`, M("3"), "\n",
	)
}
*/

/*
nach der weißen mensuralnotation (imperfekte variante)

Maxima                   = achtfach ganze    = 2048
Longa                    = vierfach ganze    = 1024
Brevis                   = doppelganze       = 512
Semibrevis               = ganze Note        = 256
Minima                   = halbe Note        = 128
Semiminima               = Viertel Note      = 64
Fusa oder Chroma         = Achtel Note       = 32
Semifusa oder Semichroma = Sechzehntel Note  = 16
                         32tel               = 8
                         64tel               = 4
                         128tel              = 2
                         256tel              = 1

ein tempo wird in 256tel pro minute angegeben und die
auflösung in ticks pro minute daraus ergibt sich das jeweilige verhältnis von 256tel zu ticks

beim takt wird angegeben, aus wievielen 256tel er besteht (wie "lang" er ist, das Measure, Taktmaß)
also Semibrevis ist ein takt, der aus 4 vierteln oder 2/2 oder 8/8 usw besteht.
Minima * 3 ist ein takt, der aus 3 halben, oder 6 vierteln usw besteht.

die frage ist, ob der takt auch betonungspositionen haben soll (unklar) und slots (auch unklar)


Auch die Ziffern 3 und 4 sind als Proportionsbezeichnung möglich (Proportio tripla bzw. quadrupla). Generell gilt dabei: Die Ziffer zeigt die Zahl der Semibreves an, die einen Tactus bilden, das Tempuszeichen zeigt an, wie viele Tactus wiederum zu einer übergeordneten Einheit zusammengefasst werden.

Semibrevis die ganze Note, aus der Minima die Halbe usw. Zudem finden die Brevis als Doppelganze und seltener auch die Longa als Vierfachganze heute noch Verwendung (z. B. in langen Schlussakkorden).

*/

/*
   Tempo wird immer angegeben in SemiMinima / Minute



   1Min = BPM * (1/4)
   60000ms = BPM * (64/256)
   60000ms = BPM * 64 * Tick
   Tick = 60000ms / (BPM * 64)
   Tick = 937,5 / BPM




   T = 64 * (1/t) * (1/1Min.)
   t = 64 * (1/T) * (1/1Min.)
   t = (64/T*Min.)
   t = (64/T*60)
   t = (64/T)


   wenn 1 Tick einem 256tel entspricht,
   geben wir also die 256tel pro Minute an
*/

// Add starts at pos and returns the number of bars needed until m can be reached
// and the position with the last bar
func (b Measure) Add(m Measure) (numBars int, posLastBar Measure) {
	if m < b {
		return 0, m
	}

	if m == b {
		return 1, 0
	}

	numBars = int(m / b)
	posLastBar = m % b
	return
}
