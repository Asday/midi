package meta

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/lib"
)

type TrackInstrument string

func (m TrackInstrument) String() string {
	return fmt.Sprintf("%T: %#v", m, string(m))
}

func (m TrackInstrument) Raw() []byte {
	return (&metaMessage{
		Typ:  byteTrackInstrument,
		Data: []byte(m),
	}).Bytes()
}

func (m TrackInstrument) readFrom(rd io.Reader) (Message, error) {
	text, err := lib.ReadText(rd)

	if err != nil {
		return nil, err
	}

	return TrackInstrument(text), nil
}

func (m TrackInstrument) Text() string {
	return string(m)
}

func (m TrackInstrument) meta() {}
