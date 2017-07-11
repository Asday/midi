package midireader_test

import (
	"bytes"
	"fmt"

	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/realtime"
	"github.com/gomidi/midi/midireader"
	"github.com/gomidi/midi/midiwriter"

	"github.com/gomidi/midi"
)

func Example() {
	var bf bytes.Buffer

	wr := midiwriter.New(&bf)
	wr.Write(channel.Ch2.NoteOn(65, 90))
	wr.Write(realtime.Reset)
	wr.Write(channel.Ch2.NoteOff(65))

	rthandler := func(m realtime.Message) {
		fmt.Printf("Realtime: %s\n", m)
	}

	rd := midireader.New(bytes.NewReader(bf.Bytes()), rthandler)

	var m midi.Message
	var err error

	for {
		m, err = rd.Read()

		// breaking at least with io.EOF
		if err != nil {
			break
		}

		switch v := m.(type) {
		case channel.NoteOn:
			fmt.Printf("NoteOn at channel %v: key: %v velocity: %v\n", v.Channel(), v.Key(), v.Velocity())
		case channel.NoteOff:
			fmt.Printf("NoteOff at channel %v: key: %v\n", v.Channel(), v.Key())
		}

	}

	// Output: NoteOn at channel 2: key: 65 velocity: 90
	// Realtime: Reset
	// NoteOff at channel 2: key: 65
}
