package database

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"encoding/json"
	_ "log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type UserStore struct {
	db *gorm.DB
}

type UserResource struct {
	Store *UserStore
	Ctl   *Control
}

type MyUser struct {
	gorm.Model
	Firebase_id string `gorm:"unique;not null;index:fid"`
	Name        string
}

func NewUserResource(db *gorm.DB, ctl *Control) *UserResource {
	db.AutoMigrate(&MyUser{})
	us := &UserStore{
		db: db,
	}

	return &UserResource{
		Store: us,
		Ctl:   ctl,
	}
}

func (c *UserResource) createMyUser(usr *MyUser) error {
	c.Ctl.IncrementUser()
	return c.Store.db.Create(&usr).Error
}

func (c *UserResource) getMyUser(usr *MyUser) error {
	return c.Store.db.Where("Firebase_id = ?", usr.Firebase_id).First(usr).Error
}

func (c *UserResource) updateMyUser(usr *MyUser) error {
	c.Ctl.IncrementUser()
	return c.Store.db.Save(usr).Error
}

func (c *UserResource) getMyUsers(limit int) ([]MyUser, error) {
	sesItem := make([]MyUser, limit)
	err := c.Store.db.Limit(limit).Find(&sesItem).Error
	return sesItem, err
}

func (c *UserResource) deleteMyUser(cc *MyUser) error {
	c.Ctl.IncrementUser()
	return c.Store.db.Delete(cc).Error
}

func (c *UserResource) getRandomUsers(count int) ([]MyUser, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	max := 0
	c.Store.db.Model(&MyUser{}).Count(&max)
	mUsers := make([]MyUser, count)
	finished := 0
	for finished < count {
		var usr MyUser
		if err := c.Store.db.Where("id = ? ", rand.Intn(max)).Find(&usr).Error; err != nil {
			continue
		}
		mUsers[finished] = usr
		finished++
	}
	return mUsers, nil
}

func (c *UserResource) List(w http.ResponseWriter, r *http.Request) {
	di, err := c.getMyUsers(20)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	c.Ctl.Render.RespondWithJSON(w, http.StatusOK, di)

}

func (c *UserResource) ListRandom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	num, err := strconv.Atoi(vars["num"])
	di, err := c.getRandomUsers(num)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	c.Ctl.Render.RespondWithJSON(w, http.StatusOK, di)

}

func (c *UserResource) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var count MyUser
	err := decoder.Decode(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.createMyUser(&count)
	c.Ctl.Render.RespondOK(w)
}

func (c *UserResource) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count := MyUser{}
	id, err := strconv.Atoi(vars["ID"])
	count.ID = uint(id)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = c.getMyUser(&count)

	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.Ctl.Render.RespondWithJSON(w, http.StatusOK, count)
}

func (c *UserResource) Update(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var count MyUser
	err := decoder.Decode(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.updateMyUser(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.Ctl.Render.RespondOK(w)

}

func (c *UserResource) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count := MyUser{}
	id, err := strconv.Atoi(vars["ID"])
	count.ID = uint(id)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = c.deleteMyUser(&count)
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.Ctl.Render.RespondOK(w)
}

func (c *UserResource) GetRouter(mw mux.MiddlewareFunc) *mux.Router {

	r := mux.NewRouter()
	r.HandleFunc("/", c.List).Methods("GET")
	rc := r.PathPrefix("/admin/").Subrouter()
	rc.Use(mw)
	r.HandleFunc("/", c.Create).Methods("POST")

	rc.HandleFunc("/random/{num}", c.List).Methods("GET")

	r.HandleFunc("/{ID}", c.Get).Methods("GET")
	r.HandleFunc("/{ID}", c.Update).Methods("PUT")
	r.HandleFunc("/{ID}", c.Delete).Methods("DELETE")

	return r

}
