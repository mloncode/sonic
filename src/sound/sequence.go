package sound

import (
	"log"
	"time"

	"github.com/rakyll/portmidi"
)

type Note struct {
	Note     int64
	Duration float64
}

type Notes []Note

type Sequence struct {
	Instrument string
	Notes      Notes
}

func NewSequence(i string, n []Note) Sequence {
	return Sequence{i, n}
}

func (s *Sequence) Play(id portmidi.DeviceID) {
	out, err := portmidi.NewOutputStream(id, 1024, 0)
	if err != nil {
		log.Fatal(err)
	}

	for _, n := range s.Notes {
		log.Print("note", n.Note)

		out.WriteShort(0x90, n.Note, 100)
		time.Sleep(time.Duration(n.Duration * float64(time.Second)))
		out.WriteShort(0x80, n.Note, 100)
	}

	out.Close()
}

func (a Notes) Len() int           { return len(a) }
func (a Notes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Notes) Less(i, j int) bool { return a[i].Duration > a[j].Duration }
