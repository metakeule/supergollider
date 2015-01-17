package supergollider

type group struct{}

func (g group) Name() string {
	return "group"
}

func newGroup(g generator) *Voice {
	return &Voice{generator: g, instrument: group{}, Group: g.newGroupId()}
}
