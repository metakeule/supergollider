// Package amp provides shortcuts for amplitude measures contrasting with the dyn package that affects how
// something is played
package amp

import "github.com/metakeule/music"

type amp float64

func (d amp) Params() map[string]float64 { return music.Amp(float64(d)).Params() }

var (
	FFFF amp = 0.5
	FFF  amp = 0.45
	FF   amp = 0.4
	F    amp = 0.35
	MF   amp = 0.3
	MP   amp = 0.25
	P    amp = 0.2
	PP   amp = 0.15
	PPP  amp = 0.1
	PPPP amp = 0.05
)
