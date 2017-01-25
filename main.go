package main

import (
	"github.com/gordonklaus/portaudio"
	//"log"
	"math"
	"time"
)

const sampleRate = 44100

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	sleepTime := 50 * time.Millisecond
	for {
		r := newStereoSine(18000, 18256, sampleRate)

		chk(r.Start())
		time.Sleep(sleepTime)
		chk(r.Stop())
		g := newStereoSine(19000, 19256, sampleRate)
		chk(g.Start())
		time.Sleep(sleepTime)
		chk(g.Stop())
		b := newStereoSine(20000, 20256, sampleRate)

		chk(b.Start())
		time.Sleep(sleepTime)
		chk(b.Stop())

	}
}

type stereoSine struct {
	*portaudio.Stream
	stepL, phaseL float64
	stepR, phaseR float64
}

func newStereoSine(freqL, freqR, sampleRate float64) *stereoSine {
	s := &stereoSine{nil, freqL / sampleRate, 0, freqR / sampleRate, 0}
	var err error
	s.Stream, err = portaudio.OpenDefaultStream(0, 4, sampleRate, 44190/10, s.processAudio)
	chk(err)
	return s
}

func (g *stereoSine) processAudio(out [][]float32) {
	numberOfPhases := 600.0
	rP := 1.0
	lP := 1.0
	p := 0.0
	diff := 0.01
	sL := false
	sR := false
	for i := range out[0] {
		if !sL || lP < numberOfPhases {
			out[0][i] = float32(math.Sin(2*math.Pi*g.phaseL)) * curveFunc(lP)
			//log.Println(out[0][i], lP, curveFunc(lP))
			p, g.phaseL = math.Modf(g.phaseL + g.stepL)
			if p == 1 {
				lP++
			}
			if math.Abs(float64(out[0][i])) < diff {
				sL = true

			}
		} else {
			out[0][i] = 0
		}
		if !sR || rP < numberOfPhases {
			out[1][i] = float32(math.Sin(2*math.Pi*g.phaseR)) * curveFunc(lP)
			p, g.phaseR = math.Modf(g.phaseR + g.stepR)
			if p == 1 {
				rP++
			}
			if math.Abs(float64(out[1][i])) < diff {
				sR = true
			}
		} else {
			out[1][i] = 0
		}
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	a = 1.0   // height of curve's peak
	b = 300   // position of the peak
	c = 0.008 // standart deviation controlling width of the curve
	//( lower abstract value of c -> "longer" curve)
)

func curveFunc(x float64) float32 {
	return float32(a * math.Exp(-math.Pow(x-b, 2)/2.0*math.Pow(c, 2)))
}

/*if g.phaseR < 0.009 && g.phaseL < 0.009 {
	log.Println(g.phaseL, g.phaseR)
	break
}*/
