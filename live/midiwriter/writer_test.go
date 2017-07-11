package midiwriter

import (
	"bytes"
	"fmt"
	"github.com/gomidi/midi/messages/channel"
	"testing"
)

func TestRunningStatus(t *testing.T) {

	var bf bytes.Buffer

	wr := New(&bf)

	wr.Write(channel.Ch0.NoteOn(50, 33))
	wr.Write(channel.Ch0.NoteOff(50))

	expected := "90 32 21 32 00"

	if got, want := fmt.Sprintf("% X", bf.Bytes()), expected; got != want {
		t.Errorf("got:\n%#v\nwanted:\n%#v\n\n", got, want)
	}
}

func TestNoRunningStatus(t *testing.T) {

	var bf bytes.Buffer

	wr := New(&bf, NoRunningStatus())

	wr.Write(channel.Ch0.NoteOn(50, 33))
	wr.Write(channel.Ch0.NoteOff(50))

	expected := "90 32 21 90 32 00"

	if got, want := fmt.Sprintf("% X", bf.Bytes()), expected; got != want {
		t.Errorf("got:\n%#v\nwanted:\n%#v\n\n", got, want)
	}
}
