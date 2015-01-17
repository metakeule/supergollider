package supergollider

import (
	"bytes"
	"github.com/metakeule/supergollider/note"
	"testing"
)

func mkVoice(instrName string) *Voice {
	var v1 = &Voice{}
	v1.scnode = 5
	var instr = &sCInstrument{}
	instr.name = instrName
	v1.instrument = instr
	return v1
}

var v1 = mkVoice("1")
var v2 = mkVoice("2")
var bpm120 = BPM(120)

func TestPrintEmpty(t *testing.T) {
	tr := &Track{}
	tr.compile()
	var bf bytes.Buffer
	tr.Print(BPM(120), "1/4", &bf)

	res := ""

	if bf.String() != res {
		t.Errorf("wrong res, expected: %#v, got: %#v", res, bf.String())
	}
}

func TestPrintOneNote(t *testing.T) {
	tr := newTrack(bpm120, M("4/4"))
	tr.At(M("0/4"), OnEvent(v1, note.A6))
	tr.At(M("2/4"), OffEvent(v1))
	tr.compile()

	var bf bytes.Buffer
	//	tr.Print(200000000, &bf)
	tr.Print(bpm120, "1/8", &bf)

	res := "1: A6-------\n"

	if bf.String() != res {
		t.Errorf("wrong res, expected: %#v, got: %#v", res, bf.String())
	}

	bf.Reset()
	tr.Print(bpm120, "1/4", &bf)
	res = "1: A6---\n"

	if bf.String() != res {
		t.Errorf("wrong res, expected: %#v, got: %#v", res, bf.String())
	}
}

func TestPrintPolyOneNote(t *testing.T) {
	tr := newTrack(bpm120, M("4/4"))
	tr.At(M("0/4"), OnEvent(v1, note.A6))
	tr.At(M("0/4"), OnEvent(v2, note.C7))
	tr.At(M("2/4"), OffEvent(v1))
	tr.At(M("2/4"), OffEvent(v2))
	tr.compile()
	var bf bytes.Buffer
	tr.Print(bpm120, "1/8", &bf)

	res := "" +
		"1: A6-------\n" +
		"2: C7-------\n"

	if bf.String() != res {
		t.Errorf("wrong res, expected: %#v, got: %#v", res, bf.String())
	}

}

func TestPrintMulti(t *testing.T) {
	tr := newTrack(bpm120, M("4/4"))
	tr.At(M("2/4"), OnEvent(v1, note.A6))
	tr.At(M("4/4"), OffEvent(v1))
	tr.At(M("6/4"), OnEvent(v1, note.C7))
	tr.At(M("8/4"), OffEvent(v1))
	tr.compile()

	var bf bytes.Buffer

	tr.Print(bpm120, "1/8", &bf)

	res := "1: " +
		"        " +
		"A6-------" +
		"        " +
		"C7-------" +
		"\n"

	if bf.String() != res {
		t.Errorf("wrong res, expected: %#v, got: %#v", res, bf.String())
	}
}

func TestPrintPolyMultiParallel(t *testing.T) {
	tr := newTrack(bpm120, M("4/4"))
	tr.At(M("0/4"), OnEvent(v1, note.A6))
	tr.At(M("0/4"), OnEvent(v2, note.C7))
	tr.At(M("2/4"), OffEvent(v1))
	tr.At(M("2/4"), OffEvent(v2))

	tr.At(M("4/4"), OnEvent(v1, note.D6))
	tr.At(M("4/4"), OnEvent(v2, note.B7))
	tr.At(M("6/4"), OffEvent(v1))
	tr.At(M("6/4"), OffEvent(v2))
	tr.compile()

	var bf bytes.Buffer

	tr.Print(bpm120, "1/8", &bf)

	res := "" +
		"1: A6-------        D6-------\n" +
		"2: C7-------        B7-------\n"

	if bf.String() != res {
		t.Errorf("wrong res, expected: %#v, got: %#v", res, bf.String())
	}
}

func TestPrintEach(t *testing.T) {
	tr := newTrack(bpm120, M("4/4"))
	tr.SetLoop("metronome", v1.Metronome(M("1/4"), note.C2))
	tr.Start()
	tr.Fill(3)

	tr.compile()
	var bf bytes.Buffer
	tr.Print(bpm120, "1/8", &bf)
	_ = tr

	res := "1: C2-  C2-  C2-  C2-  C2-  C2-  C2-  C2-  C2-  C2-  C2-  C2-  C2-  C2-  C2-  C2-\n"

	if bf.String() != res {
		t.Errorf("expected: %#v, got: %#v", res, bf.String())
	}
}

func TestPrintPolyPhon(t *testing.T) {
	tr := newTrack(bpm120, M("4/4"))
	tr.Start()

	tr.At(M("0"), OnEvent(v1, note.C8), OnEvent(v2, note.C6))
	tr.At(M("1/4"), OffEvent(v1))
	tr.At(M("2/4"), OnEvent(v1, note.C7), OffEvent(v2))
	tr.At(M("3/4"), OffEvent(v1), OnEvent(v2, note.C1))
	tr.nextBar()
	tr.At(M("0"), OnEvent(v1, note.C8), OffEvent(v2))
	tr.At(M("1/4"), OffEvent(v1), OnEvent(v2, note.C2))
	tr.At(M("2/4"), OnEvent(v1, note.C7))
	tr.SetTempo(M("2/4"), BPM(240))
	tr.At(M("3/4"), OffEvent(v1))
	tr.nextBar()
	tr.At(M("0"), OnEvent(v1, note.C8))
	tr.At(M("0"), OffEvent(v2))
	tr.At(M("1/4"), OnEvent(v2, note.C6), OffEvent(v1))
	tr.At(M("2/4"), OnEvent(v1, note.C7), OffEvent(v2))
	tr.At(M("3/4"), OffEvent(v1))
	tr.compile()

	var bf bytes.Buffer

	tr.Print(bpm120, "1/16", &bf)

	res := "" +
		"1: C8-------        C7-------        C8-------        C7---    C8---    C7---" + "\n" +
		"2: C6---------------        C1-------        C2---------------    C6---" + "\n"

	if bf.String() != res {
		t.Errorf("wrong res, expected: %#v, got: %#v", res, bf.String())
	}
}
