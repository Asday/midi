package meta

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/lib"
)

/* http://www.somascape.org/midi/tech/mfile.html
Key Signature

FF 59 02 sf mi

sf is a byte specifying the number of flats (-ve) or sharps (+ve) that identifies the key signature (-7 = 7 flats, -1 = 1 flat, 0 = key of C, 1 = 1 sharp, etc).
mi is a byte specifying a major (0) or minor (1) key.

For a format 1 MIDI file, Key Signature Meta events should only occur within the first MTrk chunk.

*/

const (
	degreeC  = 0
	degreeCs = 1
	degreeDf = degreeCs
	degreeD  = 2
	degreeDs = 3
	degreeEf = degreeDs
	degreeE  = 4
	degreeF  = 5
	degreeFs = 6
	degreeGf = degreeFs
	degreeG  = 7
	degreeGs = 8
	degreeAf = degreeGs
	degreeA  = 9
	degreeAs = 10
	degreeBf = degreeAs
	degreeB  = 11
	degreeCf = degreeB
)

// Supplied to KeySignature
const (
	majorMode = 0
	minorMode = 1
)

type KeySignature struct {
	Key     uint8
	IsMajor bool
	Num     uint8
	//	SharpsOrFlats int8
	IsFlat bool
}

/*
// NewKeySignature returns a key signature event.
// key is the key of the scale (C=0 add the corresponding number of semitones). ismajor indicates whether it is a major or minor scale
// num is the number of accidentals. isflat indicates whether the accidentals are flats or sharps
func NewKeySignature(key uint8, ismajor bool, num uint8, isflat bool) KeySignature {
	return KeySignature{Key: key, IsMajor: ismajor, Num: num, IsFlat: isflat}
}
*/

func (m KeySignature) Raw() []byte {
	mi := int8(0)
	if !m.IsMajor {
		mi = 1
	}
	sf := int8(m.Num)

	if m.IsFlat {
		sf = sf * (-1)
	}

	return (&metaMessage{
		Typ:  byteKeySignature,
		Data: []byte{byte(sf), byte(mi)},
	}).Bytes()
}

func (m KeySignature) String() string {
	return fmt.Sprintf("%T: %s", m, m.Text())
}

func (m KeySignature) Note() (note string) {
	switch m.Key {
	case degreeC:
		note = "C"
	case degreeD:
		note = "D"
	case degreeE:
		note = "E"
	case degreeF:
		note = "F"
	case degreeG:
		note = "G"
	case degreeA:
		note = "A"
	case degreeB:
		note = "B"
	case degreeCs:
		note = "C♯"
		if m.IsFlat {
			note = "D♭"
		}
	case degreeDs:
		note = "D♯"
		if m.IsFlat {
			note = "E♭"
		}
	case degreeFs:
		note = "F♯"
		if m.IsFlat {
			note = "G♭"
		}
	case degreeGs:
		note = "G♯"
		if m.IsFlat {
			note = "A♭"
		}
	case degreeAs:
		note = "A♯"
		if m.IsFlat {
			note = "B♭"
		}
	default:
		panic("unreachable")
	}

	return
}

func (m KeySignature) Text() string {
	if m.IsMajor {
		return m.Note() + " maj."
	}

	return m.Note() + " min."
}

// Taking a signed number of sharps or flats (positive for sharps, negative for flats) and a mode (0 for major, 1 for minor)
// decide the key signature.
func keyFromSharpsOrFlats(sharpsOrFlats int8, mode uint8) uint8 {
	tmp := int(sharpsOrFlats * 7)

	// Relative Minor.
	if mode == minorMode {
		tmp -= 3
	}

	// Clamp to Octave 0-11.
	for tmp < 0 {
		tmp += 12
	}

	return uint8(tmp % 12)
}

func (m KeySignature) readFrom(rd io.Reader) (Message, error) {

	// fmt.Println("Key signature")
	// TODO TEST
	var sharpsOrFlats int8
	var mode uint8

	length, err := lib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 2 {
		err = lib.UnexpectedMessageLengthError("KeySignature expected length 2")
		return nil, err
	}

	// Signed int, positive is sharps, negative is flats.
	var b byte
	b, err = lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	sharpsOrFlats = int8(b)

	// Mode is Major or Minor.
	mode, err = lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	num := sharpsOrFlats
	if num < 0 {
		num = num * (-1)
	}

	key := keyFromSharpsOrFlats(sharpsOrFlats, mode)

	return KeySignature{
		Key:     key,
		Num:     uint8(num),
		IsMajor: mode == majorMode,
		IsFlat:  sharpsOrFlats < 0,
	}, nil

}

func (m KeySignature) meta() {}
