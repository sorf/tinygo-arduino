package main

import (
	m "machine"
	"time"

	"github.com/sorf/tinygo-arduino/pkg/sevseg"
)

func require(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	m.Serial.Configure(m.UARTConfig{BaudRate: 9600})

	d, err := sevseg.NewSevSeq(m.D2, m.D3, m.D4, m.D5, m.D6, m.D7, m.D8, m.D9)
	if err != nil {
		println("error: create seven-segment", err.Error())
		return
	}
	if err := d.Configure(); err != nil {
		println("error: configure seven-segment", err.Error())
		return
	}

	dot := false
	for {
		for i := 0; i < 16; i++ {
			require(d.DisplayHexDot(uint8(i), dot))
			time.Sleep(time.Second)
		}
		require(d.DisplayJustDot())
		time.Sleep(time.Second)
		d.Clear()
		time.Sleep(time.Second)

		for c := sevseg.FirstLetter; c <= sevseg.LastLetter; c++ {
			require(d.DisplayDot(c, dot))
			time.Sleep(time.Second)
		}

		dot = !dot
	}
}
