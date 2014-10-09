package music

// the outer invoker may use the first voices instrument to query loadcode etc
func newRoute(g generator, name, path string, numVoices int) []*Voice {
	instr := &sCInstrument{
		name: name,
		Path: path,
	}
	return _voices(numVoices, g, instr, 1200)
}
