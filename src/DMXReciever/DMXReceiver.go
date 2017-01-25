package DMXReciever

import (
	"log"
	"strconv"
	"strings"

	"github.com/tarm/serial"
)

type DmxSignal struct {
	R, G, B int
}

func ReadSerial(dmxChan chan DmxSignal, qChan chan bool) {
	log.Println("TJA")

	c := &serial.Config{Name: "/dev/ttyACM0", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 128)

	rr := -1
	bb := -1
	gg := -1
	for {
		select {
		case <-qChan:
			return
		default:

			n, err := s.Read(buf)
			if err != nil {
				log.Fatal(err)
			}

			dmxOutput := string(buf[:n])
			haveNewLine := strings.Index(dmxOutput, "\r")
			if haveNewLine >= 0 {
				dmxOutput = dmxOutput[:haveNewLine]

			}

			r, err := strconv.Atoi(dmxOutput[:strings.Index(dmxOutput, ",")])
			if err != nil {
				log.Println(err, "1", dmxOutput[:strings.Index(dmxOutput, ",")])
				return
			}
			g, err := strconv.Atoi(dmxOutput[strings.Index(dmxOutput, ",")+1 : strings.LastIndex(dmxOutput, ",")])
			if err != nil {
				log.Println(err, "2", dmxOutput[:strings.Index(dmxOutput, ",")])
				return

			}

			b, err := strconv.Atoi(dmxOutput[strings.LastIndex(dmxOutput, ",")+1 : len(dmxOutput)])
			if err != nil {
				log.Println(err, "3", dmxOutput[:strings.Index(dmxOutput, ",")])
				return

			}
			if rr != r || gg != g || bb != b {
				dmxChan <- DmxSignal{R: r, G: g, B: b}
				rr = r
				gg = g
				bb = b
			}

			break
		}
	}

}
