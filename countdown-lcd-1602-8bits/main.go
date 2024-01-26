package main

import (
	"bytes"
	m "machine"
	"strconv"
	"time"

	"tinygo.org/x/drivers/hd44780"
)

func main() {
	m.Serial.Configure(m.UARTConfig{BaudRate: 9600})
	lcd, err := hd44780.NewGPIO8Bit([]m.Pin{m.D7, m.D8, m.D9, m.D10, m.D5, m.D4, m.D3, m.D2},
		m.D11, m.D12, m.NoPin)
	if err != nil {
		println("error: create LCD", err.Error())
		return
	}
	if err := lcd.Configure(hd44780.Config{Width: 16, Height: 2}); err != nil {
		println("error: configure LCD", err.Error())
		return

	}

	println("Start")
	lcd.Write([]byte("Countdown (8b):"))
	lcd.Display()
	time.Sleep(1 * time.Second)

	var b bytes.Buffer
	b.Grow(16)
	for {
		for i := 12; i > 0; i-- {
			lcd.SetCursor(0, 1)
			b.Reset()
			b.Write([]byte("          "))
			if i < 10 {
				b.WriteByte('0')
			}
			b.Write([]byte(strconv.Itoa(i)))
			lcd.Write(b.Bytes())
			lcd.Display()
			time.Sleep(1 * time.Second)
		}
	}
}
