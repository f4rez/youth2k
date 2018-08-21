package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"
	//"firebase.google.com/go/auth"
	"google.golang.org/api/option"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	api    *API
	DB     *gorm.DB
	Fir    *firebase.App
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

/*func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var usr users.MyUser
	err := decoder.Decode(&usr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = usr.CreateUser(a.DB)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("CreateUser finished")
	fmt.Fprint(w, http.StatusOK, "OK")
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	usr := users.MyUser{Firebase_id: vars["id"]}
	log.Println("GetUser: " + usr.Firebase_id)

	usr.GetUser(a.DB)
	respondWithJSON(w, http.StatusOK, usr)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	usr := users.MyUser{Firebase_id: vars["id"]}
	log.Println("DeleteUser: " + usr.Firebase_id)

	usr.DeleteUser(a.DB)
}

func (a *App) clearUserTable(w http.ResponseWriter, r *http.Request) {
	users.ClearUserTable(a.DB)
	fmt.Fprint(w, "Deleted all users")

}



func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	usr, err := users.GetUsers(a.DB, 3)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	respondWithJSON(w, http.StatusOK, usr)
}

*/

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
	a.Router.HandleFunc("/", a.allHandler)
	countString := "/countdowns"
	a.Router.PathPrefix(countString).Handler(http.StripPrefix(countString, a.api.Cntd.GetRouter(a.pwMiddleware)))
	downloadsString := "/downloads"
	a.Router.PathPrefix(downloadsString).Handler(http.StripPrefix(downloadsString, a.api.Dwnl.GetRouter(a.pwMiddleware)))
	scheduleString := "/schedule"
	a.Router.PathPrefix(scheduleString).Handler(http.StripPrefix(scheduleString, a.api.Shdl.GetRouter(a.pwMiddleware)))
	userString := "/user"
	a.Router.PathPrefix(userString).Handler(http.StripPrefix(userString, a.api.Usr.GetRouter(a.pwMiddleware)))
	speakerString := "/speakers"
	a.Router.PathPrefix(speakerString).Handler(http.StripPrefix(speakerString, a.api.Spek.GetRouter(a.pwMiddleware)))
	controlString := "/control"
	a.Router.PathPrefix(controlString).Handler(http.StripPrefix(controlString, a.api.Ctlr.GetRouter(a.pwMiddleware)))

}

func (a *App) Initialize(user, password, dbname string) {
	var db, err = gorm.Open("postgres", "host=localhost"+" user="+user+" dbname="+dbname+" password="+password+" sslmode=disable")
	//db.AutoMigrate(&users.MyUser{})
	if err != nil {
		log.Fatal(err)
	}
	opt := option.WithCredentialsFile("youth-conf-firebase-adminsdk-0cwca-7d1d7464f1.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	a.Fir = app
	a.DB = db
	a.Router = mux.NewRouter()

	api, err := NewApi(db)

	a.api = api
	if err != nil {
		log.Fatalf("Error initializing api: %v\n", err)
	}
	a.InitializeRouters()

}

func (a *App) Run(addr string) {
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
