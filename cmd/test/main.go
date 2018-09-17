package main

import (
	"log"

	"github.com/MLonCode/sonic"
	"github.com/MLonCode/sonic/src/sound"
	"github.com/rakyll/portmidi"
)

func main() {
	if err := portmidi.Initialize(); err != nil {
		log.Fatal("can't initializer portmidi", err)
	}
	defer portmidi.Terminate()

	if portmidi.CountDevices() == 0 {
		log.Fatal("no midi devices")
	}

	m1 := sound.NewMarkov("song1.midi")
	m2 := sound.NewMarkov("song2.midi")

	oldChanges := sound.NewSequence("prophet", sonic.ConvertMarkov(m1, sonic.File1.Old))
	newChanges := sound.NewSequence("prophet", sonic.ConvertMarkov(m2, sonic.File1.New))

	deviceID := portmidi.DefaultOutputDeviceID()

	oldChanges.Play(deviceID)
	newChanges.Play(deviceID)
}
