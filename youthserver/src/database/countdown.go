package database

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"encoding/json"
	_ "log"
	"net/http"
	"strconv"
	"time"
)

type CountdownStore struct {
	db *gorm.DB
}

type CountdownResource struct {
	Store *CountdownStore
	Ctl   *Control
}

type Countdown struct {
	gorm.Model
	Date        time.Time
	StartDate   time.Time
	EndDate     time.Time
	Title       string
	Description string
	ImgLink     string
}

func NewCountdownResource(db *gorm.DB, ctl *Control) *CountdownResource {
	db.AutoMigrate(&Countdown{})
	cs := &CountdownStore{
		db: db,
	}

	return &CountdownResource{
		Store: cs,
		Ctl:   ctl,
	}
}

func (c *CountdownResource) createCountdown(countdown *Countdown) error {
	c.Ctl.IncrementCountdown()
	return c.Store.db.Create(&countdown).Error
}

func (c *CountdownResource) getCountdown(countdown *Countdown) error {
	return c.Store.db.Where("ID = ?", countdown.ID).First(countdown).Error
}

func (c *CountdownResource) updateContdown(countdown *Countdown) error {
	c.Ctl.IncrementCountdown()
	return c.Store.db.Save(countdown).Error
}

func (c *CountdownResource) getCountdowns(limit int) ([]Countdown, error) {
	countdowns := make([]Countdown, limit)
	err := c.Store.db.Limit(limit).Find(&countdowns).Error
	return countdowns, err
}

func (c *CountdownResource) deleteCountdown(cc *Countdown) error {
	c.Ctl.IncrementCountdown()
	return c.Store.db.Delete(cc).Error
}

func (c *CountdownResource) List(w http.ResponseWriter, r *http.Request) {
	cd, err := c.getCountdowns(20)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	c.Ctl.Render.RespondWithJSON(w, http.StatusOK, cd)

}

func (c *CountdownResource) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var count Countdown
	err := decoder.Decode(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.createCountdown(&count)
	c.Ctl.Render.RespondOK(w)
}

func (c *CountdownResource) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count := Countdown{}
	id, err := strconv.Atoi(vars["ID"])
	count.ID = uint(id)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = c.getCountdown(&count)

	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.Ctl.Render.RespondWithJSON(w, http.StatusOK, count)
}

func (c *CountdownResource) Update(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var count Countdown
	err := decoder.Decode(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.updateContdown(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.Ctl.Render.RespondOK(w)

}

func (c *CountdownResource) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count := Countdown{}
	id, err := strconv.Atoi(vars["ID"])
	count.ID = uint(id)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = c.deleteCountdown(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.Ctl.Render.RespondOK(w)
}

func (c *CountdownResource) GetRouter(mw mux.MiddlewareFunc) *mux.Router {
	r := mux.NewRouter()
	rc := r.PathPrefix("/admin/").Subrouter()
	rc.Use(mw)
	r.HandleFunc("/", c.List).Methods("GET")

	rc.HandleFunc("/", c.Create).Methods("POST")

	rc.HandleFunc("/{ID}", c.Get).Methods("GET")
	rc.HandleFunc("/{ID}", c.Update).Methods("PUT")
	rc.HandleFunc("/{ID}", c.Delete).Methods("DELETE")

	return r

}
