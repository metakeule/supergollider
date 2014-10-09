package music

/*
TODO:

1. gleichmäßiges Verteilen von Werten einer best. Spanne und schrittweite
   auf eine Distanz
   z.B. von 1-15 in 8 Schritten
   |1-3-5-7-9-11-13-15|

2. gleichmäßiges Verteilen von schrittweitenänderungen einer bestimmten spanne und
   distanz mit minimaler schrittweitenänderung
   z.B. von schrittweite 1 auf schrittweite 5 in minimaler schrittweitenänderung von 1
   |x-x--x---x----x-----x|
*/

// LinearDistributedValues fill duration dur linearily with values from from to to in steps steps
func LinearDistributedValues(from, to float64, steps int, dur Measure) (width Measure, diff float64) {
	diff = (to - from) / float64(steps-1)
	width = Measure(int(float64(int(dur)) / float64(steps)))
	return
}

/*

f(x) = ax²

start =  1
ziel  = 16

ergebnis: 1² = 1 | 2² = 4 | 3² = 9 | 4² = 16
4 schritte
schrittweiten sind 0, 3, 5, 7, wir wollen also
diffs := []float64(1,4,9,16)

alle Punkte liegen auf der Kurve, also auch 1 und 4

to = a (steps)²
a = to / (steps)²


from to
4 => 19

n²+from-1
1²+4-1   2²+4-1    3²+4-1     4²+4-1
4             7        12         19

a1²+from-1 = from
a1 = 1

a(steps)²+from-1 = to
a(steps)² = to-from+1
a = to-from+1 / (steps)²

*/

// ExponentialDistributedValues returns the width (time difference) and values for and exponential growth from
// the given value from to the given value to in the total time of dur with the given number of steps
func ExponentialDistributedValues(from, to float64, steps int, dur Measure) (width Measure, diffs []float64) {
	reverse := false
	if to < from {
		to, from = from, to
		reverse = true
	}
	diff := (to - from) + 1.0
	factor := diff / float64(steps) / float64(steps)
	diffs = make([]float64, steps)

	for i := 0; i < steps; i++ {
		key := i
		if reverse {
			key = steps - i - 1
		}
		diffs[key] = (float64((i+1)*(i+1)) * factor) + from - 1
	}
	width = Measure(int(float64(int(dur)) / float64(steps)))
	return
}
