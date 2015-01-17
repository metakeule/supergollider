package supergollider

import (
	"bytes"
	"fmt"
	"github.com/metakeule/supergollider/note"
	"testing"
)

type tempoChange struct {
	position string
	tempo    float64
	// expected string
}

func TestTempoChange(t *testing.T) {

	tests := map[string]tempoChange{
		"1: C1--    C1--    C1--    C1--": {"0/4", 120},
		"1: C1--    C1--    C1--    C1":   {"3/4", 240},
		"1: C1--    C1--    C1  C1":       {"2/4", 240},
		"1: C1--    C1  C1  C1":           {"1/4", 240},
		"1: C1  C1  C1  C1":               {"0/4", 240},
	}

	for expected, tempoChange := range tests {
		tr := newTrack(bpm120, M("4/4"))
		tr.SetLoop("m1", v1.Metronome(M("1/4"), note.C1))
		tr.Start()
		tr.At(M(tempoChange.position), BPM(tempoChange.tempo).Event())
		var bf bytes.Buffer
		tr.Print(bpm120, "1/16", &bf)
		if got, want := bf.String(), expected+"\n"; got != want {
			t.Errorf("SetTempo(M(%#v), BPM(%v))\n\ngot:\n%s\n\nwanted:\n%s\n\n", tempoChange.position, tempoChange.tempo, got, want)
		}
	}
}

func TestMultipleTempoChanges(t *testing.T) {
	tests := map[string][]tempoChange{
		"1: C1--    C1--    C1--    C1--": {{"0/4", 120}},
		"1: C1--    C1  C1  C1--":         {{"1/4", 240}, {"3/4", 120}},
		"1: C1--  C1  C1C1":               {{"1/8", 240}, {"5/8", 480}},
		"1: C1--  C1  C1    C1--":         {{"1/8", 240}, {"5/8", 120}},
	}

	for expected, tempoChanges := range tests {
		tr := newTrack(bpm120, M("4/4"))
		tr.SetLoop("m1", v1.Metronome(M("1/4"), note.C1))
		tr.Start()

		var info bytes.Buffer

		for _, tempoChange := range tempoChanges {
			tr.At(M(tempoChange.position), BPM(tempoChange.tempo).Event())
			fmt.Fprintf(&info, "SetTempo(M(%#v), BPM(%v));", tempoChange.position, tempoChange.tempo)
		}

		var bf bytes.Buffer
		tr.Print(bpm120, "1/16", &bf)
		if got, want := bf.String(), expected+"\n"; got != want {
			t.Errorf("%s\n\ngot:\n%s\n\nwanted:\n%s\n\n", info.String(), got, want)
		}
	}

}

func TestSeqTempo(t *testing.T) {
	tests := map[string][]string{
		"1: C1--    C1  C1  C1": {"1/4"},
		"1: C1--  C1  C1  C1":   {"1/8"},
		"1: C1--  C1  C1C1":     {"1/8", "4/8"},
	}

	for expected, tempoPos := range tests {
		sT := SeqTempo(120, 120, StepAdd)

		tr := newTrack(bpm120, M("4/4"))
		tr.SetLoop("m1", v1.Metronome(M("1/4"), note.C1))
		patterns := []Pattern{}

		for _, pos := range tempoPos {
			patterns = append(patterns, sT.SetTempo(pos))
		}

		tr.Start(patterns...)
		var bf bytes.Buffer
		tr.Print(bpm120, "1/16", &bf)

		if got, want := bf.String(), expected+"\n"; got != want {
			t.Errorf("\ngot:\n%s\n\nwanted:\n%s\n\n", got, want)
		}
	}
}
