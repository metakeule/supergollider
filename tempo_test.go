package supergollider

import (
	"bytes"
	"github.com/metakeule/supergollider/note"
	"testing"
)

func TestTempoChange(t *testing.T) {

	tests := []struct {
		position string
		tempo    float64
		expected string
	}{
		{"0/4", 120, "1: C1--    C1--    C1--    C1--"},
		{"3/4", 240, "1: C1--    C1--    C1--    C1"},
		{"2/4", 240, "1: C1--    C1--    C1  C1"},
		{"1/4", 240, "1: C1--    C1  C1  C1"},
		{"0/4", 240, "1: C1  C1  C1  C1"},
	}

	for _, test := range tests {
		tr := newTrack(bpm120, M("4/4"))
		tr.SetLoop("m1", v1.Metronome(M("1/4"), note.C1))
		tr.Start()
		tr.SetTempo(M(test.position), BPM(test.tempo))
		tr.compile()
		var bf bytes.Buffer
		tr.Print(bpm120, "1/16", &bf)
		if got, want := bf.String(), test.expected+"\n"; got != want {
			t.Errorf("SetTempo(M(%#v), BPM(%v))\n\ngot:\n\n%s\n\nwanted:\n%s\n\n", test.position, test.tempo, got, want)
		}
	}
}

// SeqTempo
