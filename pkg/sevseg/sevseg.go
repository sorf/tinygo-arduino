package sevseg

import (
	m "machine"

	"github.com/sorf/tinygo-arduino/pkg/common"
)

var (
	// Bitmask for each segment
	segments = []byte{
		0b01000000, // A
		0b00100000, // B
		0b00010000, // C
		0b00001000, // D
		0b00000100, // E
		0b00000010, // F
		0b00000001, // G
		0b10000000, // DP
	}

	// Supported characters
	Characters = []byte{
		'0',
		'1',
		'2',
		'3',
		'4',
		'5',
		'6',
		'7',
		'8',
		'9',
		'A',
		'B',
		'C',
		'D',
		'E',
		'F',
		'G',
		'H',
		'I',
		'.',
	}

	// Map from supported characters to the segments that represent it
	charactersMap = []byte{
		//DP
		// ABCDEFG
		0b01111110, // 0
		0b00110000, // 1
		0b01101101, // 2
		0b01111001, // 3
		0b00110011, // 4
		0b01011011, // 5
		0b01011111, // 6
		0b01110000, // 7
		0b01111111, // 8
		0b01111011, // 9
		0b01110111, // A
		0b00011111, // b/B
		0b01001110, // C
		0b00111101, // d/D
		0b01001111, // E
		0b01000111, // F
		0b01011111, // G
		0b00110111, // H
		0b00110000, // I
		0b10000000, // .
	}

	FirstLetter = byte('A')
	LastLetter  = byte('I')
)

type Error uint8

const (
	_                      = iota
	ErrWrongPinCount Error = iota
	ErrCharacterNotSupported
)

func (err Error) Error() string {
	switch err {
	case ErrWrongPinCount:
		return "invalid number of pins, 8 are required (A-to-G and DP)"
	case ErrCharacterNotSupported:
		return "character not supported"
	default:
		return "unspecified error"
	}
}

// Reference:
// https://docs.wokwi.com/parts/wokwi-7segment
// https://cdn.sparkfun.com/datasheets/Components/LED/YSD-160AR4B-8.pdf
type Device struct {
	pins [8]m.Pin
}

// Creates a new seven-segment device.
//
// Pins & Display:
//
//	10   9   CA   7   6
//	+-----------------+
//	|      -A(7)-     |
//	||F(9)       B(6)||
//	|     -G(10)-     |
//	||E(1)       C(4)||
//	|      -D(2)-     |
//	|           DP(5).|
//	+-----------------+
//	 1   2   CA   4   5
func NewSevSeq(p1, p2, p4, p5, p6, p7, p9, p10 m.Pin) (Device, error) {
	return NewSevSeqPins([]m.Pin{p1, p2, p4, p5, p6, p7, p9, p10})
}

// NewSevSeqPins is a variation of NewSevSeq with the pins provided as an array.
func NewSevSeqPins(pins []m.Pin) (Device, error) {
	if len(pins) != 8 {
		return Device{}, ErrWrongPinCount
	}
	// convert from pin to segment
	return Device{pins: [8]m.Pin{
		pins[5], // A
		pins[4], // B
		pins[2], // C
		pins[1], // D
		pins[0], // E
		pins[6], // F
		pins[7], // G
		pins[3], // DP
	}}, nil
}

// Configure initializes the device.
func (d *Device) Configure() error {
	for _, p := range d.pins {
		p.Configure(common.PinConfigOutput)
		p.High()
	}
	return nil
}

// Clear clears the display.
func (d *Device) Clear() {
	for _, p := range d.pins {
		p.High()
	}
}

// DisplayJustDot displays just the dot.
func (d *Device) DisplayJustDot() error {
	return d.DisplayDot('.', false)
}

// Display displays a character from the supported Characters set.
func (d *Device) Display(c byte) error {
	return d.DisplayDot(c, false)
}

// DisplayWithDot displays a character from the supported Characters set with a dot.
func (d *Device) DisplayWithDot(c byte) error {
	return d.DisplayDot(c, true)
}

// DisplayDot same as Display but with or without a dot.
func (d *Device) DisplayDot(c byte, dot bool) error {
	for i, sc := range Characters {
		if sc == c {
			d.display(i, dot)
			return nil
		}
	}
	return ErrCharacterNotSupported
}

// DisplayHex displays a number (up to 15).
func (d *Device) DisplayHex(n uint8) error {
	return d.DisplayHexDot(n, false)
}

// DisplayHexWithDot displays a number (up to 15) with a dot.
func (d *Device) DisplayHexWithDot(n uint8) error {
	return d.DisplayHexDot(n, true)
}

// DisplayHexDot displays a number (up to 15) with or without a dot.
func (d *Device) DisplayHexDot(n uint8, dot bool) error {
	if n > 15 {
		return ErrCharacterNotSupported
	}
	d.display(int(n), dot)
	return nil
}

func (d *Device) display(index int, dot bool) {
	segmentBitset := charactersMap[index]
	if dot {
		segmentBitset |= segments[7]
	}
	for s, sb := range segments {
		d.pins[s].Set(segmentBitset&sb == 0) // reversed: Low means ON
	}
}
