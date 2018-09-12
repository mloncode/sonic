package sound

import (
	"fmt"
	"time"

	"github.com/hypebeast/go-osc/osc"
)

type Note struct {
	Note     string
	Duration float64
}

type Notes []Note

type Sequence struct {
	Instrument string
	Cutoff     float64
	Attack     float64
	Release    float64
	Notes      Notes
}

func NewSequence(i string, c, a, r float64, n []Note) Sequence {
	return Sequence{i, c, a, r, n}
}

func (s *Sequence) Play(client *osc.Client) {
	// sort.Sort(s.Notes)
	for _, n := range s.Notes {
		addr := fmt.Sprintf("/trigger/%s", s.Instrument)
		msg := osc.NewMessage(addr)

		// message: note, duration, cutoff, release

		msg.Append(n.Note)
		msg.Append(n.Duration / 1.3)
		msg.Append(s.Cutoff)
		msg.Append(s.Attack)
		msg.Append(s.Release)
		err := client.Send(msg)

		println("note", n.Note)
		if err != nil {
			println("err", err.Error())
		}

		time.Sleep(time.Duration(n.Duration * float64(time.Second)))
	}
}

/*

sonic pi script:


live_loop :foo do
  use_real_time
  note, duration, cutoff, attack, release = sync "/osc/trigger/prophet"
  synth :prophet,
    note: note,
    cutoff: cutoff,
    sustain: duration,
    release: release,
    attack: attack,
    amp: 8
end

*/

func (a Notes) Len() int           { return len(a) }
func (a Notes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Notes) Less(i, j int) bool { return a[i].Duration > a[j].Duration }
