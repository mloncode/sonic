package main

import (
	"github.com/MLonCode/sonic"
	"github.com/MLonCode/sonic/src/sound"
	"github.com/hypebeast/go-osc/osc"
)

func main() {
	m1 := sound.CreateMarkov("song1.midi")
	m2 := sound.CreateMarkov("song2.midi")

	client := osc.NewClient("localhost", 4559)

	oldChanges := sound.NewSequence("prophet", 100.1, 0.10, 0.05,
		sonic.ConvertMarkov(m1, sonic.File1.Old))
	newChanges := sound.NewSequence("prophet", 100.1, 0.10, 0.05,
		sonic.ConvertMarkov(m2, sonic.File1.New))

	oldChanges.Play(client)
	newChanges.Play(client)
}
