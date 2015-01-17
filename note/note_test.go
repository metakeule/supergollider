package note

import (
	"testing"
)

var tests = []struct {
	midi float64
	freq float64
}{
	// 40hz - 10000hz
	{27, 38.89},
	{31, 49},
	{32, 51.91},
	{33, 55},
	{34, 58.27},
	{45, 110},
	{57, 220},
	{69, 440},
	{71.214, 500},
	{70, 466.2},
	{71, 493.90},
	{80, 830.60},
	{77.037, 700},
	{81, 880},
	{124, 10550},
}

var precision = 4

func TestMidiToFreq(t *testing.T) {
	for _, test := range tests {
		if got, want := roundViaFloat(MidiToFreq(test.midi), precision), roundViaFloat(test.freq, precision); got != want {
			t.Errorf("MidiToFreq(%v) = %+.6f; want %+.6f", test.midi, got, want)
		}
	}
}

func TestFreqToMidi(t *testing.T) {
	for _, test := range tests {
		if got, want := roundViaFloat(FreqToMidi(test.freq), precision), roundViaFloat(test.midi, precision); got != want {
			t.Errorf("FreqToMidi(%v) = %+.6f; want %+.6f", test.freq, got, want)
		}
	}
}

func TestMidiToNote(t *testing.T) {

	tests := []struct {
		midi float64
		note Note
	}{
		{float64(C0), C0},
		{float64(C4), C4},
		{float64(B9), B9},
	}

	for _, test := range tests {

		if got, want := MidiToNote(test.midi), test.note; got != want {
			t.Errorf("MidiToNote(%v) = %v; want %v", test.midi, got, want)
		}
	}

}

func TestFreqToNote(t *testing.T) {
	tests := []struct {
		freq float64
		note Note
	}{
		{C0.Frequency(), C0},
		{C4.Frequency(), C4},
		{B9.Frequency(), B9},
		{440, A4},
		{B4.Frequency(), B4},
		{500, B4},
		{510, C5},
		{520, C5},
		{C5.Frequency(), C5},
	}

	for _, test := range tests {

		if got, want := FreqToNote(test.freq), test.note; got != want {
			t.Errorf("FreqToNote(%v) = %v; want %v", test.freq, got, want)
		}
	}
}

func TestRoundDec(t *testing.T) {

	tests := []struct {
		input    float64
		decs     int
		expected float64
	}{
		{2023.343, 2, 2023.34},
		{2023.346, 2, 2023.35},
		{23.346, 2, 23.35},
		{23.346, 0, 23},
		{2023.346, 0, 2023},
		{23.846, 0, 24},
	}

	for _, test := range tests {
		if got, want := roundDec(test.input, test.decs), test.expected; got != want {
			t.Errorf("roundDec(%.7f, %v) = %.7f; want %.7f", test.input, test.decs, got, want)
		}
	}

}
