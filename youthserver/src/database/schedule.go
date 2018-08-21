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

type ScheduleStore struct {
	db *gorm.DB
}

type ScheduleResource struct {
	Store *ScheduleStore
	Ctl   *Control
}

type ScheduleEntry struct {
	gorm.Model
	Date        time.Time
	StartDate   time.Time
	EndDate     time.Time
	Title       string
	Description string
	ImgLink     string
}

func NewScheduleResource(db *gorm.DB, ctl *Control) *ScheduleResource {
	db.AutoMigrate(&ScheduleEntry{})
	ss := ScheduleStore{
		db: db,
	}

	return &ScheduleResource{
		Store: &ss,
		Ctl:   ctl,
	}
}

func (c *ScheduleResource) createScheduleEntry(se *ScheduleEntry) error {
	c.Ctl.IncrementSchdule()
	return c.Store.db.Create(&se).Error
}

func (c *ScheduleResource) getScheduleEntry(se *ScheduleEntry) error {
	return c.Store.db.Where("ID = ?", se.ID).First(se).Error
}

func (c *ScheduleResource) updateScheduleEntry(se *ScheduleEntry) error {
	c.Ctl.IncrementSchdule()
	return c.Store.db.Save(se).Error
}

func (c *ScheduleResource) getSchduleEntries(limit int) ([]ScheduleEntry, error) {
	sesItem := make([]ScheduleEntry, limit)
	err := c.Store.db.Limit(limit).Find(&sesItem).Error
	return sesItem, err
}

func (c *ScheduleResource) deleteScheduleEntry(cc *ScheduleEntry) error {
	c.Ctl.IncrementSchdule()
	return c.Store.db.Delete(cc).Error
}

func (c *ScheduleResource) List(w http.ResponseWriter, r *http.Request) {
	di, err := c.getSchduleEntries(20)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	c.Ctl.Render.RespondWithJSON(w, http.StatusOK, di)

}

func (c *ScheduleResource) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var count ScheduleEntry
	err := decoder.Decode(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.createScheduleEntry(&count)
	c.Ctl.Render.RespondOK(w)
}

func (c *ScheduleResource) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count := ScheduleEntry{}
	id, err := strconv.Atoi(vars["ID"])
	count.ID = uint(id)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = c.getScheduleEntry(&count)

	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.Ctl.Render.RespondWithJSON(w, http.StatusOK, count)
}

func (c *ScheduleResource) Update(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var count ScheduleEntry
	err := decoder.Decode(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.updateScheduleEntry(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.Ctl.Render.RespondOK(w)

}

func (c *ScheduleResource) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count := ScheduleEntry{}
	id, err := strconv.Atoi(vars["ID"])
	count.ID = uint(id)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = c.deleteScheduleEntry(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.Ctl.Render.RespondOK(w)
}

func (c *ScheduleResource) GetRouter(mw mux.MiddlewareFunc) *mux.Router {
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
