package DMXReciever

import (
	"github.com/tarm/serial"
	"log"
	"net"
	"strconv"
	"strings"
)

type DmxSignal struct {
	R, G, B int
}

func CheckError(err error) {
	if err != nil {
		log.Println("Error: ", err)

	}
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

func ReadArtnet(dmxChan chan DmxSignal) {
	log.Println("UDP readArtnet")
	ServerAddr, err := net.ResolveUDPAddr("udp", ":6454")
	CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	dmxOutput := make([]byte, 1024)

	log.Println("Listning for UDP")
	numberOfDropped := 0
	for {
		_, _, err := ServerConn.ReadFromUDP(dmxOutput)
		r := int(dmxOutput[18])
		g := int(dmxOutput[19])
		b := int(dmxOutput[20])

		if err != nil {
			log.Println("Error: ", err)
		}
		select {
		case dmxChan <- DmxSignal{R: r, G: g, B: b}:
			log.Println("Got to send, number of droppend:", numberOfDropped)
			numberOfDropped = 0
		default:
			numberOfDropped++

		}

	}
}
