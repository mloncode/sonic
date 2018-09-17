package sonic

import (
	"hash/fnv"
	"math"

	"github.com/MLonCode/sonic/src/sound"
)

const logarithmic = false

func Convert(scale sound.Scale, nodes []sonicNode) []sound.Note {
	var max uint32
	for _, n := range nodes {
		if n.Lenght > max {
			max = n.Lenght
		}
	}

	notes := make([]sound.Note, len(nodes))
	hash := fnv.New32a()

	for i, n := range nodes {
		duration := float64(n.Lenght)
		duration = toLog(float64(max), duration) * 0.25

		hash.Reset()
		hash.Write([]byte(n.Token))
		note := scale.Get(hash.Sum32())

		notes[i] = sound.Note{Note: int64(note), Duration: duration}
	}

	return notes
}

func ConvertMarkov(m sound.Markov, nodes []sonicNode) []sound.Note {
	var max uint32
	for _, n := range nodes {
		if n.Lenght > max {
			max = n.Lenght
		}
	}

	notes := make([]sound.Note, len(nodes))
	hash := fnv.New32a()

	last := m.Rand(uint32(len(nodes)))
	for i, n := range nodes {
		var duration float64
		if logarithmic {
			duration = float64(n.Lenght)
			duration = toLog(float64(max), duration) * 0.25
		} else {
			duration = float64(n.Lenght) / float64(max)
		}

		hash.Reset()
		hash.Write([]byte(n.Token))
		note := m.Get(last, hash.Sum32())
		last = note

		notes[i] = sound.Note{Note: int64(note), Duration: duration}
	}

	return notes
}

func toLog(max, value float64) float64 {
	return (math.Log(value) - math.Log(0.1)) / (math.Log(max) - math.Log(0.1))
}
