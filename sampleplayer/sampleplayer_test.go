package sampleplayer

import (
	"testing"
)

func TestDynMapper(t *testing.T) {
	m := map[float64]string{
		0:   "pianissimo",
		0.3: "piano",
		0.4: "mezzoforte",
		0.5: "forte",
		0.7: "fortissimo",
	}

	corpus := []struct {
		dyn float64
		str string
	}{
		{0.2, "pianissimo"},
		{0.3, "piano"},
		{0.45, "mezzoforte"},
		{0.9, "fortissimo"},
	}

	for _, test := range corpus {
		got := dynMapper(m, test.dyn)
		if got != test.str {
			t.Errorf("wrong string for %0.2f: %#v (wanted %#v)", test.dyn, got, test.str)
		}
	}
	/*
		if got := ampMapper(nil, 0.6); got != "" {
			t.Errorf("empty map must return empty string, but returns %#v", got)
		}
	*/
}

func TestScanner(t *testing.T) {
	sc := NewPhilharmonicScanner()
	sc.ScanDir("/media/usb0/philharmonic-orchestra/violin/")
}

func TestAmpMapperEmpty(t *testing.T) {

	defer func() {
		if x := recover(); x == nil {
			t.Fatalf("no panic")
		}
	}()

	dynMapper(nil, 0.6)
}

func TestFreqMapper(t *testing.T) {
	m := map[float64]string{
		440:  "A1",
		880:  "A2",
		1760: "A3",
		3520: "A4",
	}

	corpus := []struct {
		searchFreq float64
		targetFreq float64
		targetStr  string
	}{
		{220, 440, "A1"},
		{660, 440, "A1"},
		{1000, 880, "A2"},
		{4000, 3520, "A4"},
	}

	for _, test := range corpus {
		gotStr, gotFreq := freqMapper(m, test.searchFreq)
		if gotStr != test.targetStr {
			t.Errorf("wrong string for %0.2f: %#v (wanted %#v)", test.searchFreq, gotStr, test.targetStr)
		}

		if gotFreq != test.targetFreq {
			t.Errorf("wrong freq for %0.2f: %0.2f (wanted %0.2f)", test.searchFreq, gotFreq, test.targetFreq)
		}
	}

}

func TestFreqMapperEmpty(t *testing.T) {

	defer func() {
		if x := recover(); x == nil {
			t.Fatalf("no panic")
		}
	}()

	freqMapper(nil, 0.6)
}
