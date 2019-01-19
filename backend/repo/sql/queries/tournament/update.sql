UPDATE tournaments SET
	updated_at = :updatedAt,
	format = :format,
	start_date = :start,
	end_date = :end,
	name = :name,
	league = :league,
	link = :link,
	entry_link = :entrylink,
	status = :status,
	registration_open = :registrationopen,
	location = :location,
	html_notes = :htmlnotes,
	mode = :mode,
	max_points = :maxpoints,
	min_teams = :minteams,
	max_teams = :maxteams,	
	end_registration = :endregistration,
	organiser = :organiser,
	phone = :phone,
	email = :email,
	website = :website,
	current_points = :currentpoints,
	live_scoring_link = :livescoringlink,
	loc_lat = :latitude,
	loc_lon = :longitude,
	season = :season,
	signedup_teams = :signedupteams
WHERE id = :id