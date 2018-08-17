package countdown

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type CountdownStore struct {
	db *pg.DB
}

type Countdown struct {
	gorm.Model
	Date        time.time
	Title       string
	Description string
	ImgLink     string
}

func NewContdownStore(db *pg.DB) *CountdownStore {
	return &CountdownStore{
		db: db,
	}
}

func (c *CountdownStore) createCountdown(countdown *Countdown) error {
	return c.db.Create(c)
}

func (c *CountdownStore) getCountdown(countdown *Countdown) error {
	return c.db.Where("ID = ?", countdown.ID).First(countdown).Error
}

func (c *CountdownStore) updateContdown(countdown *Countdown) error {
	return c.db.Save(c)
}

func (c *CountdownStore) getCountdowns(countdown *Countdown, limit int) ([]Countdown, error) {
	countdowns := make([]Countdown, limit)
	err := c.db.Limit(limit).Find(&countdowns)
	return countdowns, err
}
