package tuner

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"math"
)

var frequency float64

type Sinewave struct {
	phase float64
}

func (wave *Sinewave) callback(out []float32) {
	for i := range out {
		out[i] = 0.5 * float32(math.Sin(2*math.Pi*wave.phase))
		wave.phase += frequency / 44100
		if wave.phase >= 1.0 {
			wave.phase -= 1.0
		}
	}
}

func StartTonegen(noteStr string) {

	initSemitones()

	note, err := fromString(noteStr)
	if err != nil {
		panic(err)
	}

	frequency = note.pitch

	fmt.Printf("Initializing tone generator at %s ... ^C to stop", note.toString())

	portaudio.Initialize()
	defer portaudio.Terminate()

	sine := Sinewave{phase: 0}
	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, 512, sine.callback)
	if err != nil {
		panic(err)
	}

	defer stream.Close()

	err = stream.Start()
	if err != nil {
		panic(err)
	}

	select {}
}
