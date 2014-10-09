package music

import (
	"fmt"
	"math/rand"
	"time"
)

/*
type Pattern interface {
	Pattern(Tracker)
}
*/

type Pattern interface {
	//Events(barNum int, barMeasure Measure) map[Measure][]*Event
	Events(barNum int, t Tracker) map[Measure][]*Event
	NumBars() int
}

//type PatternFunc func(Tracker)

type PatternFunc func(barNum int, t Tracker) map[Measure][]*Event

func (tf PatternFunc) Events(barNum int, t Tracker) map[Measure][]*Event {
	return tf(barNum, t)
}

func (tf PatternFunc) NumBars() int {
	return 1
}

/*
func (tf PatternFunc) Pattern(tr Tracker) {
	tf(tr)
}
*/

/*
type seqModTrafo struct {
	*seqPlay
	pos            Measure
	overrideParams Parameter
	params         Parameter
}

func (sm *seqModTrafo) Pattern(tr *Track) {
	tr.At(sm.pos, ChangeEvent(sm.seqPlay.v, Params(sm.params, sm.overrideParams)))
}

type seqPlay struct {
	seq        []Parameter
	initParams Parameter
	v          *Voice
	Pos        int
}

func (sp *seqPlay) Modify(pos string, params ...Parameter) Pattern {
	params_ := sp.seq[sp.Pos]
	if sp.Pos < len(sp.seq)-1 {
		sp.Pos++
	} else {
		sp.Pos = 0
	}
	return &seqModTrafo{seqPlay: sp, pos: M(pos), overrideParams: Params(params...), params: params_}
}

func (sp *seqPlay) PlayDur(pos, dur string, params ...Parameter) Pattern {
	params_ := sp.seq[sp.Pos]
	if sp.Pos < len(sp.seq)-1 {
		sp.Pos++
	} else {
		sp.Pos = 0
	}
	return &seqPlayTrafo{seqPlay: sp, pos: M(pos), dur: M(dur), overrideParams: Params(params...), params: params_}
}
*/

/*
func ParamSequence(v *Voice, initParams Parameter, paramSeq ...Parameter) *seqPlay {
	return &seqPlay{
		initParams: initParams,
		seq:        paramSeq,
		v:          v,
	}
}
*/

/*
type seqPlayTrafo struct {
	*seqPlay
	pos            Measure
	dur            Measure
	overrideParams Parameter
	params         Parameter
}
*/

/*
func (spt *seqPlayTrafo) Params() (p map[string]float64) {
	return Params(spt.seqPlay.initParams, spt.seqPlay.seq[spt.seqPlay.Pos], spt.overrideParams).Params()
}
*/

/*
func (spt *seqPlayTrafo) Pattern(tr *Track) {
	params := Params(spt.seqPlay.initParams, spt.params, spt.overrideParams)
	tr.At(spt.pos, OnEvent(spt.seqPlay.v, params))
	tr.At(spt.pos+spt.dur, OffEvent(spt.seqPlay.v))
*/
/*
	if spt.seqPlay.Pos < len(spt.seqPlay.seq)-1 {
		spt.seqPlay.Pos++
	} else {
		spt.seqPlay.Pos = 0
	}
*/
/*
}
*/

// type end struct{}

func (e End) Pattern(t Tracker) {
	t.At(M(string(e)), fin)
}

func (e End) Events(barNum int, t Tracker) map[Measure][]*Event {
	return map[Measure][]*Event{
		M(string(e)): []*Event{fin},
	}
}

func (e End) NumBars() int {
	return 1
}

// var End = end{}
type End string

type Start string

func (s Start) Events(barNum int, t Tracker) map[Measure][]*Event {
	return map[Measure][]*Event{
		M(string(s)): []*Event{start},
	}
}

func (s Start) NumBars() int {
	return 1
}

// var Start = begin{}
/*
type stopAll struct {
	pos    Measure
	Voices []*Voice
}

func StopAll(pos string, vs ...[]*Voice) *stopAll {
	s := &stopAll{pos: M(pos)}

	for _, v := range vs {
		s.Voices = append(s.Voices, v...)
	}

	return s
}

func (p *stopAll) Pattern(t Tracker) {
	for i := 0; i < len(p.Voices); i++ {
		t.At(p.pos, OffEvent(p.Voices[i]))
	}
}
*/

/*
type setTempo struct {
	Tempo Tempo
	Pos   Measure
}

func SetTempo(at string, t Tempo) *setTempo {
	return &setTempo{t, M(at)}
}

// Tracker must be a *Track
func (s *setTempo) Pattern(t Tracker) {
	t.(*Track).SetTempo(s.Pos, s.Tempo)
}
*/

