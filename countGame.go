package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type mUser struct {
	Username     string
	Firebase_id  string
	Count_clicks int
	Team         string
}

type Team struct {
	Name              string
	Count_clicks      int
	Number_of_members int
}

var db, EEE = sql.Open("postgres", "user=Farez dbname=Farez password=passpass sslmode=disable")

func StoreUserToDB(db *sql.DB, u *mUser) (string, error) {
	res, err := db.Exec(`INSERT INTO mUser (username, firebase_id, count_game_clicks, team) 
							VALUES ($1, $2, $3, $4)`,
		u.Username, u.Firebase_id, u.Count_clicks, u.Team)
	if err != nil {
		log.Println("Store userError: ", err, res, EEE)
		return "", err
	}
	return "title", err
}

func countHandler(w http.ResponseWriter, r *http.Request) {
	t := new(mUser)
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

	muser := mUser{"Farez", "mFirebase_id", 0, "Team Sävsjö"}
	log.Println(StoreUserToDB(db, &muser))

	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}

}
