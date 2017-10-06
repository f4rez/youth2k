package users

import (
	"database/sql"
	"errors"
)

type MyUser struct {
	Firebase_id string `json:"fid"`
}

const tableCreationQuery = `CREATE TABLE IF NOT EXSIST mUser(	
					firebase_id text PRIMARY KEY,
					);
)`

func (usr *MyUser) GetUser(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM mUser WHERE firebase_id=$1",
		usr.Firebase_id).Scan(&usr.Firebase_id)
}

func (usr *MyUser) DeleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM mUser WHERE firebase_id=$1", usr.Firebase_id)

	return err
}

func (usr *MyUser) CreateUser(db *sql.DB) error {
	_, err := db.Exec(`INSERT INTO mUser (firebase_id) 
							VALUES ($1)`,
		usr.Firebase_id)
	return err
}

//Probably not needed
func GetUsers(db *sql.DB, start, count int) ([]MyUser, error) {
	return nil, errors.New("Not implemented")
}

func ClearUserTable(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM mUser")
	return err
}
