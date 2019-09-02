package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/robfig/cron"
	"github.com/stianeikeland/go-rpio"
	//"github.com/stianeikeland/go-rpio"
)

// The library use the raw BCM2835 pinouts, not the ports as they are mapped
// on the output pins for the raspberry pi, and not the wiringPi convention.
//             Rev 2 and 3 Raspberry Pi
//   +-----+---------+----------+---------+-----+
//   | BCM |   Name  | Physical | Name    | BCM |
//   +-----+---------+----++----+---------+-----+
//   |     |    3.3v |  1 || 2  | 5v      |     |
//   |   2 |   SDA 1 |  3 || 4  | 5v      |     |
//   |   3 |   SCL 1 |  5 || 6  | 0v      |     |
//   |   4 | GPIO  7 |  7 || 8  | TxD     | 14  |
//   |     |      0v |  9 || 10 | RxD     | 15  |
//   |  17 | GPIO  0 | 11 || 12 | GPIO  1 | 18  |
//   |  27 | GPIO  2 | 13 || 14 | 0v      |     |
//   |  22 | GPIO  3 | 15 || 16 | GPIO  4 | 23  |
//   |     |    3.3v | 17 || 18 | GPIO  5 | 24  |
//   |  10 |    MOSI | 19 || 20 | 0v      |     |
//   |   9 |    MISO | 21 || 22 | GPIO  6 | 25  |
//   |  11 |    SCLK | 23 || 24 | CE0     | 8   |
//   |     |      0v | 25 || 26 | CE1     | 7   |
//   |   0 |   SDA 0 | 27 || 28 | SCL 0   | 1   |
//   |   5 | GPIO 21 | 29 || 30 | 0v      |     |
//   |   6 | GPIO 22 | 31 || 32 | GPIO 26 | 12  |
//   |  13 | GPIO 23 | 33 || 34 | 0v      |     |
//   |  19 | GPIO 24 | 35 || 36 | GPIO 27 | 16  |
//   |  26 | GPIO 25 | 37 || 38 | GPIO 28 | 20  |
//   |     |      0v | 39 || 40 | GPIO 29 | 21  |
//   +-----+---------+----++----+---------+-----+
//
//
//
//
//
//
//
//
//
//
//
//
//
//

// DelayForSeconds ,delay for several seconds
func DelayForSeconds(pin rpio.Pin, delay int) {
	fmt.Println("Turn on after", delay, "seconds.")

	end := time.Now().Add(time.Second * time.Duration(delay))
	for {
		if time.Now().After(end) {
			fmt.Println("Delay done.")
			break
		}
	}
}

// TurnOnForSeconds ,Turn on for several seconds
func TurnOnForSeconds(pin rpio.Pin, lasts int) {
	fmt.Println("Run cron job.")
	if lasts == 0 {
		return
	}

	fmt.Println("Turn on.")
	pin.High()
	defer pin.Low()

	endTime := time.Now().Add(time.Second * time.Duration(lasts))
	for {
		if time.Now().After(endTime) {
			fmt.Println("Last for", lasts, "seconds done.")
			break
		}
	}

	fmt.Println("Turn off.")
}

func main() {
	rawPin := flag.Uint("pin", 10, "Raw pinouts, not the ports as they are mapped")
	timeDelay := flag.Int("delay", 0, "Delay time in seconds")
	timeLast := flag.Int("lastfor", 0, "Last for time in seconds")
	useCron := flag.Bool("cron", false, "Whether use cron job")
	startHour := flag.Int("hour", 21, "Start hour from 0 to 23")

	flag.Parse()

	err := rpio.Open()
	if err != nil {
		fmt.Println("Got error while opening GPIO, error: ", err.Error())
		os.Exit(1)
	}
	defer rpio.Close()

	fmt.Println("Using pin ", *rawPin)
	pin := rpio.Pin(*rawPin)
	pin.Output()
	fmt.Println("Setting up output done.")

	if *useCron {
		c := cron.New()
		cronStr := fmt.Sprintf("0 0 %d * * *", *startHour)
		c.AddFunc(cronStr, func() { TurnOnForSeconds(pin, *timeLast) })
		c.Start()
		for {
			time.Sleep(time.Second)
		}
	} else {
		DelayForSeconds(pin, *timeDelay)
		TurnOnForSeconds(pin, *timeLast)
	}

	os.Exit(0)
}
