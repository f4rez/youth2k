drop TABLE team;
drop TABLE mUser;

CREATE TABLE team(	teamname text PRIMARY KEY,
					count_game_clicks int,
					team_members_count int
					);

CREATE TABLE mUser(	username text,
					firebase_id text PRIMARY KEY,
					count_game_clicks int,
					team text,
					last_vote time
					);
