package main

import (
	"github.com/gordonklaus/portaudio"
	//"log"
	"math"
	"src/DMXReciever"
	"time"
)

/*

		18000	18064	18128	18192	18256	18320	18384	18448	18512	18576	18640	18704	18768	18832	18896	18960
18000		0		1		2		3		4		5		6		7		8		9		10		11		12		13		14		15
18064		16		17		18		19		20		21		22		23		24		25		26		27		28		29		30		31
18192
18256
18320
18384
18448
18512
18576
18640
18704
18768
18832
18896
18960


*/

const sampleRate = 44100

func main() {
	dmxChan := make(chan DMXReciever.DmxSignal)
	go DMXReciever.ReadArtnet(dmxChan)

	portaudio.Initialize()
	defer portaudio.Terminate()

	sleepTime := 50 * time.Millisecond
	for {
		dmx := <-dmxChan
		first, second := makeColorHz(18000, dmx.R)
		r := newStereoSine(first, second, sampleRate)
		chk(r.Start())
		time.Sleep(sleepTime)
		chk(r.Stop())
		first, second = makeColorHz(19000, dmx.G)
		g := newStereoSine(first, second, sampleRate)
		chk(g.Start())
		time.Sleep(sleepTime)
		chk(g.Stop())
		first, second = makeColorHz(20000, dmx.B)

		b := newStereoSine(first, second, sampleRate)

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

func makeColorHz(start int, color int) (int, int) {
	first := color % 16
	second := color / 16
	first = start + first*64
	second = start + 1000 + second*64
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
