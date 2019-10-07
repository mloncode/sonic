package sound

import (
	"log"
	"time"

	"gitlab.com/gomidi/midi/mid"
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

func (s *Sequence) Play(out mid.Out) {
	out.Open()
	wr := mid.ConnectOut(out)

	for _, n := range s.Notes {
		log.Print("note", n.Note)

		note := uint8(n.Note)
		wr.NoteOn(note, 100)
		time.Sleep(time.Duration(n.Duration * float64(time.Second)))
		wr.NoteOff(note)
	}

	out.Close()
}

func (a Notes) Len() int           { return len(a) }
func (a Notes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Notes) Less(i, j int) bool { return a[i].Duration > a[j].Duration }
