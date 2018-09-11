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

type Sequence struct {
	Instrument string
	Cutoff     float64
	Attack     float64
	Release    float64
	Notes      []Note
}

func NewSequence(i string, c, a, r float64, n []Note) Sequence {
	return Sequence{i, c, a, r, n}
}

func (s *Sequence) Play(client *osc.Client) {
	for _, n := range s.Notes {
		addr := fmt.Sprintf("/trigger/%s", s.Instrument)
		msg := osc.NewMessage(addr)

		// message: note, duration, cutoff, release

		msg.Append(n.Note)
		msg.Append(n.Duration)
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
