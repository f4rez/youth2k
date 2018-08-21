package database

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"youth2k/youthserver/src/render"

	"log"
	"net/http"
)

type Control struct {
	DB     *gorm.DB
	Render *render.Render
}

type ControlItem struct {
	ID                                            int `gorm:"primary_key"`
	Countdown, Downloads, Schedule, Speaker, User int
}

type ControlResource struct {
	db  *gorm.DB
	Ctl *Control
}

func NewControl(db *gorm.DB) *Control {
	db.AutoMigrate(&ControlItem{})
	c := &ControlItem{
		ID:        0,
		Countdown: 0,
		Downloads: 0,
		Schedule:  0,
		Speaker:   0,
		User:      0,
	}
	if db.NewRecord(c) {
		log.Println("First boot on database, creating control")
		db.Create(c)
	} else {
		log.Println("Already have a control object, won't create a new one")
	}
	return &Control{
		DB: db,
	}
}

func NewControlResource(db *gorm.DB, ctl *Control) *ControlResource {
	return &ControlResource{
		db:  db,
		Ctl: ctl,
	}
}

func (c *Control) IncrementCountdown() error {
	ctl, err := c.getControl()
	if err != nil {
		return err
	}
	ctl.Countdown++
	return c.DB.Save(ctl).Error
}
func (c *Control) IncrementDownloads() error {
	ctl, err := c.getControl()
	if err != nil {
		return err
	}
	ctl.Downloads++
	return c.DB.Save(ctl).Error
}
func (c *Control) IncrementSchdule() error {
	ctl, err := c.getControl()
	if err != nil {
		return err
	}
	ctl.Schedule++
	return c.DB.Save(ctl).Error
}
func (c *Control) IncrementSpeaker() error {
	ctl, err := c.getControl()
	if err != nil {
		return err
	}
	ctl.Speaker++
	return c.DB.Save(ctl).Error
}
func (c *Control) IncrementUser() error {
	ctl, err := c.getControl()
	if err != nil {
		return err
	}
	ctl.User++
	return c.DB.Save(ctl).Error
}

func (c *Control) getControl() (*ControlItem, error) {
	ctl := &ControlItem{}
	err := c.DB.First(ctl).Error
	return ctl, err
}

func (c *ControlResource) Get(w http.ResponseWriter, r *http.Request) {
	ctl := &ControlItem{}
	err := c.db.First(ctl).Error
	if err != nil {
		c.Ctl.Render.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.Ctl.Render.RespondWithJSON(w, http.StatusOK, ctl)
}

func (c *ControlResource) GetRouter(_ mux.MiddlewareFunc) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", c.Get).Methods("GET")

	return r

}
