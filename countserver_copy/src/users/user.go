package users

import (
	"database/sql"
	"errors"
	"time"
)

type MyUser struct {
	Username     string    `json:"username"`
	Firebase_id  string    `json:"fid"`
	Count_clicks int       `json:"clicks"`
	Team         string    `json:"team"`
	LastVote     time.Time `json:"lastVote"`
}

const tableCreationQuery = `CREATE TABLE IF NOT EXSIST mUser(	
					username text,
					firebase_id text PRIMARY KEY,
					count_game_clicks int,
					team text
					);
)`

func (usr *MyUser) GetUser(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM mUser WHERE firebase_id=$1",
		usr.Firebase_id).Scan(&usr.Username, &usr.Firebase_id, &usr.Count_clicks, &usr.Team, &usr.LastVote)
}

func (usr *MyUser) UpdateUser(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE mUser SET username=$1, count_game_clicks = $2, team = $3 WHERE firebase_id=$4",
			usr.Username, usr.Count_clicks, usr.Team, usr.Firebase_id)

	return err
}

func (usr *MyUser) DeleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM mUser WHERE firebase_id=$1", usr.Firebase_id)

	return err
}

func (usr *MyUser) CreateUser(db *sql.DB) error {
	_, err := db.Exec(`INSERT INTO mUser (username, firebase_id, count_game_clicks, team) 
							VALUES ($1, $2, $3, $4)`,
		usr.Username, usr.Firebase_id, usr.Count_clicks, usr.Team)
	return err
}

func (usr *MyUser) updateVote(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE mUser SET count_game_clicks = count_game_clicks +1, last_vote = $1 WHERE firebase_id=$2",
			time.Now, usr.Firebase_id)

	return err
}

func (usr *MyUser) VoteIfPossible(db *sql.DB) (bool, error) {
	canVote := time.Since(usr.LastVote) > 60.0
	if canVote {
		err := usr.updateVote(db)
		return true, err
	}
	return false, nil
}

//Probably not needed
func GetUsers(db *sql.DB, start, count int) ([]MyUser, error) {
	return nil, errors.New("Not implemented")
}

func ClearUserTable(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM mUser")
	return err
}
