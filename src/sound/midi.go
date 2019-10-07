package sound

import (
	"strings"
	"errors"
	"gitlab.com/gomidi/midi/mid"
	driver "gitlab.com/gomidi/portmididrv"
)


func MidiOut() (mid.Out, error) {
	drv, err := driver.New()
	defer drv.Close()
	if err != nil {
		return nil, errors.New("could not initialize the midi driver")
	}

	outs, _ := drv.Outs()
	var out mid.Out
	out = nil

	for _, current := range(outs) {
		// Check a midi out port that is not named 'Midi Through Port-0'
		// FluidSynth port would be named 'Synth input port' for example
		if !strings.Contains(current.String(), "Midi Through Port-") {
			out = current
			break;
		}
	}

	if out == nil {
		return nil, errors.New("could not find suitable out midi port")
	}

	return out, nil
}
