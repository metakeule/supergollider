package music

import (
	"sort"
	"testing"
)

func TestLoop(t *testing.T) {
	st := New()
	instr := st.Instrument("test", "", 1)[0]

	l := Loop(
		instr.Play("0", Freq(220)),
		instr.Play("1/4", Freq(440)),
	)

	events := st.Track("4/4", BPM(120)).SetLoop("x", "0", l).Start().Fill(2).Events

	es := eventsSorted(events)
	sort.Sort(es)
	_ = events

	if len(es) != 6 {
		t.Fatalf("wrong number of events; is %d, must be 6", len(es))
	}

	type result struct {
		position  Measure
		frequency float64
	}

	results := []result{
		{M("0"), 220},
		{M("1/4"), 440},
		{M("1"), 220},
		{M("1 + 1/4"), 440},
		{M("2"), 220},
		{M("2 + 1/4"), 440},
	}

	for i, res := range results {
		if es[i].absPosition != res.position {
			t.Errorf("wrong position at index %d: got %s instead of %s", i, es[i].absPosition, res.position)
		}

		if fr := es[i].Params.Params()["freq"]; fr != res.frequency {
			t.Errorf("wrong frequency at index %d: got %v instead of %v", i, fr, res.frequency)
		}
	}
}

func TestLoopOfLoops(t *testing.T) {
	st := New()
	instr := st.Instrument("test", "", 1)[0]

	l1 := Loop(
		instr.Play("0", Freq(220)),
		instr.Play("1/4", Freq(440)),
	)

	l2 := Loop(
		instr.Play("0", Freq(550)),
		instr.Play("1/4", Freq(880)),
	)

	l := SeqLoop(SeqLoop(l1, l1), l2)

	events := st.Track("4/4", BPM(120)).SetLoop("x", "0", l).Start().Fill(3).Events

	es := eventsSorted(events)
	sort.Sort(es)
	_ = events

	type result struct {
		position  Measure
		frequency float64
	}

	results := []result{
		{M("0"), 220},
		{M("1/4"), 440},
		{M("1"), 220},
		{M("1 + 1/4"), 440},
		{M("2"), 550},
		{M("2 + 1/4"), 880},
		{M("3"), 220},
		{M("3 + 1/4"), 440},
	}

	if len(es) != len(results) {
		t.Fatalf("wrong number of events; is %d, must be %d", len(es), len(results))
	}

	for i, res := range results {
		if es[i].absPosition != res.position {
			t.Errorf("wrong position at index %d: got %s instead of %s", i, es[i].absPosition, res.position)
		}

		if fr := es[i].Params.Params()["freq"]; fr != res.frequency {
			t.Errorf("wrong frequency at index %d: got %v instead of %v", i, fr, res.frequency)
		}
	}
}

func TestLoopWithBarChange(t *testing.T) {
	st := New()
	instr := st.Instrument("test", "", 1)[0]

	l := Loop(
		instr.Play("0", Freq(220)),
		instr.Play("1/4", Freq(440)),
	).Next(
		instr.Play("1/8", Freq(330)),
		instr.Play("2/8", Freq(660)),
	)

	events := st.Track("4/4", BPM(120)).SetLoop("x", "0", l).Start().Fill(2).Events

	es := eventsSorted(events)
	sort.Sort(es)
	_ = events

	if len(es) != 6 {
		t.Fatalf("wrong number of events; is %d, must be 6", len(es))
	}

	/*
		for _, ev := range es {
			fmt.Println(ev.absPosition, ev.Params)
		}
	*/

	type result struct {
		position  Measure
		frequency float64
	}

	results := []result{
		{M("0"), 220},
		{M("1/4"), 440},
		{M("1 + 1/8"), 330},
		{M("1 + 1/4"), 660},
		{M("2"), 220},
		{M("2 + 1/4"), 440},
	}

	for i, res := range results {
		if es[i].absPosition != res.position {
			t.Errorf("wrong position at index %d: got %s instead of %s", i, es[i].absPosition, res.position)
		}

		if fr := es[i].Params.Params()["freq"]; fr != res.frequency {
			t.Errorf("wrong frequency at index %d: got %v instead of %v", i, fr, res.frequency)
		}
	}
}
