package main

import (
	m "machine"
	"time"

	"github.com/sorf/tinygo-arduino/pkg/sevseg"
	"tinygo.org/x/drivers/hd44780i2c"
)

func main() {
	m.Serial.Configure(m.UARTConfig{BaudRate: 9600})

	i2c := m.I2C0
	err := i2c.Configure(m.I2CConfig{})
	if err != nil {
		println("could not configure I2C:", err)
		return
	}

	lcd := hd44780i2c.New(i2c, 0x27)
	if err != nil {
		println("error: create LCD", err.Error())
		return
	}
	if err := lcd.Configure(hd44780i2c.Config{Width: 16, Height: 2}); err != nil {
		println("error: configure LCD", err.Error())
		return
	}

	d1 := sevseg.NewSevSeq(m.ADC0, m.ADC1, m.ADC2, m.ADC3, m.D10, m.D11, m.D12, m.D13)
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

	chars2 := make([]byte, 0, int('Z'-'A'))
	for l := byte('A'); l <= 'Z'; l++ {
		chars2 = append(chars2, l)
	}

	dot := false
	counter := 0
	maxCounter := len(chars2)
	for {
		if counter == maxCounter {
			d1.Clear()
			d2.Clear()
			lcd.ClearDisplay()
			counter = 0
			dot = !dot
		} else {
			c1 := chars1[counter%len(chars1)]
			c2 := chars2[counter%len(chars2)]
			d := byte(' ')
			if dot {
				d = '.'
			}

			// Display on the LCD first due to its latency
			lcd.SetCursor(7, 0)
			lcd.Print([]byte{c1, d})
			lcd.SetCursor(7, 1)
			lcd.Print([]byte{c2, d})
			// Then the 7-segments
			d1.DisplayDot(c1, dot)
			d2.DisplayDot(c2, dot)
			counter++
		}
		time.Sleep(1500 * time.Millisecond)
	}
}
