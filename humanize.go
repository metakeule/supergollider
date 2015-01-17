package supergollider

import (
	"math/rand"
	"time"
)

// a simple idea to humanize offset by +/10 and amp by +- 0.1
type humanize_V1 struct {
	params       Parameter
	offsetFactor float64
	ampFactor    float64
	freqFactor   float64
	rateFactor   float64
}

func (h *humanize_V1) Params() map[string]float64 {
	p := h.params.Params()
	// between 0 and 1

	if h.offsetFactor > 0 {
		src1 := rand.NewSource(time.Now().UTC().UnixNano())
		r1 := rand.New(src1).Float64()

		offsetAdd := r1 * h.offsetFactor * (-1)
		p["offset"] = p["offset"] + offsetAdd
	}

	if h.ampFactor > 0 {
		src2 := rand.NewSource(time.Now().UTC().UnixNano() * time.Now().UTC().UnixNano())
		r2 := rand.New(src2).Float64()

		ampAdd := r2 * h.ampFactor

		amp, hasAmp := p["amp"]
		if !hasAmp {
			amp = 1.0
		}
		p["amp"] = amp + ampAdd
	}

	/*
		TODO also for freq
		offsetAdd := (r1 - 0.5) * h.offsetFactor
			if r1 < 0.5 {
				offsetAdd = r1 * h.offsetFactor * (-1)
			}
			p["offset"] = p["offset"] + offsetAdd
	*/

	if h.freqFactor > 0 && p["freq"] != 0 {
		// was := p["freq"]
		src3 := rand.NewSource(time.Now().UTC().UnixNano() * time.Now().UTC().UnixNano())
		r3 := rand.New(src3).Float64()

		freqAdd := (r3 - 0.5) * h.freqFactor

		if r3 <= 0.5 {
			freqAdd = r3 * h.freqFactor * (-1)
		}

		if x := p["freq"] + freqAdd; x > 0 {
			p["freq"] = x
			// fmt.Printf("freq: %v => %v \n", was, p["freq"])
		}

	}

	if h.rateFactor > 0 {
		src4 := rand.NewSource(time.Now().UTC().UnixNano() * time.Now().UTC().UnixNano())
		r4 := rand.New(src4).Float64()

		rateAdd := r4 * h.rateFactor

		rate, hasRate := p["rate"]
		if !hasRate {
			rate = 1.0
		}
		// fmt.Printf("rate was : %v\n", rate)

		// fmt.Printf("rate: %v\n", rate+rateAdd)
		p["rate"] = rate + rateAdd
	}

	return p
}

type HumanizeV1 struct {
	OffsetFactor float64
	AmpFactor    float64
	FreqFactor   float64
	RateFactor   float64
}

func (h HumanizeV1) Modify(params Parameter) Parameter {
	return &humanize_V1{params, h.OffsetFactor, h.AmpFactor, h.FreqFactor, h.RateFactor}
}
