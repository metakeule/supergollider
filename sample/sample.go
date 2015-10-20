// Package sample provides shortcuts for parameters that are useful for samples
package sample

type Release float64

func (r Release) Params() map[string]float64 {
	return map[string]float64{
		"release": float64(r),
	}
}

type Attack float64

func (r Attack) Params() map[string]float64 {
	return map[string]float64{
		"attack": float64(r),
	}
}

type Skip float64

func (r Skip) Params() map[string]float64 {
	return map[string]float64{
		"skip": float64(r),
	}
}

type Rate float64

func (r Rate) Params() map[string]float64 {
	return map[string]float64{
		"rate": float64(r),
	}
}

type Backwards float64

func (r Backwards) Params() map[string]float64 {
	return map[string]float64{
		"rate": (-1.0) * float64(r),
	}
}

var Reverse = Rate(-1.0)
