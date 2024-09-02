package tuner

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"github.com/mjibson/go-dsp/fft"
	"math"
	"math/cmplx"
	"time"
)

type Note struct {
	pitch  float64 // frequency in Hz
	octave uint8
	ltr    string
}

var semitoneRatio = math.Pow(2, 1.0/12.0)

var lowC = Note{
	pitch:  16.35,
	octave: 0,
	ltr:    "C",
}

func getPitch(samples []float32) float64 {
	data := make([]complex128, len(samples))

	for i, v := range samples {
		data[i] = complex128(complex(v, 0))
	}

	freqs := fft.FFT(data)

	magnitude := make([]float64, len(freqs))
	maxMagI := 0

	for i, freq := range freqs {
		// Distance from origin in complex plane
		magnitude[i] = cmplx.Abs(freq)
	}

	for i, mag := range magnitude {
		if mag > magnitude[maxMagI] {
			maxMagI = i
		}
	}

	return (float64(maxMagI) * 44100) / float64(len(samples))
}

func sufficientVolume(samples []float32, cutoff float64) bool {
	return RMS(samples) > cutoff
}

func RMS(samples []float32) float64 {
	var sum float64
	for _, v := range samples {
		sum += math.Pow(float64(v), 2)
	}
	mean := sum / float64(len(samples))
	return math.Sqrt(mean)
}

// Returns a tone from octave 0 to octave 8, e.g. a0, a1 ... a8
func getNoteSet(starter Note) []Note {
	noteSet := make([]Note, 9)
	for i := 0; i < 9; i++ {
		newNote := Note{
			pitch:  starter.pitch * math.Pow(2, float64(i)),
			octave: uint8(i),
			ltr:    starter.ltr,
		}

		noteSet[i] = newNote
	}
	return noteSet
}

// Returns all semitones in an octave
func getFullOctave(starter *Note) []*Note {

	// Between semitone pitches
	ltrs := []string{"C", "C#", "D", "D#", "E", "E#", "F", "F#", "G", "G#", "A", "A#", "B"}

	noteSet := make([]*Note, 13)

	for i := 0; i < 13; i++ {
		newNote := Note{
			ltr:    ltrs[i],
			pitch:  starter.pitch * math.Pow(semitoneRatio, float64(i)),
			octave: starter.octave,
		}

		noteSet[i] = &newNote
	}

	return noteSet
}

var octaveLowerBounds = getNoteSet(lowC)
var octaveSemitones []Note // 9 * 13 tones

func matchToNote(pitchIn float64) (string, int, int) {
	for i := 0; i < 9; i++ {
		octave := getFullOctave(&octaveLowerBounds[i])
		for _, note := range octave {
			octaveSemitones = append(octaveSemitones, *note)
		}
	}

	minI := 0
	maxI := len(octaveSemitones) - 1

	for minI < maxI-1 { // log n! not sure if this contributes to A#8 noise issue
		mid := (minI + maxI) / 2
		if pitchIn < octaveSemitones[mid].pitch {
			maxI = mid
		} else if pitchIn > octaveSemitones[mid].pitch {
			minI = mid
		}
	}

	var note string
	var octave int
	var cents int

	if math.Abs(pitchIn-octaveSemitones[minI].pitch) > math.Abs(pitchIn-octaveSemitones[maxI].pitch) {
		note = octaveSemitones[maxI].ltr
		octave = int(octaveSemitones[maxI].octave)
		cents = int(1200 * math.Log2(pitchIn/octaveSemitones[maxI].pitch))
	} else {
		note = octaveSemitones[minI].ltr
		octave = int(octaveSemitones[minI].octave)
		cents = int(1200 * math.Log2(pitchIn/octaveSemitones[maxI].pitch))
	}

	return note, octave, cents
}

func Main() {

	println("Initializing tuner... ^C to stop")

	err := portaudio.Initialize()
	defer func() {
		err := portaudio.Terminate()
		if err != nil {
			panic(err)
		}
	}()

	soundIn := make([]float32, 1024)

	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(soundIn), soundIn)
	if err != nil {
		panic(err)
	}

	defer func(stream *portaudio.Stream) {
		err := stream.Close()
		if err != nil {
			panic(err)
		}
	}(stream)

	err = stream.Start()
	defer func(stream *portaudio.Stream) {
		err := stream.Stop()
		if err != nil {
			panic(err)
		}
	}(stream)

	println("Play a note!")

	for {
		err = stream.Read()
		if err != nil {
			panic(err)
		}

		pitch := getPitch(soundIn)

		if pitch > 0 && sufficientVolume(soundIn, 0.055) {
			note, octave, cents := matchToNote(pitch)
			fmt.Print("\r")
			fmt.Printf("\r%s %d %d", note, octave, cents)
			time.Sleep(15 * time.Millisecond)
		}
	}
}
