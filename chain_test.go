package supergollider

import (
	"github.com/metakeule/supergollider/note"

	"testing"
)

func TestChain(t *testing.T) {
	tr := newTrack(BPM(120), M("1/4")*4)
	// fmt.Println(tr.Tempi[0].Tempo.MilliSecs(M("1/4")))

	tr.At(M("0"), OnEvent(nil, note.A2))
	tr.At(M("1/4"), OffEvent(nil))

	tr.compile()

	if tr.Events[1].tick != 500000000 {
		t.Errorf("expecting 500000000, got %v", tr.Events[1].tick)
	}

	tr.At(M("0"), BPM(240).Event())
	// tr.SetTempo(M("0"), BPM(240))
	tr.compile()

	if tr.Events[1].tick != 250000000 {
		t.Errorf("expecting 250000000, got %v", tr.Events[1].tick)
	}
}

func TestAtSeq(t *testing.T) {
	var v1 *Voice
	tr := newTrack(BPM(120), M("1/4")*4)
	tr.Start()
	tr.At(M("0"), OnEvent(v1, note.A2))
	tr.At(M("1/4"), OffEvent(v1))
	// tr.At(M("2/4"), SetBPM(240))
	tr.At(M("2/4"), BPM(240).Event())
	// tr.SetTempo(M("2/4"), BPM(240))
	tr.At(M("2/4"), OnEvent(v1, note.A2))
	tr.At(M("3/4"), OffEvent(v1))
	tr.nextBar()
	tr.At(M("0"), OnEvent(v1, note.A2))
	tr.At(M("1/4"), OffEvent(v1))

	tr.compile()

	if tr.Events[0].tick != 0 {
		t.Errorf("expecting 0, got %v", tr.Events[0].tick)
	}

	if tr.Events[1].tick != 500000000 {
		t.Errorf("expecting 500000000, got %v", tr.Events[1].tick)
	}

	if tr.Events[3].tick != 1000000000 {
		t.Errorf("expecting 1000000000, got %v", tr.Events[3].tick)
	}

	if tr.Events[4].tick != 1250000000 {
		t.Errorf("expecting 1250000000, got %v", tr.Events[4].tick)
	}

	if tr.Events[5].tick != 1500000000 {
		t.Errorf("expecting 1500000000, got %v", tr.Events[5].tick)
	}

	if tr.Events[6].tick != 1750000000 {
		t.Errorf("expecting 1750000000, got %v", tr.Events[6].tick)
	}
}

func TestBarChange(t *testing.T) {
	var v1 *Voice
	tr := newTrack(BPM(120), M("4/4"))
	tr.Start()
	tr.At(M("0"), OnEvent(v1, note.A2))
	tr.At(M("1/4"), OffEvent(v1))
	tr.changeBar(M("3/4"))
	tr.At(M("1/4"), OnEvent(v1, note.A2))
	tr.nextBar()
	tr.At(M("0"), OnEvent(v1, note.A2))
	tr.nextBar()
	tr.At(M("0"), OffEvent(v1))

	tr.compile()

	fac := uint(500000000)

	if tr.Events[0].tick != 0 {
		t.Errorf("expecting 0, got %v", tr.Events[0].tick)
	}

	if tr.Events[1].tick != fac {
		t.Errorf("expecting %d, got %v", fac, tr.Events[1].tick)
	}

	if tr.Events[2].tick != fac*4 {
		t.Errorf("expecting %d, got %v", fac*4, tr.Events[2].tick)
	}

	if tr.Events[3].tick != fac*6 {
		t.Errorf("expecting %d, got %v", fac*6, tr.Events[3].tick)
	}

	if tr.Events[4].tick != fac*9 {
		t.Errorf("expecting %d, got %v", fac*9, tr.Events[4].tick)
	}
}

func TestBarOverflow(t *testing.T) {
	var v1 *Voice
	tr := newTrack(BPM(120), M("4/4"))
	tr.At(M("0"), OnEvent(v1, note.A2))
	tr.At(M("1/4"), OffEvent(v1))
	tr.At(M("4/4"), OnEvent(v1, note.A2))
	tr.At(M("6/4"), OffEvent(v1))

	tr.compile()

	fac := uint(500000000)

	if tr.Events[0].tick != 0 {
		t.Errorf("expecting 0, got %v", tr.Events[0].tick)
	}

	if tr.Events[1].tick != fac {
		t.Errorf("expecting %d, got %v", fac, tr.Events[1].tick)
	}

	if tr.Events[2].tick != fac*4 {
		t.Errorf("expecting %d, got %v", fac*4, tr.Events[2].tick)
	}

	if tr.Events[3].tick != fac*6 {
		t.Errorf("expecting %d, got %v", fac*6, tr.Events[2].tick)
	}
}
