package metro

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"math"
	"time"
)

const samples = int(44100 * 0.02)

var tone = make([]float32, samples)

func generateClick(freq float64) []float32 {
	for i := range tone {
		tone[i] = float32(math.Sin(2*math.Pi*freq*float64(i)/44100) * 0.6)
	}

	return tone
}

func playClick(tone []float32) {

	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, len(tone)/4, tone)
	if err != nil {
		panic(err)
	}
	stream.Info().OutputLatency = time.Duration(1000*len(tone)*(len(tone)/4/44100)) * time.Millisecond
	defer stream.Close()

	err = stream.Start()
	err = stream.Write()
	if err != nil {
		return
	}
	defer stream.Stop()
}

func Start(bpm int) {
	portaudio.Initialize()
	defer portaudio.Terminate()

	interval := time.Minute / time.Duration(bpm)

	tone := generateClick(400)
	downBeat := generateClick(440)

	ticker := time.NewTicker(interval - 10)

	beatCount := 0
	numerator := 4

	fmt.Printf("Playing at %d bpm... ^C to stop", bpm)

	for {
		select {
		case <-ticker.C:
			if beatCount == numerator || beatCount == 0 {
				beatCount = 0
				playClick(downBeat)
			} else {
				playClick(tone)
			}
			beatCount++
		}
	}
}
