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
		0b00011111, // b
		0b01001110, // C
		0b00111101, // d
		0b01001111, // E
		0b01000111, // F
		0b01011111, // G
		0b00110111, // H
		0b00110000, // I
		0b00111000, // J
		0b00110111, // K  Same as 'H'
		0b00001110, // L
		0b00000000, // M  N/A
		0b00010101, // n
		0b01111110, // O
		0b01100111, // P
		0b01110011, // q
		0b00000101, // r
		0b01011011, // S
		0b00001111, // t
		0b00111110, // U
		0b00111110, // V  Same as 'U'
		0b00000000, // W  N/A
		0b00110111, // X  Same as 'H'
		0b00111011, // y
		0b01101101, // Z  Same as '2'
		0b00000000, //    Space
		0b00000001, // -  Dash
		0b10000000, // .  Dot
		0b01100011, // *  Star
		0b00001000, // _  Underscore
	}

	SpecialCharacters = []byte{' ', '-', '.', '*', '_'}
)

const (
	indexDigits     = 0
	indexA          = indexDigits + 10
	indexZ          = indexA + 25
	indexSpace      = indexZ + 1
	indexDash       = indexSpace + 1
	indexDot        = indexDash + 1
	indexStar       = indexDot + 1
	indexUnderscore = indexStar + 1
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

// Seven-Segment device.
type Device struct {
	pins [8]m.Pin
}

// NewSevSeq creates a new seven-segment device.
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
//
// Reference:
// https://docs.wokwi.com/parts/wokwi-7segment
// https://cdn.sparkfun.com/datasheets/Components/LED/YSD-160AR4B-8.pdf
func NewSevSeq(p1, p2, p4, p5, p6, p7, p9, p10 m.Pin) Device {
	return Device{pins: [8]m.Pin{
		p7,  // A
		p6,  // B
		p4,  // C
		p2,  // D
		p1,  // E
		p9,  // F
		p10, // G
		p5,  // DP
	}}
}

// NewPinSet creates a new seven-segment device. from an array.
func NewSevSeqPins(p []m.Pin) (Device, error) {
	if len(p) != 8 {
		return Device{}, ErrWrongPinCount
	}
	return NewSevSeq(p[0], p[1], p[2], p[3], p[4], p[5], p[6], p[7]), nil
}

// Configure initializes the device.
func (d *Device) Configure() {
	for _, p := range d.pins {
		p.Configure(common.PinConfigOutput)
		p.High()
	}
}

// Clear clears the display.
func (d *Device) Clear() {
	for _, p := range d.pins {
		p.High()
	}
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
	if c >= '0' && c <= '9' {
		d.display(int(c-'0')+indexDigits, dot)
	} else if c >= 'a' && c <= 'z' {
		d.display(int(c-'a')+indexA, dot)
	} else if c >= 'A' && c <= 'Z' {
		d.display(int(c-'A')+indexA, dot)
	} else {
		switch c {
		case ' ':
			d.display(indexSpace, dot)
		case '-':
			d.display(indexDash, dot)
		case '.':
			d.display(indexDot, dot)
		case '*':
			d.display(indexStar, dot)
		case '_':
			d.display(indexUnderscore, dot)
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
