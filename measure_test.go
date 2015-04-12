package supergollider

import "testing"

func TestM(t *testing.T) {
	corpus := map[string]int{
		"1/2 + 1/3":       int(M("1/2")) + int(M("1/3")),
		"1/2 + 1/2":       int(M("1")),
		"1/2 + 1/4 + 1/4": int(M("1")),
	}

	for s, m := range corpus {
		if int(M(s)) != m {
			t.Errorf("expecting %v for %s, but got %v", int(M(s)), s, m)
		}
	}

}

func equals(a, b []Measure) bool {
	for i, aa := range a {
		if aa != b[i] {
			return false
		}
	}
	return true
}

func TestMeasureScale(t *testing.T) {

	tests := []struct {
		total     Measure
		relations []float64
		expected  []Measure
	}{
		{M("10"), []float64{1, 3, 1}, []Measure{M("2"), M("6"), M("2")}},
		{M("10"), []float64{3, 1, 1}, []Measure{M("6"), M("2"), M("2")}},
		{M("2 + 1/2"), []float64{3, 1, 1}, []Measure{M("1 + 1/2"), M("1/2"), M("1/2")}},
		{M("7"), []float64{2, 4, 4, 4}, []Measure{M("1"), M("2"), M("2"), M("2")}},
		{M("3/4"), []float64{1, 2, 2, 1}, []Measure{M("1/8"), M("1/4"), M("1/4"), M("1/8")}},
	}

	for _, test := range tests {
		if got, want := test.total.Scale(test.relations...), test.expected; !equals(got, want) {
			t.Errorf(`M("%v").Scale(%v...) = %v; want %v`, test.total, test.relations, got, want)
		}
	}

}
