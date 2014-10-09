package chord

import (
	"github.com/metakeule/music/note"
	"github.com/metakeule/music/scale"
)

// a Chord is a periodic scale but typically has less tones
func Chord(base note.Note, steps ...uint) *scale.Periodic {
	return &scale.Periodic{Steps: steps, NumChromaticSteps: 12, BaseNote: base}
}

func Dim(base note.Note) *scale.Periodic      { return Chord(base, 3, 3, 6) }
func Aug(base note.Note) *scale.Periodic      { return Chord(base, 4, 4, 4) }
func Dur(base note.Note) *scale.Periodic      { return Chord(base, 4, 3, 5) }
func DurMaj7(base note.Note) *scale.Periodic  { return Chord(base, 4, 3, 4, 1) }
func DurMin7(base note.Note) *scale.Periodic  { return Chord(base, 4, 3, 3, 2) }
func Moll(base note.Note) *scale.Periodic     { return Chord(base, 3, 4, 5) }
func MollMaj7(base note.Note) *scale.Periodic { return Chord(base, 3, 4, 4, 1) }
func MollMin7(base note.Note) *scale.Periodic { return Chord(base, 3, 4, 3, 2) }
