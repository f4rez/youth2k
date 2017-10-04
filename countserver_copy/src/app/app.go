package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"youth2k/countserver/src/teams"
	"youth2k/countserver/src/users"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	usr := users.MyUser{Firebase_id: vars["id"]}
	log.Println("CreateUser: " + usr.Firebase_id)

	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		log.Println(err)
	}
	log.Println("decoder")

	err := usr.CreateUser(a.DB)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("CreateUser finished")
	respondWithJSON(w, http.StatusOK, usr)
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	usr := users.MyUser{Firebase_id: vars["id"]}
	log.Println("GetUser: " + usr.Firebase_id)

	usr.GetUser(a.DB)
	respondWithJSON(w, http.StatusOK, usr)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	usr := users.MyUser{Firebase_id: vars["id"]}
	log.Println("UpdateUser: " + usr.Firebase_id)

	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		log.Println(err)
	}
	err := usr.UpdateUser(a.DB)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}
	team := new(teams.Team)
	team.Name = usr.Team
	err = team.AddMember(a.DB)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}
	respondWithJSON(w, http.StatusOK, usr)

}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	usr := users.MyUser{Firebase_id: vars["id"]}
	log.Println("DeleteUser: " + usr.Firebase_id)

	usr.DeleteUser(a.DB)
}

func (a *App) castVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	usr := users.MyUser{Firebase_id: vars["id"]}
	could, err := usr.VoteIfPossible(a.DB)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	if could {
		return
	}

	team := teams.Team{Name: vars["team"]}
	team.Vote(a.DB, vars["vote"])

}

func (a *App) clearTeams(w http.ResponseWriter, r *http.Request) {
	teams.ClearTeams(a.DB)
	fmt.Fprint(w, "Deleted all teams")

}

func (a *App) clearUserTable(w http.ResponseWriter, r *http.Request) {
	users.ClearUserTable(a.DB)
	fmt.Fprint(w, "Deleted all users")

}

func (a *App) allHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("All")
	fmt.Fprint(w, "all")
}

func (a *App) pwMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		pw := os.Getenv("APP_PW_FUNCTIONS")
		if pw != vars["pw"] {
			log.Println(pw, vars["pw"])
			respondWithError(w, http.StatusForbidden, "Not Allowed")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *App) InitializeRouters() {
	a.Router.Handle("/user/deleteAll/{pw}", a.pwMiddleware(http.HandlerFunc(a.clearUserTable)))
	a.Router.HandleFunc("/user/{id}", a.createUser).Methods("POST")
	a.Router.HandleFunc("/user/{id}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/user/{id}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/user/{id}", a.deleteUser).Methods("DELETE")

	a.Router.Handle("/team/deleteAll/{pw}", a.pwMiddleware(http.HandlerFunc(a.clearTeams)))
	a.Router.HandleFunc("/team/{id}", a.createUser).Methods("POST")
	a.Router.HandleFunc("/team/{id}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/team/{id}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/team/{id}", a.deleteUser).Methods("DELETE")

	a.Router.HandleFunc("/vote/{id}/{team}/{vote}", a.castVote)

	a.Router.HandleFunc("/", a.allHandler)
}

func (a *App) Initialize(user, password, dbname string) {
	var db, err = sql.Open("postgres", "host=youth2k.c7ygvmc8gmni.eu-central-1.rds.amazonaws.com"+" user="+user+" dbname="+dbname+" password="+password+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	a.DB = db
	a.Router = mux.NewRouter()
	a.InitializeRouters()

}

func (a *App) Run(addr string) {
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}

}
