package music

import "math/big"

// how long does a 256th play in milliseconds
type milliSecsOf256th struct {
	*big.Rat
}

// returns number of  Nanoseconds
func (t *milliSecsOf256th) MilliSecs(m Measure) uint {
	r := new(big.Rat)
	r.Mul(t.Rat, big.NewRat(int64(int(m)), 1))
	f, _ := r.Float64()
	// fmt.Printf("Nanoseconds: %0.2f\n", f*1000000)
	return uint(f * 1000000)
	/*
		a := big.NewRat(64, 60)
		b := big.NewRat(1, int64(uint(t)))
		z := new(big.Rat)
		z = z.Mul(a, b)
		z = z.Mul(z, big.NewRat(int64(int(m)), 1))
		f, _ := z.Float64()
		fmt.Printf("ticksPerMinute: %0.2f", f*1000)
		return uint(f * 1000)
	*/
}

type Tempo interface {
	// returns the milliseconds of the measure
	MilliSecs(m Measure) float64 //uint
	// return the beats per minute
	BPM() float64
}

type BPM float64

func (b BPM) BPM() float64 {
	return float64(b)
}

func (b BPM) MilliSecs(m Measure) float64 {
	// r := new(big.Rat)
	// r.SetFloat64(f)
	//ms := &milliSecsOf256th{big.NewRat(9375, int64(uint(b))*10)}
	// microseconds instead of millis
	//ms := &milliSecsOf256th{big.NewRat(9375, int64(uint(b))*10)}
	//return ms.MilliSecs(m)
	//return uint(RoundFloat(9375.0*float64(b)*100000, 0))

	return 9375.0 * float64(int(m)) * 100000.0 / float64(b)

	// uint(f * 1000000)

	// 937.5 / b
	// return ticksPerMinute(uint(b) * 64).MilliSecs(m)
}
