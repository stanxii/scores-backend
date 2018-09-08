package sqlite

import (
	"database/sql"
	"time"

	"github.com/raphi011/scores"
)

var _ scores.StatisticService = &StatisticService{}

type StatisticService struct {
	DB *sql.DB
}

const (
	statisticFieldsSelectSQL = `
			max(s.name) as name,
			cast((sum(s.won) / cast(count(1) as decimal) * 100) as unsigned) as percentage_won,
			sum(s.points_won) as points_won,
			sum(s.points_lost) as points_lost,
			count(1) as played,
			sum(s.won) as games_won,
			sum(1) - sum(s.won) as games_lost
	`
	ungroupedPlayerStatisticSelectSQL = `
		SELECT 
			s.player_id,
			COALESCE(u.profile_image_url, "") as profile_image,
	` + statisticFieldsSelectSQL + `
		FROM player_statistics s
		JOIN players p ON s.player_id = p.id
		LEFT JOIN users u ON p.user_id = u.id 
		WHERE s.created_at > ?
	`
	groupedPlayerStatisticSQL = `
		GROUP BY s.player_id 
		ORDER BY percentage_won DESC
	`
	playersStatisticSelectSQL = ungroupedPlayerStatisticSelectSQL + groupedPlayerStatisticSQL

	groupPlayersStatisticSelectSQL = ungroupedPlayerStatisticSelectSQL +
		" and s.group_id = ? " +
		groupedPlayerStatisticSQL

	playerStatisticSelectSQL = ungroupedPlayerStatisticSelectSQL +
		" and s.player_id = ? " + groupedPlayerStatisticSQL
)

func parseTimeFilter(filter string) time.Time {
	timeFilter := time.Now()
	year := timeFilter.Year()
	month := timeFilter.Month()
	day := timeFilter.Day()
	loc := timeFilter.Location()

	switch filter {
	case "today":
		timeFilter = time.Date(year, month, day, 0, 0, 0, 0, loc)
	case "month":
		timeFilter = time.Date(year, month-1, day, 0, 0, 0, 0, loc)
	case "thisyear":
		timeFilter = time.Date(year, 1, 1, 0, 0, 0, 0, loc)
	default: // "all"
		timeFilter = time.Unix(0, 0)
	}

	return timeFilter
}

func scanPlayerStatistics(db *sql.DB, query string, args ...interface{}) (scores.PlayerStatistics, error) {
	statistics := scores.PlayerStatistics{}
	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		statistic, err := scanPlayerStatistic(rows)

		if err != nil {
			return nil, err
		}

		statistics = append(statistics, *statistic)
	}

	return statistics, nil
}

func scanPlayerStatistic(scanner scan) (*scores.PlayerStatistic, error) {
	s := &scores.PlayerStatistic{
		Player: &scores.Player{},
	}

	err := scanner.Scan(
		&s.PlayerID,
		&s.Player.ProfileImageURL,
		&s.Player.Name,
		&s.PercentageWon,
		&s.PointsWon,
		&s.PointsLost,
		&s.Played,
		&s.GamesWon,
		&s.GamesLost,
	)

	if err != nil {
		return nil, err
	}

	s.Rank = scores.CalculateRank(int(s.PercentageWon))
	s.Player.ID = s.PlayerID

	return s, nil
}

func (s *StatisticService) Players(filter string) (scores.PlayerStatistics, error) {
	timeFilter := parseTimeFilter(filter)

	statistics, err := scanPlayerStatistics(s.DB, playersStatisticSelectSQL, timeFilter)

	if err != nil {
		return nil, err
	}

	return statistics, nil
}

func (s *StatisticService) PlayersByGroup(groupID uint, filter string) (scores.PlayerStatistics, error) {
	timeFilter := parseTimeFilter(filter)

	statistics, err := scanPlayerStatistics(s.DB, groupPlayersStatisticSelectSQL, timeFilter, groupID)

	if err != nil {
		return nil, err
	}

	return statistics, nil
}

func (s *StatisticService) Player(playerID uint, filter string) (*scores.PlayerStatistic, error) {
	timeFilter := parseTimeFilter(filter)

	row := s.DB.QueryRow(playerStatisticSelectSQL, timeFilter, playerID)

	st, err := scanPlayerStatistic(row)

	if err != nil {
		return nil, err
	}

	return st, nil
}

const (
	playerTeamsStatisticSelectSQL = `
		SELECT 
			MAX(CASE WHEN s.player1_id = ? THEN s.player2_id ELSE s.player1_id END) AS player_id,
			COALESCE(MAX(CASE WHEN s.player1_id = ? THEN u2.profile_image_url ELSE u1.profile_image_url END), "") AS profile_image,
			MAX(CASE WHEN s.player1_id = ? THEN p2.name ELSE p1.name END) AS name,
			CAST((SUM(s.won) / CAST(COUNT(1) AS decimal) * 100) AS unsigned) AS percentage_won,
			SUM(s.points_won) AS points_won,
			SUM(s.points_lost) AS points_lost,
			COUNT(1) AS played,
			SUM(s.won) AS games_won,
			SUM(1) - SUM(s.won) AS games_lost
		FROM team_statistics s
		JOIN players p1 ON s.player1_id = p1.id
		JOIN players p2 ON s.player2_id = p2.id
		LEFT JOIN users u1 ON p1.user_id = u1.id 
		LEFT JOIN users u2 ON p2.user_id = u2.id 
		WHERE (s.player1_id = ? OR s.player2_id = ?) and s.created_at > ?
		GROUP BY s.player1_id, s.player2_id 
		ORDER BY percentage_won DESC
	`
)

func (s *StatisticService) PlayerTeams(playerID uint, filter string) (scores.PlayerStatistics, error) {
	timeFilter := parseTimeFilter(filter)

	statistics, err := scanPlayerStatistics(s.DB, playerTeamsStatisticSelectSQL, playerID, timeFilter)

	if err != nil {
		return nil, err
	}

	return statistics, nil
}
