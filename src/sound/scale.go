package sound

import (
	"fmt"
	"strconv"
)

type Scale struct {
	Notes []string
}

func NewScale(notes []string, base, width int) Scale {
	scale := make([]string, len(notes)*width)

	for i := range scale {
		n := notes[i%len(notes)]
		o := i/len(notes) + base

		scale[i] = fmt.Sprintf("%s%d", n, o)
	}

	return Scale{scale}
}

func (c *Scale) Get(number uint32) int {
	return keys[c.Notes[number%uint32(len(c.Notes))]]
}

var (
	EMajor = []string{"E", "Fs", "Gs", "A", "B", "Cs", "Ds"}
	AMajor = []string{"A", "B", "Cs", "D", "E", "Fs", "Gs"}
	CMinor = []string{"C", "D", "Es", "F", "G", "As", "Bs"}
)

var keyBase = map[string]int{
	"C": 24,
	"D": 26,
	"E": 28,
	"F": 29,
	"G": 31,
	"A": 33,
	"B": 35,
}

var keys map[string]int

func init() {
	keys = make(map[string]int)
	for octave := 0; octave < 9; octave++ {
		octaveStr := ""
		if octave != 3 {
			octaveStr = strconv.Itoa(octave + 1)
		}
		for bk, bv := range keyBase {
			key := bv + (octave * 12)
			keys[bk+"b"+octaveStr] = key - 1
			keys[bk+octaveStr] = key
			keys[bk+"s"+octaveStr] = key + 1
		}
	}
}
