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

type DownloadsStore struct {
	db *gorm.DB
}

type DownloadsResource struct {
	Store *DownloadsStore
	Ctl   *Control
}

type DownloadsItem struct {
	gorm.Model
	StartDate   time.Time
	EndDate     time.Time
	Title       string
	Description string
	ImgLink     string
}

func NewDownloadsResource(db *gorm.DB, ctl *Control) *DownloadsResource {
	db.AutoMigrate(&DownloadsItem{})
	ds := &DownloadsStore{
		db: db,
	}

	return &DownloadsResource{
		Store: ds,
		Ctl:   ctl,
	}
}

func (c *DownloadsResource) createDownloadsItem(download *DownloadsItem) error {
	c.Ctl.IncrementDownloads()
	return c.Store.db.Create(&download).Error
}

func (c *DownloadsResource) getDownloadsItem(download *DownloadsItem) error {
	return c.Store.db.Where("ID = ?", download.ID).First(download).Error
}

func (c *DownloadsResource) updateDownloadsItem(download *DownloadsItem) error {
	c.Ctl.IncrementDownloads()
	return c.Store.db.Save(download).Error
}

func (c *DownloadsResource) getDownloadsItems(limit int) ([]DownloadsItem, error) {
	downloadsItem := make([]DownloadsItem, limit)
	err := c.Store.db.Limit(limit).Find(&downloadsItem).Error
	return downloadsItem, err
}

func (c *DownloadsResource) deleteDownloadsItem(cc *DownloadsItem) error {
	c.Ctl.IncrementDownloads()
	return c.Store.db.Delete(cc).Error
}

func (c *DownloadsResource) List(w http.ResponseWriter, r *http.Request) {
	di, err := c.getDownloadsItems(20)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	c.Ctl.Render.RespondWithJSON(w, http.StatusOK, di)

}

func (c *DownloadsResource) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var count DownloadsItem
	err := decoder.Decode(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.createDownloadsItem(&count)
	c.Ctl.Render.RespondOK(w)
}

func (c *DownloadsResource) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count := DownloadsItem{}
	id, err := strconv.Atoi(vars["ID"])
	count.ID = uint(id)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = c.getDownloadsItem(&count)

	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.Ctl.Render.RespondWithJSON(w, http.StatusOK, count)
}

func (c *DownloadsResource) Update(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var count DownloadsItem
	err := decoder.Decode(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.updateDownloadsItem(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.Ctl.Render.RespondOK(w)

}

func (c *DownloadsResource) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count := DownloadsItem{}
	id, err := strconv.Atoi(vars["ID"])
	count.ID = uint(id)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = c.deleteDownloadsItem(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.Ctl.Render.RespondOK(w)
}

func (c *DownloadsResource) GetPublicRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", c.List).Methods("GET")
	return r
}

func (c *DownloadsResource) GetRouter(mw mux.MiddlewareFunc) *mux.Router {
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
