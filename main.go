package main

import (
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"youth2k/src/DMXReciever"

	"github.com/gordonklaus/portaudio"
)

const (
	numberOfPhases  = 1000
	sampleRate      = 192000
	bufferSize      = sampleRate / 16
	a               = 1.0            // height of curve's peak
	b               = bufferSize / 2 // position of the peak
	c               = 0.03           // standart deviation controlling width of the curve
	startHz         = 18023
	phoneSampleRate = 44100
	fftSize         = 2048
	bucketSize      = phoneSampleRate / fftSize
	diffHz          = bucketSize*18*2 + bucketSize
	plotSize        = bufferSize
)

var Arr = make([][]float32, 6)
var gArr = make([]float32, 0)

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 32)
}

func makePlot() {
	sleepTime := 1000 * time.Millisecond

	f, _ := os.Create("valuesL.csv")
	defer f.Close()
	for i, v := range Arr[0] {
		f.WriteString(FloatToString(float64(v)))
		f.WriteString(",")
		f.WriteString(strconv.FormatInt(int64(i), 10))
		f.WriteString("\n")
	} /*
		ff, _ := os.Create("valuesR.csv")
		defer ff.Close()
		for i, v := range Arr[1] {
			ff.WriteString(FloatToString(float64(v)))
			ff.WriteString(",")
			ff.WriteString(strconv.FormatInt(int64(i), 10))
			ff.WriteString("\n")
		}

		f2, _ := os.Create("valuesL2.csv")
		defer f2.Close()
		for i, v := range Arr[2] {
			f2.WriteString(FloatToString(float64(v)))
			f2.WriteString(",")
			f2.WriteString(strconv.FormatInt(int64(i), 10))
			f2.WriteString("\n")
		}
		ff2, _ := os.Create("valuesR2.csv")
		defer ff2.Close()
		for i, v := range Arr[3] {
			ff2.WriteString(FloatToString(float64(v)))
			ff2.WriteString(",")
			ff2.WriteString(strconv.FormatInt(int64(i), 10))
			ff2.WriteString("\n")
		}
		f3, _ := os.Create("valuesL3.csv")
		defer f3.Close()
		for i, v := range Arr[4] {
			f3.WriteString(FloatToString(float64(v)))
			f3.WriteString(",")
			f3.WriteString(strconv.FormatInt(int64(i), 10))
			f3.WriteString("\n")
		}
		ff3, _ := os.Create("valuesR3.csv")
		defer ff3.Close()
		for i, v := range Arr[5] {
			ff3.WriteString(FloatToString(float64(v)))
			ff3.WriteString(",")
			ff3.WriteString(strconv.FormatInt(int64(i), 10))
			ff3.WriteString("\n")
		}*/

	time.Sleep(sleepTime * 50)
}

type stereoSine struct {
	*portaudio.Stream
	stepL, phaseL   float64
	stepR, phaseR   float64
	stepL2, phaseL2 float64
	stepR2, phaseR2 float64
	stepL3, phaseL3 float64
	stepR3, phaseR3 float64
}

func newStereoSine(freqL, freqR, fl2, fl3, fr2, fr3, sampleRate float64) *stereoSine {
	s := &stereoSine{nil, freqL / sampleRate, 0, freqR / sampleRate, 0, fl2 / sampleRate, 0, fr2 / sampleRate, 0, fl3 / sampleRate, 0, fr3 / sampleRate, 0}
	var err error
	s.Stream, err = portaudio.OpenDefaultStream(0, 2, sampleRate, bufferSize, s.processAudio)
	chk(err)
	return s
}

