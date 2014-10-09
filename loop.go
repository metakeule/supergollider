package music

func RepeatPattern(n int, loop Pattern) Pattern {
	ls := make([]Pattern, n)

	for i := 0; i < n; i++ {
		ls[i] = loop
	}

	return SeqPatterns(ls...)
}

func Loop(patterns ...Pattern) *loop {
	l := &loop{}
	l.Next(patterns...)
	return l
}

/*
func (l *loop) TrackLoop(name string) *TrackLoop {
	return &TrackLoop{loop: l, Name: name}
}
*/
func (l *loop) TrackLoop() *TrackLoop {
	return &TrackLoop{loop: l}
}

type loop struct {
	currentBar int
	patterns   []Pattern
}

func (l *loop) Events(barNum int, t Tracker) map[Measure][]*Event {
	num := 0

	for _, p := range l.patterns {
		next := num + p.NumBars()
		if barNum < next {
			return p.Events(barNum-num, t)
		}
		num = next
	}
	return nil
}

func (l *loop) NumBars() int {
	num := 0

	for _, p := range l.patterns {
		num += p.NumBars()
	}
	return num
}

func (l *loop) Next(patterns ...Pattern) *loop {
	// l.nextBar()
	l.patterns = append(l.patterns, MixPatterns(patterns...))
	return l
}

type loopInTrack struct {
	start        Measure
	currentIndex int
	loop         Pattern
}

func (l *loopInTrack) setEventsForBar(t *Track) {
	for pos, events := range l.loop.Events(l.currentIndex, t) {
		t.At(l.start+pos, events...)
	}
	if l.currentIndex >= l.loop.NumBars()-1 {
		l.currentIndex = 0
	} else {
		l.currentIndex++
	}
}

type TrackLoop struct {
	// stopAt *Measure
	*loop
	// Name string
}

type TrackLoopStart struct {
	*TrackLoop
	start        Measure
	currentIndex int
}

type TrackLoopReplacer struct {
	newLoop *loop
	start   Measure
	*TrackLoop
}

func (l *TrackLoop) ReplaceAt(start string, n *loop) *TrackLoopReplacer {
	return &TrackLoopReplacer{
		newLoop:   n,
		TrackLoop: l,
		start:     M(start),
	}
}

func (l *TrackLoop) Replace(n *loop) *TrackLoopReplacer {
	return l.ReplaceAt("0", n)
}

func (l *TrackLoop) OnAt(pos string) *TrackLoopStart {
	return &TrackLoopStart{
		TrackLoop:    l,
		start:        M(pos),
		currentIndex: 0,
	}
}

func (l *TrackLoop) On() *TrackLoopStart {
	return l.OnAt("0")
}

// call this function at the end of each bar
func (l *TrackLoopStart) setEventsForBar(t *Track) {
	// fmt.Printf("loop %s T %d\n", l.Name, t.BarNum())
	for pos, events := range l.loop.Events(l.currentIndex, t) {
		// fmt.Printf("loop %s T %d Pos %d\n", l.Name, t.BarNum(), l.start+pos)
		t.At(l.start+pos, events...)
	}
	if l.currentIndex >= l.loop.NumBars()-1 {
		l.currentIndex = 0
	} else {
		l.currentIndex++
	}
}

type TrackLoopStop struct {
	// stopAt Measure
	*TrackLoop
}

func (l *TrackLoop) Off() *TrackLoopStop {
	return &TrackLoopStop{
		// stopAt:    M(pos),
		TrackLoop: l,
	}
}
