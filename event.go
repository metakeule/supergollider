package music

import "bytes"

type Event struct {
	Voice                     *Voice
	Params                    Parameter // a special parameter offset may be used to set a per event offset
	runner                    func(*Event)
	type_                     string
	tick                      uint
	absPosition               Measure // will be enabled when integrated
	offset                    float64 // offset added to the final position (includes instrument and sample offsets as well as offset set via parameter)
	sccode                    bytes.Buffer
	changedParamsPrepared     map[string]float64
	sampleInstrumentFrequency float64
	synthID                   int
	reference                 *Event
}

var fin = &Event{runner: func(*Event) {}, type_: "fin"}
var start = &Event{runner: func(*Event) {}, type_: "start"}

func newEvent(v *Voice, type_ string) *Event {
	return &Event{
		Voice:  v,
		Params: ParamsMap(map[string]float64{}),
		type_:  type_,
	}
}

// merges the given params of the event into a clone
// of ev, returning the clone
// may be used with events that have modifiers, like Scale, Rhythm etc
// the given voice is set and we get an On event
//func (ev *Event) OnMerged(voice Voice, m map[string]float64) *Event {
func (ev *Event) OnMerged(voice *Voice, ps ...Parameter) *Event {
	n := ev.Clone()
	p := []Parameter{ev.Params}
	p = append(p, ps...)
	n.Params = MixParams(p...)
	n.Voice = voice
	n.runner = voice.OnEvent
	n.type_ = "ON"
	return n
}

// merges the given params of the event into a clone
// of ev, returning the clone
// may be used with events that have modifiers, like Scale, Rhythm etc
// the given voice is set and we get a change event
//func (ev *Event) ChangeMerged(voice Voice, m map[string]float64) *Event {
func (ev *Event) ChangeMerged(voice *Voice, ps ...Parameter) *Event {
	n := ev.Clone()
	p := []Parameter{ev.Params}
	p = append(p, ps...)
	n.Params = MixParams(p...)
	n.Voice = voice
	n.runner = voice.ChangeEvent
	n.type_ = "CHANGE"
	return n
}

func (ev *Event) Clone() *Event {
	n := &Event{Voice: ev.Voice, runner: ev.runner}
	n.type_ = ev.type_
	n.absPosition = ev.absPosition
	n.Params = ev.Params
	return n
}

var OnEvent = EventGenerator(func(v *Voice, params ...Parameter) *Event {
	return &Event{
		Voice:  v,
		Params: MixParams(params...),
		runner: v.OnEvent,
		type_:  "ON",
	}
})

// params are ignored, just to fullfill the EventGenerator interface
var OffEvent = EventGenerator(func(v *Voice, params ...Parameter) *Event {
	return &Event{
		Voice:  v,
		runner: v.OffEvent,
		type_:  "OFF",
	}
})

// params are ignored, just to fullfill the EventGenerator interface
var MuteEvent = EventGenerator(func(v *Voice, params ...Parameter) *Event {
	return &Event{
		Voice:  v,
		runner: v.OffEvent,
		type_:  "MUTE",
	}
})

// params are ignored, just to fullfill the EventGenerator interface
var UnMuteEvent = EventGenerator(func(v *Voice, params ...Parameter) *Event {
	return &Event{
		Voice:  v,
		runner: v.donothing,
		type_:  "UNMUTE",
	}
})

var ChangeEvent = EventGenerator(func(v *Voice, params ...Parameter) *Event {
	return &Event{
		Voice:  v,
		Params: MixParams(params...),
		runner: v.ChangeEvent,
		type_:  "CHANGE",
	}
})

func newFreeEvent(v *Voice, reference *Event) *Event {
	return &Event{
		Voice:     v,
		type_:     "FREE",
		reference: reference,
	}
}

func CustomEvent(fn func(*Event)) *Event {
	return &Event{
		runner: fn,
		type_:  "CUSTOM",
	}
}

type setBpm struct {
	pos Measure
	bpm float64
}

func (s *setBpm) Events(barNum int, t Tracker) map[Measure][]*Event {
	return map[Measure][]*Event{
		s.pos: []*Event{
			&Event{
				Params: Param("bpm", s.bpm),
				type_:  "TEMPO_CHANGE",
			},
		},
	}
}

func (s *setBpm) NumBars() int {
	return 1
}

func SetBPM(pos string, bpm float64) *setBpm {
	return &setBpm{M(pos), bpm}
}

type EventGenerator func(v *Voice, params ...Parameter) *Event

/*
func ExecEvent(v *Voice, fn func(e *Event), params ...Parameter) *Event {
	return &Event{
		Voice:  v,
		Params: Params(params...),
		runner: fn,
		Type:   "EXEC",
	}
}
*/

// returns a Pattern for an event func at a certain position
func EventFuncPattern(pos string, fn func(e *Event)) Pattern {
	return PatternFunc(func(barNum int, t Tracker) map[Measure][]*Event {
		return map[Measure][]*Event{M(pos): []*Event{CustomEvent(fn)}}
	})
}

func Exec(pos string, fns ...func()) Pattern {
	return EventFuncPattern(pos, func(e *Event) {
		for _, fn := range fns {
			fn()
		}
	})
}
