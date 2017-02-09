package main

import (
	"github.com/gordonklaus/portaudio"
	"math"
	"time"
)

const sampleRate = 45000

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()
	s := newStereoSine(18346, 18733, 18798, 19186, 19573, 19961, sampleRate)
	defer s.Close()
	chk(s.Start())
	time.Sleep(500 * time.Second)
	chk(s.Stop())
}

type stereoSine struct {
	*portaudio.Stream
	stepL, phaseL   float64
	stepR, phaseR   float64
	stepE, phaseE   float64
	stepW, phaseW   float64
	stepE2, phaseE2 float64
	stepW2, phaseW2 float64
}

func newStereoSine(freqL, freqR, freqE, freqW, freqE2, freqW2, sampleRate float64) *stereoSine {
	s := &stereoSine{nil, freqL / sampleRate, 0, freqR / sampleRate, 0, freqE / sampleRate, 0, freqW / sampleRate, 0, freqE2 / sampleRate, 0, freqW2 / sampleRate, 0}
	var err error
	s.Stream, err = portaudio.OpenDefaultStream(0, 6, sampleRate, 0, s.processAudio)
	chk(err)
	return s
}

func (g *stereoSine) processAudio(out [][]float32) {
	for i := range out[0] {
		out[0][i] = float32(math.Sin(2 * math.Pi * g.phaseL))
		_, g.phaseL = math.Modf(g.phaseL + g.stepL)
		out[1][i] = float32(math.Sin(2 * math.Pi * g.phaseR))
		_, g.phaseR = math.Modf(g.phaseR + g.stepR)
		out[2][i] = float32(math.Sin(2*math.Pi*g.phaseE)) * 5
		_, g.phaseE = math.Modf(g.phaseE + g.stepE)
		out[3][i] = float32(math.Sin(2*math.Pi*g.phaseW)) * 5
		_, g.phaseW = math.Modf(g.phaseW + g.stepW)
		out[4][i] = float32(math.Sin(2 * math.Pi * g.phaseE2))
		_, g.phaseE2 = math.Modf(g.phaseE2 + g.stepE2)
		out[5][i] = float32(math.Sin(2 * math.Pi * g.phaseW2))
		_, g.phaseW2 = math.Modf(g.phaseW2 + g.stepW2)
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
