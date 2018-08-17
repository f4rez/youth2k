package speaker

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type SpeakerStore struct {
	db *pg.DB
}

type Speaker struct {
	gorm.Model
	Date        time.time
	Title       string
	Description string
	ImgLink     string
}

func NewSpeakerStore(db *pg.DB) *SpeakerStore {
	return &SpeakerStore{
		db: db,
	}
}

func (c *SpeakerStore) createSpeaker(spk *Speaker) error {
	return c.db.Create(spk)
}

func (c *SpeakerStore) getSpeaker(spk *Speaker) error {
	return c.db.Where("ID = ?", c.ID).First(spk).Error
}

func (c *SpeakerStore) updateSpeaker(spk *Speaker) error {
	return c.db.Save(spk)
}

func (c *SpeakerStore) getSpeakers(limit int) ([]Speaker, error) {
	speakers := make([]Speaker, limit)
	err := c.db.Limit(limit).Find(&Speaker)
	return speakers, err
}
