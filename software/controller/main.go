package main

import (
	"fmt"
	"os"
	"time"

	"github.com/robfig/cron"
	"github.com/stianeikeland/go-rpio"
)

// The library use the raw BCM2835 pinouts, not the ports as they are mapped
// on the output pins for the raspberry pi, and not the wiringPi convention.
//             Rev 2 and 3 Raspberry Pi                        Rev 1 Raspberry Pi (legacy)
//   +-----+---------+----------+---------+-----+      +-----+--------+----------+--------+-----+
//   | BCM |   Name  | Physical | Name    | BCM |      | BCM | Name   | Physical | Name   | BCM |
//   +-----+---------+----++----+---------+-----+      +-----+--------+----++----+--------+-----+
//   |     |    3.3v |  1 || 2  | 5v      |     |      |     | 3.3v   |  1 ||  2 | 5v     |     |
//   |   2 |   SDA 1 |  3 || 4  | 5v      |     |      |   0 | SDA    |  3 ||  4 | 5v     |     |
//   |   3 |   SCL 1 |  5 || 6  | 0v      |     |      |   1 | SCL    |  5 ||  6 | 0v     |     |
//   |   4 | GPIO  7 |  7 || 8  | TxD     | 14  |      |   4 | GPIO 7 |  7 ||  8 | TxD    |  14 |
//   |     |      0v |  9 || 10 | RxD     | 15  |      |     | 0v     |  9 || 10 | RxD    |  15 |
//   |  17 | GPIO  0 | 11 || 12 | GPIO  1 | 18  |      |  17 | GPIO 0 | 11 || 12 | GPIO 1 |  18 |
//   |  27 | GPIO  2 | 13 || 14 | 0v      |     |      |  21 | GPIO 2 | 13 || 14 | 0v     |     |
//   |  22 | GPIO  3 | 15 || 16 | GPIO  4 | 23  |      |  22 | GPIO 3 | 15 || 16 | GPIO 4 |  23 |
//   |     |    3.3v | 17 || 18 | GPIO  5 | 24  |      |     | 3.3v   | 17 || 18 | GPIO 5 |  24 |
//   |  10 |    MOSI | 19 || 20 | 0v      |     |      |  10 | MOSI   | 19 || 20 | 0v     |     |
//   |   9 |    MISO | 21 || 22 | GPIO  6 | 25  |      |   9 | MISO   | 21 || 22 | GPIO 6 |  25 |
//   |  11 |    SCLK | 23 || 24 | CE0     | 8   |      |  11 | SCLK   | 23 || 24 | CE0    |   8 |
//   |     |      0v | 25 || 26 | CE1     | 7   |      |     | 0v     | 25 || 26 | CE1    |   7 |
//   |   0 |   SDA 0 | 27 || 28 | SCL 0   | 1   |      +-----+--------+----++----+--------+-----+
//   |   5 | GPIO 21 | 29 || 30 | 0v      |     |
//   |   6 | GPIO 22 | 31 || 32 | GPIO 26 | 12  |
//   |  13 | GPIO 23 | 33 || 34 | 0v      |     |
//   |  19 | GPIO 24 | 35 || 36 | GPIO 27 | 16  |
//   |  26 | GPIO 25 | 37 || 38 | GPIO 28 | 20  |
//   |     |      0v | 39 || 40 | GPIO 29 | 21  |
//   +-----+---------+----++----+---------+-----+

// Use mcu pin 10, corresponds to physical pin 19 on the pi
var PIN1 = rpio.Pin(10)

// Use mcu pin 13, corresponds to physical pin 33 on the pi
var PIN2 = rpio.Pin(13)

func main() {
	err := rpio.Open()
	if err != nil {
		fmt.Println("Got error while opening GPIO, error: ", err.Error())
		os.Exit(1)
	}
	defer rpio.Close()

	PIN1.Output()
	PIN2.Output()
	fmt.Println("Setting up output done.")

	c := cron.New()
	c.AddFunc("15 * * * * *", func() { PIN2.Toggle() })
	c.Start()

	for {
		PIN1.Toggle()
		time.Sleep(time.Second / 2)
	}

	os.Exit(0)
}
