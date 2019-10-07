package main

import (
	"log"
	"github.com/mloncode/sonic"
	"github.com/mloncode/sonic/src/sound"
)

func main() {
	out, err := sound.MidiOut()

	if err != nil {
		log.Fatal(err)
		return
	}

	m1 := sound.NewMarkov("song1.midi")
	m2 := sound.NewMarkov("song2.midi")

	oldChanges := sound.NewSequence("prophet", sonic.ConvertMarkov(m1, sonic.File1.Old))
	newChanges := sound.NewSequence("prophet", sonic.ConvertMarkov(m2, sonic.File1.New))

	oldChanges.Play(out)
	newChanges.Play(out)
}
