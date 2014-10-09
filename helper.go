package music

import "strconv"

// stolen from https://groups.google.com/forum/#!topic/golang-nuts/ITZV08gAugI
func RoundFloat(x float64, prec int) float64 {
	frep := strconv.FormatFloat(x, 'g', prec, 64)
	f, _ := strconv.ParseFloat(frep, 64)
	return f
}

type tempoSorted []tempoAt

func (t tempoSorted) Len() int      { return len(t) }
func (t tempoSorted) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t tempoSorted) Less(i, j int) bool {
	return t[i].AbsPos < t[j].AbsPos
}

type eventsSorted []*Event

func (t eventsSorted) Len() int      { return len(t) }
func (t eventsSorted) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t eventsSorted) Less(i, j int) bool {
	return t[i].absPosition < t[j].absPosition
}

func MillisecsToTick(ms float64) int {
	return int(RoundFloat(ms*1000000.0, 0))
}

func TickToSeconds(tick int) float32 {
	return float32(tick) / float32(1000000000)
}
