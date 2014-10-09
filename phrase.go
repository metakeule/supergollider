package music

type phrase struct {
	voice      *Voice
	startPos   float64
	currentPos float64
	patterns   []Pattern
}

func (m *phrase) Events(barNum int, t Tracker) map[Measure][]*Event {
	res := map[Measure][]*Event{}
	pos := m.startPos
	for _, p := range m.patterns {
		//p.Pattern(t)
		// WARNING: does only work as long as each pattern has only one key (pos)
		for evPos, events := range p.Events(barNum, t) {
			currentPos := Measure(int(pos)) + evPos
			res[currentPos] = append(res[currentPos], events...)
			pos += float64(int(evPos))
		}
	}
	return res
}

/*
func (m *phrase) Pattern(t Tracker) {
	for _, p := range m.patterns {
		p.Pattern(t)
	}
}
*/

func newPhrase(v *Voice, pos string) *phrase {
	return &phrase{
		voice:    v,
		startPos: _M(pos),
	}
}

/*
 */

func (m *phrase) currentMeasure() Measure {
	return Measure(int(m.startPos + m.currentPos))
}

func (m *phrase) Play(distance string, params ...Parameter) *phrase {
	m.currentPos += _M(distance)
	m.patterns = append(m.patterns, &play{m.currentMeasure(), m.voice, MixParams(params...)})
	return m
}

func (m *phrase) PlayDur(distance string, dur string, params ...Parameter) *phrase {
	m.currentPos += _M(distance)
	m.patterns = append(m.patterns, &play{m.currentMeasure(), m.voice, MixParams(params...)})
	m.patterns = append(m.patterns, &stop{m.currentMeasure() + M(dur), m.voice})
	return m
}

func (m *phrase) Stop(distance string) *phrase {
	m.currentPos += _M(distance)
	m.patterns = append(m.patterns, &stop{m.currentMeasure(), m.voice})
	return m
}

func (m *phrase) Modify(distance string, params ...Parameter) *phrase {
	m.currentPos += _M(distance)
	m.patterns = append(m.patterns, &mod{m.currentMeasure(), m.voice, MixParams(params...)})
	return m
}
