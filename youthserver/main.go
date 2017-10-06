package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"youth2k/youthserver/src/app"
)

func main() {
	a := app.App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	var host = flag.String("host", "0.0.0.0", "IP of host to run webserver on")
	var port = flag.Int("port", 8080, "Port to run webserver on")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Listening on %s", addr)
	a.Run(addr)
}
