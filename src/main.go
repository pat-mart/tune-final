package main

import (
	"fmt"
	"os"
	"strconv"
	"tune-minal/src/metro"
	"tune-minal/src/tuner"
)

func main() {
	argCount := len(os.Args)

	if argCount < 2 {
		fmt.Println("Usage: tuneminal [-tune] [-metronome [-bpm]]")
	}

	if os.Args[1] == "tune" {
		tuner.Main()
	} else if os.Args[1] == "metronome" {
		bpm, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Usage: [-metronome [-bpm]] (bpm must be integer)")
		}
		println(bpm)
		metro.Start(bpm)
	} else {
		fmt.Println("Usage: tuneminal [-tune] [-metronome [-bpm]]")
	}
}
