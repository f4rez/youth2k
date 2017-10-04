package teams

import (
	"database/sql"
	"errors"
)

type Team struct {
	Name              string `json:"teamname"`
	Count_clicks      int    `json:"clicks"`
	Number_of_members int    `json:"member_count"`
}

const tableCreationQuery = `CREATE TABLE IF NOT EXSIST team(	
					teamname text PRIMARY KEY,
					count_game_up_clicks int,
					count_game_down_clicks int,
					team_members_count int
					);
)`

func (team *Team) GetTeam(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM team WHERE team_name=$1",
		team.Name).Scan(&team.Count_clicks, &team.Number_of_members)
}

func (team *Team) UpdateTeam(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE team SET count_game_clicks = $1, team_members_count = $2 WHERE teamname=$3",
			team.Count_clicks, team.Number_of_members, team.Name)
	return err
}

func (team *Team) DeleteTeam(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM team WHERE team_name=$1", team.Name)
	return err
}

func (team *Team) CreateTeam(db *sql.DB) error {
	_, err := db.Exec(`INSERT INTO team (teamname, count_game_clicks, team_members_count) 
							VALUES ($1, $2, $3)`,
		team.Name, team.Count_clicks, team.Number_of_members)
	return err
}

func getTeams(db *sql.DB, start, count int) ([]Team, error) {
	return nil, errors.New("Not implemented")
}

func (team *Team) Vote(db *sql.DB, vote string) error {
	intVote := 0
	if vote == "up" {
		intVote = 1
	} else {
		intVote = -1
	}
	_, err :=
		db.Exec("UPDATE team SET count_game_clicks = count_game_clicks + $1 WHERE teamname=$2",
			intVote, team.Name)
	return err

}

func (team *Team) AddMember(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE team SET team_members_count = team_members_count + $1 WHERE teamname=$2",
			1, team.Name)
	return err
}

func ClearTeams(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM team")
	return err
}
