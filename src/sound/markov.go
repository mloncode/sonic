package sound

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type midiNote struct {
	Name     string
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

	replacer := strings.NewReplacer("#", "b")
	var notes []midiNote
	for _, t := range m.Tracks {
		for _, n := range t.Notes {
			note := midiNote{
				Name:     replacer.Replace(n.Name),
				Duration: n.Duration,
			}

			notes = append(notes, note)
		}
	}

	return notes
}

type markov map[string]map[string]uint32
type Markov map[string][]probability

func NewMarkov(f string) Markov {
	notes := midiLoad(f)
	chain := make(markov)

	previous := ""
	for _, n := range notes {
		if previous != "" {
			if chain[previous] == nil {
				chain[previous] = make(map[string]uint32)
			}

			chain[previous][n.Name]++
		}

		previous = n.Name
	}

	m := make(Markov)
	for k, c := range chain {
		p := getProbabilities(c)
		m[k] = p
	}

	return m
}

type probability struct {
	Name  string
	Value uint32
}

const probMax uint32 = (1 << 32) - 1

func getProbabilities(probs map[string]uint32) []probability {
	var max uint32
	var keys []string

	for n, p := range probs {
		keys = append(keys, n)
		max += p
	}

	if max == 0 {
		return nil
	}

	// keys are sorted as map order is not stable
	sort.Strings(keys)

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

func (m Markov) Get(prev string, num uint32) string {
	chain := m[prev]
	if chain == nil {
		return m.Rand(num)
	}

	for _, p := range chain {
		if p.Value >= num {
			return p.Name
		}
	}

	return m.Rand(num)
}

func (m Markov) Rand(num uint32) string {
	// if the chain is not initialized return a valid note
	if m == nil {
		return "C3"
	}

	var keys []string

	for n, _ := range m {
		keys = append(keys, n)
	}

	return keys[num%uint32(len(keys))]
}
