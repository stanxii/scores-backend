CREATE TABLE players (
	id integer PRIMARY KEY,

	created_at datetime NOT NULL,
	updated_at datetime,
	deleted_at datetime,

	first_name varchar(128) CHARSET utf8mb4 NOT NULL,
	last_name varchar(128) CHARSET utf8mb4 NOT NULL,
	total_points integer NOT NULL,
	ladder_rank integer NOT NULL,
	country_union varchar(255) NOT NULL,
	club varchar(255) NOT NULL,
	birthday date,
	license varchar(32) NOT NULL,
	gender varchar(1) NOT NULL,

	INDEX(first_name),
	INDEX(last_name),
	INDEX(birthday),
	INDEX(gender),
	INDEX(ladder_rank)
);

CREATE TABLE users (
	id integer AUTO_INCREMENT PRIMARY KEY,

	created_at datetime NOT NULL,
	updated_at datetime,
	deleted_at datetime,

	email varchar(255) NOT NULL UNIQUE,
	profile_image_url varchar(255) NOT NULL,
	pw_hash blob,
	pw_iterations integer,
	pw_salt blob,
	role varchar(32) NOT NULL,

	player_login varchar(64),
	player_id integer UNIQUE REFERENCES players(id)
);

CREATE TABLE tournaments (
	id integer PRIMARY KEY,

	created_at datetime NOT NULL,
	updated_at datetime,
	deleted_at datetime,

	gender varchar(16) NOT NULL,
	signedup_teams integer NOT NULL,
	start_date datetime NOT NULL,
	end_date datetime NOT NULL,
	name varchar(128) CHARSET utf8mb4 NOT NULL,
	league varchar(128) NOT NULL,
	league_key varchar(128) NOT NULL,
	sub_league varchar(128) NOT NULL,
	sub_league_key varchar(128) NOT NULL,
	link varchar(255) NOT NULL,
	entry_link varchar(255) NOT NULL,
	status varchar(255) NOT NULL,
	registration_open integer NOT NULL,
	location varchar(255) NOT NULL,
	live_scoring_link varchar(255) NOT NULL,
	html_notes text CHARSET utf8mb4 NOT NULL,
	mode varchar(64) NOT NULL,
	max_points integer NOT NULL,
	min_teams integer NOT NULL,
	max_teams integer NOT NULL,
	end_registration datetime,
	organiser varchar(128) NOT NULL,
	phone varchar(128) NOT NULL,
	email varchar(128) NOT NULL,
	website varchar(128) NOT NULL,
	current_points varchar(256) NOT NULL,
	season varchar(16) NOT NULL,
	loc_lat double NOT NULL,
	loc_lon double NOT NULL,
	INDEX(season),
	INDEX(gender),
	INDEX(league_key),
	INDEX(name),
	INDEX(start_date),
	INDEX(end_date)
);


CREATE TABLE tournament_teams (
	tournament_id integer NOT NULL,
	player_1_id integer NOT NULL,
	player_2_id integer NOT NULL,

	created_at datetime NOT NULL,
	updated_at datetime,
	deleted_at datetime,

	result integer NOT NULL,
	seed integer NOT NULL,
	total_points integer NOT NULL,
	won_points integer NOT NULL,
	prize_money real NOT NULL,
	deregistered integer NOT NULL,

	FOREIGN KEY(tournament_id) REFERENCES tournaments(id),
	FOREIGN KEY(player_1_id) REFERENCES players(id),
	FOREIGN KEY(player_2_id) REFERENCES players(id),
	PRIMARY KEY(tournament_id, player_1_id, player_2_id)
);