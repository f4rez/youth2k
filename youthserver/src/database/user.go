package users

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"math/rand"
	"time"
)

type UserStore struct {
	db *pg.DB
}

type MyUser struct {
	gorm.Model
	Firebase_id string `gorm:"unique;not null;index:fid"`
	Name        string
}

func NewUserStore(db *pg.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (c *UserStore) GetUser(usr *MyUser) error {
	return c.db.Where("Firebase_id = ?", usr.Firebase_id).First(usr).Error
}

func (c *UserStore) DeleteUser(usr *MyUser) error {
	_, err := db.Exec("DELETE FROM mUser WHERE firebase_id=$1", usr.Firebase_id)

	return nil

}

func (c *UserStore) CreateUser(usr *MyUser) error {
	if usr.Firebase_id != "" {
		log.Println(usr)
		db.Create(usr)
		return nil
	}
	return errors.New("No firebaseID")

}

//Probably not needed
func GetUsers(db *gorm.DB, count int) ([]MyUser, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	max := 0
	db.Model(&MyUser{}).Count(&max)
	mUsers := make([]MyUser, count)
	finished := 0
	for finished < count {
		var usr MyUser
		if err := db.Where("id = ? ", rand.Intn(max)).Find(&usr).Error; err != nil {
			continue
		}
		mUsers[finished] = usr
		finished++
	}
	return mUsers, nil
}

func ClearUserTable(db *gorm.DB) error {
	return nil
}
