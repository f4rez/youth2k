package app

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"youth2k/youthserver/src/database"

	"log"
)

type API struct {
	Ctl  *database.Control
	Usr  *database.UserResource
	Cntd *database.CountdownResource
	Dwnl *database.DownloadsResource
	Shdl *database.ScheduleResource
	Spek *database.SpeakerResource
	Ctlr *database.ControlResource
}

func NewApi(db *gorm.DB) (*API, error) {
	log.Println("Creating resources for all types")
	ctl := database.NewControl(db)
	usr := database.NewUserResource(db, ctl)
	cntd := database.NewCountdownResource(db, ctl)
	dwnl := database.NewDownloadsResource(db, ctl)
	shdl := database.NewScheduleResource(db, ctl)
	spek := database.NewSpeakerResource(db, ctl)
	ctlr := database.NewControlResource(db, ctl)

	api := API{
		Ctl:  ctl,
		Usr:  usr,
		Cntd: cntd,
		Dwnl: dwnl,
		Shdl: shdl,
		Spek: spek,
		Ctlr: ctlr,
	}
	return &api, nil

}
