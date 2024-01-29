package main

import (
	m "machine"
	"time"

	"github.com/sorf/tinygo-arduino/pkg/sevseg"
)

func main() {
	m.Serial.Configure(m.UARTConfig{BaudRate: 9600})

	d1 := sevseg.NewSevSeq(m.ADC0, m.ADC1, m.ADC2, m.ADC3, m.ADC4, m.ADC5, m.D12, m.D13)
	d1.Configure()
	d2 := sevseg.NewSevSeq(m.D2, m.D3, m.D4, m.D5, m.D6, m.D7, m.D8, m.D9)
	d2.Configure()

	chars1 := make([]byte, 0, 10+len(sevseg.SpecialCharacters))
	for d := byte('0'); d <= '9'; d++ {
		chars1 = append(chars1, d)
	}
	for _, s := range sevseg.SpecialCharacters {
		chars1 = append(chars1, s)
	}

	chars2 := make([]byte, 0, int('z'-'a'))
	for l := byte('a'); l <= 'z'; l++ {
		chars2 = append(chars2, l)
	}

	dot := false
	counter := 0
	maxCounter := len(chars2)
	for {
		if counter == maxCounter {
			d1.Clear()
			d2.Clear()
			counter = 0
			dot = !dot
		} else {
			d1.DisplayDot(chars1[counter%len(chars1)], dot)
			d2.DisplayDot(chars2[counter%len(chars2)], dot)
			counter++
		}
		time.Sleep(time.Second)
	}
}
