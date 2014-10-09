package music

import (
	"testing"
)

func TestLinearDistributedValues(t *testing.T) {
	// from, to float64, steps int, dur Measure
	type test struct {
		from, to float64
		steps    int
		dur      Measure
	}

	type rest struct {
		width Measure
		diff  float64
	}

	var corpus = map[test]rest{
		// |0    15|
		//  1     2
		test{0, 15, 2, M("1")}: {M("1/2"), 15},

		// |0 10 20|
		//  1  2  3
		test{0, 20, 3, M("1")}: {M("1/3"), 10},

		// |0  5 10 15|
		//  1  2  3  4|
		test{0, 15, 4, M("1")}: {M("1/4"), 5},

		// |0  5 10 15 20|
		//  1  2  3  4  5
		test{0, 20, 5, M("1")}: {M("1/5"), 5},

		// |20  15 10|
		//   1   2  3
		test{20, 10, 3, M("1")}: {M("1/3"), -5},
	}

	for in, out := range corpus {
		width, diff := LinearDistributedValues(in.from, in.to, in.steps, in.dur)
		if width != out.width {
			t.Errorf("width is %f, but should be %f", width, out.width)
		}
		if diff != out.diff {
			t.Errorf("diff is %f, but should be %f", diff, out.diff)
		}
	}

}

// ergebnis: 1² = 1 | 2² = 4 | 3² = 9 | 4² = 16
func TestExponentialDistributedValues(t *testing.T) {
	// from, to float64, steps int, dur Measure
	type test struct {
		from, to float64
		steps    int
		dur      Measure
	}

	type rest struct {
		width Measure
		diffs []float64
	}

	var corpus = map[test]rest{
		// | 1² 2² 3²  4²|
		//   1  4  9  16
		//   1  2  3   4
		test{1, 16, 4, M("1")}: {M("1/4"), []float64{1, 4, 9, 16}},

		// | 1²+3 2²+3 3²+3 4²+3|
		//      4    7   12   19
		//   1    2    3    4
		test{4, 19, 4, M("1")}: {M("1/4"), []float64{4, 7, 12, 19}},

		// | 1²+3 2²+3 3²+3|
		//      4    7   12
		//   1    2    3
		test{4, 12, 3, M("1")}: {M("1/3"), []float64{4, 7, 12}},

		// | 4²+3 3²+3 2²+3 1²+3   |
		//      19  12    7   4
		//   1    2    3    4
		test{19, 4, 4, M("1")}: {M("1/4"), []float64{19, 12, 7, 4}},

		// | 4² 3² 2²  1²|
		//   16 9  4   1
		//   1  2  3   4
		test{16, 1, 4, M("1")}: {M("1/4"), []float64{16, 9, 4, 1}},
	}

	for in, out := range corpus {
		width, diffs := ExponentialDistributedValues(in.from, in.to, in.steps, in.dur)
		if width != out.width {
			t.Errorf("width is %f, but should be %f", float64(width), float64(out.width))
		}

		if len(diffs) != len(out.diffs) {
			t.Errorf("len diffs is %d, but should be %d", len(diffs), len(out.diffs))
		}

		for i, diff := range diffs {
			if diff != out.diffs[i] {
				t.Errorf("diff[%d] is %f, but should be %f", i, diff, out.diffs[i])
			}
		}
	}

}
