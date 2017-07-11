package midihandler_test

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gomidi/midi/midihandler"
	"github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/midimessage/meta"
	"github.com/gomidi/midi/smf/smfwriter"
)

func mkMIDI() io.Reader {
	var bf bytes.Buffer

	wr := smfwriter.New(&bf)
	wr.Write(channel.Ch2.NoteOn(65, 90))
	wr.SetDelta(2)
	wr.Write(channel.Ch2.NoteOff(65))
	wr.Write(meta.EndOfTrack)

	return bytes.NewReader(bf.Bytes())
}

func Example() {

	hd := midihandler.New(midihandler.NoLogger())

	// set the functions for the messages you are interested in
	hd.Message.Channel.NoteOn = func(p *midihandler.SMFPosition, channel, pitch, vel uint8) {
		fmt.Printf("[%v] NoteOn at channel %v: pitch %v velocity: %v\n", p.Delta, channel, pitch, vel)
	}

	hd.Message.Channel.NoteOff = func(p *midihandler.SMFPosition, channel, pitch, vel uint8) {
		fmt.Printf("[%v] NoteOff at channel %v: pitch %v velocity: %v\n", p.Delta, channel, pitch, vel)
	}

	hd.ReadSMF(mkMIDI())

	// Output: [0] NoteOn at channel 2: pitch 65 velocity: 90
	// [2] NoteOff at channel 2: pitch 65 velocity: 0

}
