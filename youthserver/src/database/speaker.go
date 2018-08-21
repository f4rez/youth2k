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

type SpeakerStore struct {
	db *gorm.DB
}

type SpeakerResource struct {
	Store *SpeakerStore
	Ctl   *Control
}

type Speaker struct {
	gorm.Model
	StartDate   time.Time
	EndDate     time.Time
	Title       string
	Description string
	ImgLink     string
}

func NewSpeakerResource(db *gorm.DB, ctl *Control) *SpeakerResource {
	db.AutoMigrate(&Speaker{})
	ss := &SpeakerStore{
		db: db,
	}
	return &SpeakerResource{
		Store: ss,
		Ctl:   ctl,
	}
}

func (c *SpeakerResource) createSpeaker(speaker *Speaker) error {
	c.Ctl.IncrementSpeaker()
	return c.Store.db.Create(&speaker).Error
}

func (c *SpeakerResource) getSpeaker(speaker *Speaker) error {
	return c.Store.db.Where("ID = ?", speaker.ID).First(speaker).Error
}

func (c *SpeakerResource) updateSpeaker(speaker *Speaker) error {
	c.Ctl.IncrementSpeaker()
	return c.Store.db.Save(speaker).Error
}

func (c *SpeakerResource) getSpeakers(limit int) ([]Speaker, error) {
	sesItem := make([]Speaker, limit)
	err := c.Store.db.Limit(limit).Find(&sesItem).Error
	return sesItem, err
}

func (c *SpeakerResource) deleteSpeaker(cc *Speaker) error {
	c.Ctl.IncrementSpeaker()
	return c.Store.db.Delete(cc).Error
}

func (c *SpeakerResource) List(w http.ResponseWriter, r *http.Request) {
	di, err := c.getSpeakers(20)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	c.Ctl.Render.RespondWithJSON(w, http.StatusOK, di)

}

func (c *SpeakerResource) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var count Speaker
	err := decoder.Decode(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.createSpeaker(&count)
	c.Ctl.Render.RespondOK(w)
}

func (c *SpeakerResource) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count := Speaker{}
	id, err := strconv.Atoi(vars["ID"])
	count.ID = uint(id)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = c.getSpeaker(&count)

	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.Ctl.Render.RespondWithJSON(w, http.StatusOK, count)
}

func (c *SpeakerResource) Update(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var count Speaker
	err := decoder.Decode(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.updateSpeaker(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.Ctl.Render.RespondOK(w)

}

func (c *SpeakerResource) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count := Speaker{}
	id, err := strconv.Atoi(vars["ID"])
	count.ID = uint(id)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = c.deleteSpeaker(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.Ctl.Render.RespondOK(w)
}

func (c *SpeakerResource) GetPublicRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", c.List).Methods("GET")
	return r
}

func (c *SpeakerResource) GetRouter(mw mux.MiddlewareFunc) *mux.Router {
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