/*
type times struct {
	times int
	trafo Pattern
}

func (n *times) NumBars() int {
	return n.times
}

func (n *times) Events(barNum int, barMeasure Measure) map[Measure][]*Event {
	res := map[Measure][]*Event{}

	for i := 0; i < n.times; i++ {
		for pos, events := range n.trafo.Events(barNum, barMeasure) {
			finalPos := M(fmt.Sprintf("%d", i)) + pos
			res[finalPos] = append(res[finalPos], events...)
		}
	}
	return res
}
*/

/*
func (n *times) Pattern(t Tracker) {
	for i := 0; i < n.times; i++ {
		n.trafo.Pattern(t)
	}
}
*/

/*

func Times(num int, trafo Pattern) Pattern {
	return &times{times: num, trafo: trafo}
}
*/

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type randomPattern []Pattern

func (r randomPattern) Events(barNum int, t Tracker) map[Measure][]*Event {
	return r[rand.Intn(len(r))].Events(barNum, t)
}

func (r randomPattern) NumBars() int {
	return r[0].NumBars()
}

/*
func (r randomPattern) Pattern(t Tracker) {
	r[rand.Intn(len(r))].Pattern(t)
}
*/

// each pattern must have the same NumBars
func RandomPattern(patterns ...Pattern) Pattern {
	if len(patterns) < 1 {
		panic("at least one pattern required")
	}
	numBars := patterns[0].NumBars()

	if len(patterns) == 1 {
		return randomPattern(patterns)
	}

	for i, p := range patterns[1:] {
		if p.NumBars() != numBars {
			panic(fmt.Sprintf("pattern %d does not have the same number of bars as pattern 0", i))
		}
	}

	return randomPattern(patterns)
}

type tempoSpan struct {
	current  float64
	step     float64
	modifier func(current, step float64) float64
	// t        *Track
}

type tempoSpanTrafo struct {
	*tempoSpan
	pos string
}

func (ts *tempoSpan) SetTempo(pos string) Pattern {
	return &tempoSpanTrafo{ts, pos}
}

func (ts *tempoSpanTrafo) NumBars() int {
	return 1
}

// TODO: fix it! we need to have a TempoSpan event that is handled specially by the  Tracker - OR - a
// special Event type that has a method that receives a *Track and may act upon it
// that must be somehow integrated to the tracker methods
func (ts *tempoSpanTrafo) Events(barNum int, t Tracker) (res map[Measure][]*Event) {
	var newtempo float64
	if ts.current == -1 {
		//newtempo = ts.modifier(ts.t.TempoAt(M(ts.pos)).BPM(), ts.step)
		newtempo = ts.modifier(t.TempoAt(M(ts.pos)).BPM(), ts.step)
	} else {
		ts.current = ts.modifier(ts.current, ts.step)
		newtempo = ts.current
	}
	//rounded := RoundFloat(newtempo, 4)
	//ts.t.SetTempo(M(ts.pos), BPM(newtempo))
	t.SetTempo(M(ts.pos), BPM(newtempo))
	return nil
}

/*
// Tracker must be a *Track
func (ts *tempoSpanTrafo) Pattern(t Tracker) {
	var newtempo float64
	if ts.current == -1 {
		newtempo = ts.modifier(t.(*Track).TempoAt(M(ts.pos)).BPM(), ts.step)
	} else {
		ts.current = ts.modifier(ts.current, ts.step)
		newtempo = ts.current
	}
	//rounded := RoundFloat(newtempo, 4)
	t.(*Track).SetTempo(M(ts.pos), BPM(newtempo))
}
*/

func StepAdd(current, step float64) float64 {
	return current + step
}

func StepMultiply(current, step float64) float64 {
	return current * step
}

// for start = -1 takes the current tempo
func SeqTempo(start float64, step float64, modifier func(current, step float64) float64) *tempoSpan {
	return &tempoSpan{current: start, step: step, modifier: modifier}
}

type seqBool struct {
	seq   []bool
	pos   int
	trafo Pattern
}

func (s *seqBool) Events(barNum int, t Tracker) (res map[Measure][]*Event) {
	if s.trafo == nil {
		return nil
	}
	if s.seq[s.pos] {
		res = s.trafo.Events(barNum, t)
	}
	if s.pos < len(s.seq)-1 {
		s.pos++
	} else {
		s.pos = 0
	}
	return
}

func (s *seqBool) NumBars() int {
	if s.trafo == nil {
		return 1
	}
	return s.trafo.NumBars()
}

func SeqSwitch(trafo Pattern, seq ...bool) Pattern {
	return &seqBool{seq: seq, pos: 0, trafo: trafo}
}

type sequence []Pattern

func (s sequence) NumBars() int {
	num := 0

	for _, p := range s {
		num += p.NumBars()
	}
	return num
}

