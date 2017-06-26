package meta

import (
	"io"
)

const (
	// End of track
	// the handler is supposed to keep track of the current track

	byteEndOfTrack            = byte(0x2F)
	byteSequenceNumber        = byte(0x00)
	byteText                  = byte(0x01)
	byteCopyright             = byte(0x02)
	byteSequence              = byte(0x03)
	byteTrackInstrument       = byte(0x04)
	byteLyric                 = byte(0x05)
	byteMarker                = byte(0x06)
	byteCuePoint              = byte(0x07)
	byteMIDIChannel           = byte(0x20)
	byteDevicePort            = byte(0x9)
	byteMIDIPort              = byte(0x21)
	byteTempo                 = byte(0x51)
	byteTimeSignature         = byte(0x58)
	byteKeySignature          = byte(0x59)
	byteSequencerSpecificInfo = byte(0x7F)
)

var metaMessages = map[byte]Message{
	byteEndOfTrack:            EndOfTrack,
	byteSequenceNumber:        SequenceNumber(0),
	byteText:                  Text(""),
	byteCopyright:             Copyright(""),
	byteSequence:              Sequence(""),
	byteTrackInstrument:       TrackInstrument(""),
	byteLyric:                 Lyric(""),
	byteMarker:                Marker(""),
	byteCuePoint:              CuePoint(""),
	byteMIDIChannel:           MIDIChannel(0),
	byteDevicePort:            DevicePort(""),
	byteMIDIPort:              MIDIPort(0),
	byteTempo:                 Tempo(0),
	byteTimeSignature:         TimeSignature{},
	byteKeySignature:          KeySignature{},
	byteSequencerSpecificInfo: nil, // SequencerSpecificInfo
}

// Reader reads a Meta Message
type Reader interface {
	// Read reads a single Meta Message.
	// It may just be called once per Reader. A second call returns io.EOF
	Read() (Message, error)
}

// NewReader returns a reader that can read a single Meta Message
// Read may just be called once per Reader. A second call returns io.EOF
func NewReader(input io.Reader, typ byte) Reader {
	return &reader{input, typ, false}
}

type reader struct {
	input io.Reader
	typ   byte
	done  bool
}

// Read may just be called once per Reader. A second call returns io.EOF
func (r *reader) Read() (Message, error) {
	if r.done {
		return nil, io.EOF
	}

	m := metaMessages[r.typ]
	if m == nil {
		m = Undefined{Typ: r.typ}
	}

	return m.readFrom(r.input)
}
