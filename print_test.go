package music

/*
import (
	"bytes"
	"testing"
)

func note(n float64) map[string]float64 {
	return map[string]float64{"note": n}
}

type voice string

func (p voice) On(ev *Event)     {}
func (p voice) Off(ev *Event)    {}
func (p voice) Change(ev *Event) {}
func (p voice) Mute(ev *Event)   {}
func (p voice) UnMute(ev *Event) {}

func (v voice) Name() string {
	return string(v)
}

var v1 = voice("1")
var v2 = voice("2")

func TestPrintEmpty(t *testing.T) {
	tr := &Track{}
	tr.compiled = true

	var bf bytes.Buffer
	tr.Print(BPM(120), "1/4", &bf)

	res := ""

	if bf.String() != res {
		t.Errorf("wrong res, expected: %#v, got: %#v", res, bf.String())
	}
}

func TestPrintOneNote(t *testing.T) {
	tr := &Track{}
	tr.compiled = true
	e1 := On(v1, note(2))
	e2 := Off(v1)
	e1.Tick = 0
	e2.Tick = 1 * bil
	tr.Events = []*Event{e1, e2}

	var bf bytes.Buffer
	//	tr.Print(200000000, &bf)
	tr.Print(BPM(120), "1/8", &bf)

	res := "1: 2-------\n"

	if bf.String() != res {
		t.Errorf("wrong res, expected: %#v, got: %#v", res, bf.String())
	}

	bf.Reset()
	tr.Print(BPM(120), "1/4", &bf)
	res = "1: 2---\n"

	if bf.String() != res {
		t.Errorf("wrong res, expected: %#v, got: %#v", res, bf.String())
	}
}

var bil = uint(1000000000)

func TestPrintPolyOneNote(t *testing.T) {
	tr := &Track{}
	tr.compiled = true
	e1 := On(v1, note(2))
	e2 := On(v2, note(3))
	e3 := Off(v1)
	e4 := Off(v2)
	e1.Tick = 0
	e2.Tick = 0
	e3.Tick = 1 * bil
	e4.Tick = 1 * bil
	tr.Events = []*Event{e1, e2, e3, e4}

	var bf bytes.Buffer
	tr.Print(BPM(120), "1/8", &bf)

	res := "" +
		"1: 2-------\n" +
		"2: 3-------\n"

	if bf.String() != res {
		t.Errorf("wrong res, expected: %#v, got: %#v", res, bf.String())
	}

}

func TestPrintMulti(t *testing.T) {
	tr := &Track{}
	e1 := On(v1, note(2))
	e2 := Off(v1)
	e3 := On(v1, note(3))
	e4 := Off(v1)
	e1.Tick = 1 * bil
	e2.Tick = 2 * bil
	e3.Tick = 3 * bil
	e4.Tick = 4 * bil
	tr.Events = []*Event{e1, e2, e3, e4}
	tr.compiled = true

	var bf bytes.Buffer

	tr.Print(BPM(120), "1/8", &bf)

	res := "1: " +
		"        " +
		"2-------" +
		"        " +
		"3-------" +
		"\n"

	if bf.String() != res {
		t.Errorf("wrong res, expected: %#v, got: %#v", res, bf.String())
	}
}

func TestPrintPolyMultiParallel(t *testing.T) {
	tr := &Track{}
	e1 := On(v1, note(2))
	e2 := On(v2, note(4))
	e3 := Off(v1)
	e4 := Off(v2)
	e5 := On(v1, note(3))
	e6 := On(v2, note(5))
	e7 := Off(v1)
	e8 := Off(v2)
	e1.Tick = 0
	e2.Tick = 0
	e3.Tick = 1 * bil
	e4.Tick = 1 * bil
	e5.Tick = 2 * bil
	e6.Tick = 2 * bil
	e7.Tick = 3 * bil
	e8.Tick = 3 * bil
	tr.Events = []*Event{e1, e2, e3, e4, e5, e6, e7, e8}
	tr.compiled = true

	var bf bytes.Buffer

	tr.Print(BPM(120), "1/8", &bf)

	res := "" +
		"1: 2-------        3-------\n" +
		"2: 4-------        5-------\n"

	if bf.String() != res {
		t.Errorf("wrong res, expected: %#v, got: %#v", res, bf.String())
	}
}

// TODO: fill with too less values
// prefill with pause does not work
func TestPrintEach(t *testing.T) {

	tr := New("4/4", BPM(120),
		EachBar(
			Metronome(v1, M("1/4"), note(float64(2))),
		),
	).Fill(3)

	tr.Compile()
	var bf bytes.Buffer
	tr.Print(BPM(120), "1/8", &bf)
	_ = tr

	res := "1: 2-  2-  2-  2-  2-  2-  2-  2-  2-  2-  2-  2-  2-  2-  2-  2-\n"

	if bf.String() != res {
		t.Errorf("expected: %#v, got: %#v", res, bf.String())
	}
}


func TestPrintPolyPhon(t *testing.T) {
	tr := New("4/4", BPM(120))

	tr.At(M("0"), On(v1, note(8)), On(v2, note(6)))
	tr.At(M("1/4"), Off(v1))
	tr.At(M("2/4"), On(v1, note(7)), Off(v2))
	tr.At(M("3/4"), Off(v1), On(v2, note(1)))
	tr.nextBar()
	tr.At(M("0"), On(v1, note(8)), Off(v2))
	tr.At(M("1/4"), Off(v1), On(v2, note(2)))
	tr.At(M("2/4"), On(v1, note(7)))
	tr.SetTempo(M("2/4"), BPM(240))
	tr.At(M("3/4"), Off(v1))
	tr.nextBar()
	tr.At(M("0"), On(v1, note(8)))
	tr.At(M("0"), Off(v2))
	tr.At(M("1/4"), On(v2, note(6)), Off(v1))
	tr.At(M("2/4"), On(v1, note(7)), Off(v2))
	tr.At(M("3/4"), Off(v1))
	tr.Compile()

	var bf bytes.Buffer

	tr.Print(BPM(120), "1/16", &bf)

	res := "" +
		"1: 8-------        7-------        8-------        7---    8---    7---" + "\n" +
		"2: 6---------------        1-------        2---------------    6---" + "\n"

	if bf.String() != res {
		t.Errorf("wrong res, expected: %#v, got: %#v", res, bf.String())
	}
}

*/
