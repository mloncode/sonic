package main

import (
	"math/rand"

	"github.com/MLonCode/sonic/src/sound"
	"github.com/hypebeast/go-osc/osc"
)

func main() {
	client := osc.NewClient("localhost", 4559)

	scale := sound.NewScale(sound.AMajor, 4, 2)

	var notes []sound.Note

	for i := 0; i < 10; i++ {
		n := rand.Intn(100)
		note := scale.Get(n)
		duration := rand.Float64() * 0.5

		notes = append(notes, sound.Note{Note: note, Duration: duration})
	}

	sequence := sound.NewSequence("prophet", 50.1, 0.1, 0.05, notes)
	sequence.Play(client)
}
