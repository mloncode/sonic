package sound

import "fmt"

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

func (c *Scale) Get(number int) string {
	return c.Notes[number%len(c.Notes)]
}

var (
	EMajor = []string{"E", "Fs", "Gs", "A", "B", "Cs", "Ds"}
	AMajor = []string{"A", "B", "Cs", "D", "E", "Fs", "Gs"}
)
