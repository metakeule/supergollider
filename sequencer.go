package music

type sequencer struct {
	seq Sequencer
	v   *Voice
}

func (s *sequencer) NumBars() int {
	return s.seq.Next(s.v).NumBars()
}

func (s *sequencer) Events(barNum int, t Tracker) map[Measure][]*Event {
	return s.seq.Next(s.v).Events(barNum, t)
}

/*
func (s *sequencer) Pattern(t Tracker) {
	s.seq.Next(s.v).Pattern(t)
}
*/

type Sequencer interface {
	Next(v *Voice) Pattern
}

type sequencePairs struct {
	startPos        string
	positions       []string
	paramsSequences [][]Parameter
}

func newSequencePairs(startPos string, positions ...string) *sequencePairs {
	return &sequencePairs{positions: positions}
}

func (s *sequencePairs) AddSequence(paramsSeq ...Parameter) {
	s.paramsSequences = append(s.paramsSequences, paramsSeq)
}

/*

	bass[0].SequencePairs("0", "1/4", "1/8").
		AddSequence(note.C1, note.D1, note.E1).
		AddSequence(amp.FF, amp.PP,amp.PP,amp.PP)



*/

/*
type Sequencer struct {
	r *rhythm
	paramSequence []Parameter
}
*/

/*
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
*/
