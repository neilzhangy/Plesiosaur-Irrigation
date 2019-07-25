package main

import (
	"flag"
	"fmt"
	"os"
	"time"

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

func main() {
	rawPin := flag.Uint("pin", 10, "Raw pinouts, not the ports as they are mapped")
	turnOn := flag.Bool("turnon", true, "Whether turn on target GPIO")
	timeDelay := flag.Int("delay", 0, "Delay time in seconds")
	lastFor := flag.Int("lastfor", 0, "Last for time in seconds")

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

	if *turnOn {
		fmt.Println("Turn On after", *timeDelay, "seconds.")
	} else {
		fmt.Println("Turn off after", *timeDelay, "seconds.")
	}

	end := time.Now().Add(time.Second * time.Duration(*timeDelay))
	for {
		if time.Now().After(end) {
			fmt.Println("Delay done.")
			break
		}
	}

	end = time.Now().Add(time.Second * time.Duration(*lastFor))
	if 0 == *lastFor {
		end = time.Now().Add(time.Hour * time.Duration(1000000))
	}
	if *turnOn {
		fmt.Println("Turn on.")
		pin.High()
	} else {
		fmt.Println("Turn off.")
		pin.Low()
	}
	for {
		if time.Now().After(end) {
			fmt.Println("Last for", *lastFor, "seconds done.")
			break
		}
	}
	if *turnOn {
		fmt.Println("Turn off.")
		pin.Low()
	} else {
		fmt.Println("Turn on.")
		pin.High()
	}

	// c := cron.New()
	// c.AddFunc("15 * * * * *", func() { PIN2.Toggle() })
	// c.Start()

	// for {
	// 	PIN1.Toggle()
	// 	time.Sleep(time.Second / 2)
	// }

	os.Exit(0)
}
