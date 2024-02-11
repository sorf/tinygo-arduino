package main

import (
	"bytes"
	m "machine"
	"strconv"
	"time"

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

	println("Start")
	lcd.Print([]byte("Countdown (i2c):"))
	time.Sleep(time.Millisecond * 250)
	for i := 0; i < 5; i++ {
		lcd.BacklightOn(false)
		lcd.Print([]byte("BacklightOn off")) // debug wokwi?
		time.Sleep(time.Millisecond * 250)
		lcd.BacklightOn(true)
		lcd.Print([]byte("BacklightOn on"))
		time.Sleep(time.Millisecond * 250)
	}
	var b bytes.Buffer
	b.Grow(16)
	for {
		for i := 12; i > 0; i-- {
			println(i)
			lcd.SetCursor(0, 1)
			b.Reset()
			b.Write([]byte("          "))
			if i < 10 {
				b.WriteByte('0')
			}
			b.Write([]byte(strconv.Itoa(i)))
			lcd.Print(b.Bytes())
			time.Sleep(1 * time.Second)
		}
	}
}
