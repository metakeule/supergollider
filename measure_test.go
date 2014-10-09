package music

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
