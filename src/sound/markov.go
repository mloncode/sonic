package sound

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
)

type midiNote struct {
	Key      int `json:"midi"`
	Duration float64
}

type midiTrack struct {
	Notes []midiNote
}

type midi struct {
	Tracks []midiTrack
}

func midiLoad(f string) []midiNote {
	j, err := os.Open(f)
	if err != nil {
		panic(err.Error())
	}

	defer j.Close()

	b, err := ioutil.ReadAll(j)
	if err != nil {
		panic(err.Error())
	}

	var m midi

	json.Unmarshal(b, &m)

	var notes []midiNote
	for _, t := range m.Tracks {
		for _, n := range t.Notes {
			note := midiNote{
				Key:      n.Key,
				Duration: n.Duration,
			}

			notes = append(notes, note)
		}
	}

	return notes
}

type markov map[int]map[int]uint32
type Markov map[int][]probability

func NewMarkov(f string) Markov {
	notes := midiLoad(f)
	chain := make(markov)

	previous := 0
	for _, n := range notes {
		if previous != 0 {
			if chain[previous] == nil {
				chain[previous] = make(map[int]uint32)
			}

			chain[previous][n.Key]++
		}

		previous = n.Key
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
		t := probs[k]
		if t == max {
			// maximum probability is manually set to overcome rounding errors
			t = probMax
		} else {
			t = uint32(float64(t)*scale) + previous
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
