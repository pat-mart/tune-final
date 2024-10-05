package tuner

import (
	"errors"
	"math"
	"strconv"
)

var semitoneRatio = math.Pow(2, 1.0/12.0)

var octaveLowerBounds = getNoteSet(lowC)
var octaveSemitones []Note // 9 * 13 tones

type Note struct {
	pitch  float64 // frequency in Hz
	octave uint8
	ltr    string
}

var lowC = Note{
	pitch:  16.35,
	octave: 0,
	ltr:    "C",
}

func initSemitones() {
	for i := 0; i < 9; i++ {
		octave := getFullOctave(&octaveLowerBounds[i])
		for _, note := range octave {
			octaveSemitones = append(octaveSemitones, *note)
		}
	}
}

func (note Note) toString() string {
	octave := strconv.Itoa(int(note.octave))

	return note.ltr + octave
}

func fromString(str string) (Note, error) {
	for _, v := range octaveSemitones {
		if v.toString() == str {
			return v, nil
		}
	}

	return Note{}, errors.New("invalid string input: expected capitalized semitone from octave 0 to 8")
}
