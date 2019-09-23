package sound

import (
	"os"
	"sort"
    "gitlab.com/gomidi/midi/mid"
)

func midiLoad(f string) []int {
	j, err := os.Open(f)
	if err != nil {
		panic(err.Error())
	}

	defer j.Close()

	var notes []int
	rd := mid.NewReader(mid.NoLogger())
	rd.Msg.Channel.NoteOn = func(p *mid.Position, channel, key, vel uint8) {
		notes = append(notes, int(key))
	}

	if err := rd.ReadAllSMF(j); err != nil {
		panic(err.Error())
	}

	return notes
}

type Transitions map[int]map[int]uint32
type Markov map[int][]probability

// Returns a Markov chain
func NewMarkov(f string) Markov {
	notes := midiLoad(f)
	chain := make(Transitions)

	// Parse the transitions, i.e. store for a state (note identifier in midi)
	// the number of transitions to each other note, as observed in the midi file
	previous := -1
	for _, n := range notes {
		if previous != -1 {
			if chain[previous] == nil {
				chain[previous] = make(map[int]uint32)
			}

			chain[previous][n]++
		}

		previous = n
	}

	m := make(Markov)
	for k, c := range chain {
		p := getProbabilities(c)
		m[k] = p
	}

	return m
}

type probability struct {
	Key   int
	Value uint32
}

const probMax uint32 = (1 << 32) - 1

func getProbabilities(probs map[int]uint32) []probability {
	var max uint32
	var keys []int

	for n, p := range probs {
		keys = append(keys, n)
		max += p
	}

	if max == 0 {
		return nil
	}

	// keys are sorted as map order is not stable
	sort.Ints(keys)

	// probabilities are scaled from 0 to max(uint32)
	scale := float64(probMax) / float64(max)
	var pList []probability
	var previous uint32
	for _, k := range keys {
		t := uint32(float64(probs[k]) * scale)
		// take care of overflow
		if t > probMax - previous {
			t = probMax
		} else {
			t += previous
		}

		previous = t

		p := probability{k, t}
		pList = append(pList, p)
	}

	return pList
}

func (m Markov) Get(prev int, num uint32) int {
	chain := m[prev]
	if chain == nil {
		return m.Rand(num)
	}

	for _, p := range chain {
		if p.Value >= num {
			return p.Key
		}
	}

	return m.Rand(num)
}

func (m Markov) Rand(num uint32) int {
	// if the chain is not initialized return a valid note
	if m == nil {
		return 24 * 3 * 12 // C3
	}

	var keys []int

	for n := range m {
		keys = append(keys, n)
	}

	return keys[num%uint32(len(keys))]
}
