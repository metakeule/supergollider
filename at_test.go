package music

/*
import (
	"testing"
)

func TestAtOne(t *testing.T) {

	tr := newTrack(BPM(120), M("1/4")*4)
	// fmt.Printf("%d\n", int(tr.Tempi[0].Tempo.MilliSecs(M("1/4"))))

	tr.At(M("0"), On(v1, note(2)))
	tr.At(M("1/4"), Off(v1))
	tr.Compile()

	if tr.Events[1].Tick != 500000000 {
		t.Errorf("expecting 500000000, got %v", tr.Events[1].Tick)
	}

	tr.SetTempo(M("0"), BPM(240))
	tr.Compile()

	if tr.Events[1].Tick != 250000000 {
		t.Errorf("expecting 250000000, got %v", tr.Events[1].Tick)
	}

}

func TestAtSeq(t *testing.T) {
	tr := newTrack(BPM(120), M("1/4")*4)

	tr.At(M("0"), On(v1, note(2)))
	tr.At(M("1/4"), Off(v1))
	tr.SetTempo(M("2/4"), BPM(240))
	tr.At(M("2/4"), On(v1, note(2)))
	tr.At(M("3/4"), Off(v1))
	tr.nextBar()
	tr.At(M("0"), On(v1, note(2)))
	tr.At(M("1/4"), Off(v1))

	tr.Compile()

	if tr.Events[0].Tick != 0 {
		t.Errorf("expecting 0, got %v", tr.Events[0].Tick)
	}

	if tr.Events[1].Tick != 500000000 {
		t.Errorf("expecting 500000000, got %v", tr.Events[1].Tick)
	}

	if tr.Events[2].Tick != 1000000000 {
		t.Errorf("expecting 1000000000, got %v", tr.Events[2].Tick)
	}

	if tr.Events[3].Tick != 1250000000 {
		t.Errorf("expecting 1250000000, got %v", tr.Events[3].Tick)
	}

	if tr.Events[4].Tick != 1500000000 {
		t.Errorf("expecting 1500000000, got %v", tr.Events[4].Tick)
	}

	if tr.Events[5].Tick != 1750000000 {
		t.Errorf("expecting 1750000000, got %v", tr.Events[5].Tick)
	}
}

func TestBarChange(t *testing.T) {
	tr := newTrack(BPM(120), M("4/4"))

	tr.At(M("0"), On(v1, note(2)))
	tr.At(M("1/4"), Off(v1))
	tr.changeBar(M("3/4"))
	tr.At(M("1/4"), On(v1, note(2)))
	tr.nextBar()
	tr.At(M("0"), On(v1, note(2)))
	tr.nextBar()
	tr.At(M("0"), Off(v1))

	tr.Compile()

	fac := uint(500000000)

	if tr.Events[0].Tick != 0 {
		t.Errorf("expecting 0, got %v", tr.Events[0].Tick)
	}

	if tr.Events[1].Tick != fac {
		t.Errorf("expecting %d, got %v", fac, tr.Events[1].Tick)
	}

	if tr.Events[2].Tick != fac*4 {
		t.Errorf("expecting %d, got %v", fac*4, tr.Events[2].Tick)
	}

	if tr.Events[3].Tick != fac*6 {
		t.Errorf("expecting %d, got %v", fac*6, tr.Events[3].Tick)
	}

	if tr.Events[4].Tick != fac*9 {
		t.Errorf("expecting %d, got %v", fac*9, tr.Events[4].Tick)
	}
}

func TestBarOverflow(t *testing.T) {
	tr := newTrack(BPM(120), M("4/4"))
	tr.At(M("0"), On(v1, note(2)))
	tr.At(M("1/4"), Off(v1))
	tr.At(M("4/4"), On(v1, note(2)))
	tr.At(M("6/4"), Off(v1))

	tr.Compile()

	fac := uint(500000000)

	if tr.Events[0].Tick != 0 {
		t.Errorf("expecting 0, got %v", tr.Events[0].Tick)
	}

	if tr.Events[1].Tick != fac {
		t.Errorf("expecting %d, got %v", fac, tr.Events[1].Tick)
	}

	if tr.Events[2].Tick != fac*4 {
		t.Errorf("expecting %d, got %v", fac*4, tr.Events[2].Tick)
	}

	if tr.Events[3].Tick != fac*6 {
		t.Errorf("expecting %d, got %v", fac*6, tr.Events[2].Tick)
	}
}

*/
