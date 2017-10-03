package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type countStruct struct {
	User string
	Team string
}

func countHandler(w http.ResponseWriter, r *http.Request) {
	t := new(countStruct)
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		fmt.Println(err)
	}
	log.Println(t)
	fmt.Fprint(w, t)

}

func main() {
	router := mux.NewRouter()
	var host = flag.String("host", "0.0.0.0", "IP of host to run webserver on")
	var port = flag.Int("port", 8080, "Port to run webserver on")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Listening on %s", addr)

	router.HandleFunc("/count", countHandler)

	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
