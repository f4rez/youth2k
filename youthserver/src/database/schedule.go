package countdown

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type ScheduleStore struct {
	db *pg.DB
}

type ScheduleEntry struct {
	gorm.Model
	Date        time.time
	Title       string
	Description string
	ImgLink     string
}

func NewScheduleStore(db *pg.DB) *ScheduleStore {
	return &ScheduleStore{
		db: db,
	}
}

func (c *ScheduleEntry) createScheduleEntry(se *ScheduleEntry) error {
	return c.db.Create(c)
}

func (c *ScheduleStore) getScheduleEntry(se *ScheduleEntry) error {
	return c.db.Where("ID = ?", c.ID).First(se).Error
}

func (c *ScheduleStore) updateScheduleEntry(se *ScheduleEntry) error {
	return c.db.Save(c)
}

func (c *ScheduleStore) getScheduleEntrys(limit int) ([]ScheduleEntry, error) {
	se := make([]ScheduleEntry, limit)
	err := c.db.Limit(limit).Find(&se)
	return se, err
}
