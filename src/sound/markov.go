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

type Markov map[string]map[string]int

func CreateMarkov(f string) Markov {
	notes := midiLoad(f)
	m := make(Markov)

	previous := ""
	for _, n := range notes {
		println(n.Name)
		if previous != "" {
			if m[previous] == nil {
				m[previous] = make(map[string]int)
			}

			m[previous][n.Name] += 1
		}

		previous = n.Name
	}

	return m
}

type probability struct {
	Name  string
	Value int
}

func (m Markov) Get(prev string, rnd int) string {
	probs := m[prev]

	max := 0
	var keys []string

	for n, p := range probs {
		keys = append(keys, n)
		max += p
	}

	if max == 0 {
		return m.Rand(rnd)
	}

	sort.Strings(keys)

	var pList []probability
	for _, k := range keys {
		t := probs[k]
		p := probability{k, t}
		pList = append(pList, p)
	}

	num := rnd % max
	cur := 0
	for _, p := range pList {
		cur += p.Value
		println(max, cur, num)
		if cur >= num {
			return p.Name
		}
	}

	return ""
}

func (m Markov) Rand(rnd int) string {
	var keys []string

	for n, _ := range m {
		keys = append(keys, n)
	}

	return keys[rnd%len(keys)]

}
