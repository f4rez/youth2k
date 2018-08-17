package app

import (
	"database"
)

type API struct {
	Usr  *UserStore
	Cntd *CountdownStore
	Dwnl *DownloadStore
	Shdl *ScheduleStore
	Spek *SpeakerStore
}

func NewApi(db *pg.DB) (*API, error) {
	usr := NewUserStore(db)
	cntd := NewCoundownStore(db)
	dwnl := NewDownloadStore(db)
	shdl := NewScheduleStore(db)
	spek := NewSpeakerSchedule(db)

	api := API{
		Usr:  usr,
		Cntd: cntd,
		Dwnl: dwnl,
		Shdl: shdl,
		Spek: spek,
	}
	return api, nil

}
