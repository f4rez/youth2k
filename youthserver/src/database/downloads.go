package downloads

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type DownloadsStore struct {
	db *pg.DB
}

type DownloadsItem struct {
	gorm.Model
	Title       string
	Description string
	ImgLink     string
}

func NewDownloadsStore(db *pg.DB) *DownloadsStore {
	return &DownloadsStore{
		db: db,
	}
}

func (c *DownloadsStore) createDownloadsItem(dwnl *DownloadsItem) error {
	return c.db.Create(c)
}

func (c *DownloadsStore) getDownloadsItem(dwnl *DownloadsItem) error {
	return c.db.Where("ID = ?", dwnl.ID).First(dwnl).Error
}

func (c *DownloadsStore) updateContdown(dwnl *DownloadsItem) error {
	return c.db.Save(c)
}

func (c *DownloadsStore) getDownloadsItems(limit int) ([]DownloadsItem, error) {
	di := make([]DownloadsItem, limit)
	err := c.db.Limit(limit).Find(&di)
	return di, err
}
