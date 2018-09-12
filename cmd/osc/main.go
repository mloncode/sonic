package main

import (
	"github.com/MLonCode/sonic"
	"github.com/MLonCode/sonic/src/sound"
	"github.com/hypebeast/go-osc/osc"
)

func main() {
	client := osc.NewClient("localhost", 4559)

	amajor := sound.NewScale(sound.AMajor, 4, 2)
	cminor := sound.NewScale(sound.CMinor, 4, 2)

	oldChanges := sound.NewSequence("prophet", 80.1, 0.10, 0.05,
		sonic.Convert(amajor, sonic.File1.Old))
	newChanges := sound.NewSequence("prophet", 100.1, 0.10, 0.05,
		sonic.Convert(cminor, sonic.File1.New))

	oldChanges.Play(client)
	newChanges.Play(client)
}