func (s sequence) Events(barNum int, t Tracker) (res map[Measure][]*Event) {
	num := 0

	for _, p := range s {
		next := num + p.NumBars()
		if barNum < next {
			return p.Events(barNum-num, t)
		}
		num = next
	}
	return
}

func SeqPatterns(seq ...Pattern) Pattern {
	return sequence(seq)
}

type compose []Pattern

func (c compose) Events(barNum int, t Tracker) map[Measure][]*Event {
	res := map[Measure][]*Event{}
	for _, pattern := range c {
		if pattern != nil {
			for pos, events := range pattern.Events(barNum, t) {
				res[pos] = append(res[pos], events...)
			}
		}
	}
	return res
}

func (c compose) NumBars() int {
	max := 1
	for _, pattern := range c {
		if pattern != nil {
			if pattern.NumBars() > max {
				max = pattern.NumBars()
			}
		}
	}
	return max
}

func MixPatterns(trafos ...Pattern) Pattern {
	return compose(trafos)
}

type linearDistribute struct {
	from, to float64
	steps    int
	dur      Measure
	key      string
	// from, to float64, steps int, dur Measure
}

// LinearDistribution creates a transformer that modifies the given parameter param
// from the value from to the value to in n steps in linear growth for a total duration dur
func LinearDistribution(param string, from, to float64, n int, dur Measure) *linearDistribute {
	return &linearDistribute{from, to, n, dur, param}
}

func (l *linearDistribute) ModifyDistributed(position string, v *Voice) Pattern {
	//return &linearDistributeTrafo{l, v, M(position)}
	p := []Pattern{}
	width, diff := LinearDistributedValues(l.from, l.to, l.steps, l.dur)
	// tr.At(ld.pos, Change(ld.v, ))
	pos := M(position)
	val := l.from
	for i := 0; i < l.steps; i++ {
		// println(pos.String())
		p = append(p, &mod{pos, v, ParamsMap(map[string]float64{l.key: val})})
		// tr.At(pos, ChangeEvent(ld.v, ParamsMap(map[string]float64{ld.linearDistribute.key: val})))
		pos += width
		val += diff
	}

	return MixPatterns(p...)
}

type linearDistributeTrafo struct {
	*linearDistribute
	v   *Voice
	pos Measure
}

/*
func (ld *linearDistributeTrafo) Pattern(tr *Track) {
	width, diff := LinearDistributedValues(ld.linearDistribute.from, ld.linearDistribute.to, ld.linearDistribute.steps, ld.linearDistribute.dur)
	// tr.At(ld.pos, Change(ld.v, ))
	pos := ld.pos
	val := ld.linearDistribute.from
	for i := 0; i < ld.linearDistribute.steps; i++ {
		tr.At(pos, ChangeEvent(ld.v, ParamsMap(map[string]float64{ld.linearDistribute.key: val})))
		pos += width
		val += diff
	}
}
*/

// ---------------------------------------

// func ExponentialDistributedValues(from, to float64, steps int, dur Measure) (width Measure, diffs []float64) {

type expDistribute struct {
	from, to float64
	steps    int
	dur      Measure
	key      string
	// from, to float64, steps int, dur Measure
}

// ExponentialDistribution creates a transformer that modifies the given parameter param
// from the value from to the value to in n steps in exponential growth for a total duration dur
func ExponentialDistribution(param string, from, to float64, n int, dur Measure) *expDistribute {
	return &expDistribute{from, to, n, dur, param}
}

func (l *expDistribute) ModifyDistributed(position string, v *Voice) Pattern {
	p := []Pattern{}
	width, diffs := ExponentialDistributedValues(l.from, l.to, l.steps, l.dur)
	// tr.At(ld.pos, Change(ld.v, ))
	pos := M(position)
	for i := 0; i < l.steps; i++ {
		p = append(p, &mod{pos, v, ParamsMap(map[string]float64{l.key: diffs[i]})})
		//tr.At(pos, ChangeEvent(ld.v, ParamsMap(map[string]float64{ld.expDistribute.key: diffs[i]})))
		pos += width
		//val += diff
	}
	//return &expDistributeTrafo{l, v, M(position)}
	return MixPatterns(p...)
}

type expDistributeTrafo struct {
	*expDistribute
	v   *Voice
	pos Measure
}

/*
func (ld *expDistributeTrafo) Pattern(tr *Track) {
	width, diffs := ExponentialDistributedValues(ld.expDistribute.from, ld.expDistribute.to, ld.expDistribute.steps, ld.expDistribute.dur)
	// tr.At(ld.pos, Change(ld.v, ))
	pos := ld.pos
	for i := 0; i < ld.expDistribute.steps; i++ {
		tr.At(pos, ChangeEvent(ld.v, ParamsMap(map[string]float64{ld.expDistribute.key: diffs[i]})))
		pos += width
		//val += diff
	}
}
*/