func (g *stereoSine) processAudio(out [][]float32) {

	for i := range out[0] {
		/*
			out[0][i] = float32(math.Sin(2*math.Pi*g.phaseL)) * curveFunc(float64(i))

			_, g.phaseL = math.Modf(g.phaseL + g.stepL)

			out[1][i] = float32(math.Sin(2*math.Pi*g.phaseR)) * curveFunc(float64(i))
			_, g.phaseR = math.Modf(g.phaseR + g.stepR)

			out[2][i] = float32(math.Sin(2*math.Pi*g.phaseL2)) * curveFunc(float64(i)) * 5
			_, g.phaseL2 = math.Modf(g.phaseL2 + g.stepL2)

			out[3][i] = float32(math.Sin(2*math.Pi*g.phaseR2)) * curveFunc(float64(i)) * 5
			_, g.phaseR2 = math.Modf(g.phaseR2 + g.stepR2)

			out[4][i] = float32(math.Sin(2*math.Pi*g.phaseL3)) * curveFunc(float64(i))
			_, g.phaseL3 = math.Modf(g.phaseL3 + g.stepL3)

			out[5][i] = float32(math.Sin(2*math.Pi*g.phaseR3)) * curveFunc(float64(i))
			_, g.phaseR3 = math.Modf(g.phaseR3 + g.stepR3)

			if i < plotSize {
				Arr[0][i] = out[0][i]
				Arr[1][i] = out[1][i]
				Arr[2][i] = out[2][i]
				Arr[3][i] = out[3][i]
				Arr[4][i] = out[4][i]
				Arr[5][i] = out[5][i]
			}
		*/
		if i < len(gArr) {
			out[0][i] = gArr[i] * curveFunc(float64(i))
			out[1][i] = gArr[i] * curveFunc(float64(i))
			log.Println(curveFunc(float64(i)))
			Arr[0][i] = out[0][i]
		}

	}
	/*g.phaseL = 0
	g.phaseL2 = 0
	g.phaseL3 = 0
	g.phaseR = 0
	g.phaseR2 = 0
	g.phaseR3 = 0*/

}

func makeColorHz(start int, color int) (float64, float64) {
	first := color % 16
	second := color / 16
	first = start + int(float32(first)*bucketSize)
	second = start + bucketSize*19 + int(float32(second)*bucketSize)
	return float64(first), float64(second)
}

func genArr(hz int) []float32 {
	arr := make([]float32, bufferSize, bufferSize)
	phase := 0.0
	step := float64(hz) / float64(sampleRate)
	for i := range arr {
		if i < bufferSize-200 {
			arr[i] = float32(math.Sin(2 * math.Pi * phase))
			_, phase = math.Modf(phase + step)
		} else if phase != 0 {
			arr[i] = float32(math.Sin(2 * math.Pi * phase))
			_, phase = math.Modf(phase + step)
		} else {
			arr[i] = 0.0
		}
	}
	return arr

}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func curveFunc(x float64) float32 {
	if x > (bufferSize - 300.0) {
		return float32(1.0 - (300-bufferSize+x)/100)
	}
	return 1 //float32(a * math.Exp(-math.Pow(x-b, 2)/2.0*math.Pow(c, 2)))
}

func main() {
	dmxChan := make(chan DMXReciever.DmxSignal)
	go DMXReciever.ReadArtnet(dmxChan)
	for i := 0; i < 6; i++ {
		Arr[i] = make([]float32, plotSize)
	}
	a := genArr(18000)
	/*for _, i := range a {
		log.Println(i)
	}*/
	gArr = a
	portaudio.Initialize()
	defer portaudio.Terminate()

	sleepTime := 50 * time.Millisecond
	index := 0
	for {
		dmx := <-dmxChan
		first, second := makeColorHz(startHz, dmx.R)
		first2, second2 := makeColorHz(startHz+int(diffHz), dmx.G)
		first3, second3 := makeColorHz(startHz+diffHz*2, dmx.B)

		log.Println(first, second)
		r := newStereoSine(first, second, first2, second2, first3, second3, sampleRate)
		chk(r.Start())
		time.Sleep(sleepTime)
		chk(r.Stop())
		makePlot()
		r.Close()
		index++
		log.Println(index, dmx)

	}
}
