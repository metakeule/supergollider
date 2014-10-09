package music

type bus int

func (b bus) Name() string {
	return "bushub"
}

var bushub = bus(0)
var busses = map[string]int{}

func newBus(g generator, name string) *Voice {
	if _, has := busses[name]; has {
		panic("bus with name " + name + " already defined")
	}
	busses[name] = g.newBusId()
	return &Voice{generator: g, instrument: bushub, Bus: busses[name]}
}
