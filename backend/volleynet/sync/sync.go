package sync

import (
	"time"

	"github.com/pkg/errors"
	"github.com/raphi011/scores/events"
	"github.com/raphi011/scores/repo"
	"github.com/raphi011/scores/volleynet"
	"github.com/raphi011/scores/volleynet/client"
)

// Changes contains metrics of a scrape job
type Changes struct {
	Tournament     *TournamentChanges
	Team           *TeamChanges
	ScrapeDuration time.Duration
	Success        bool
}

// Service allows loading and synchronizing of the volleynetpage.
type Service struct {
	TeamRepo repo.TeamRepository
	TournamentRepo repo.TournamentRepository
	PlayerRepo repo.PlayerRepository
	Client              client.Client
	Subscriptions       events.Publisher
}

// Tournaments loads tournaments of a certain `gender`, `league` and `season` and
// synchronizes + updates them (if necessary) in the repository.
func (s *Service) Tournaments(gender, league string, season int) error {
	report := &Changes{Tournament: &TournamentChanges{}, Team: &TeamChanges{}}
	s.publishStartScrapeEvent("tournaments", time.Now())

	current, err := s.Client.AllTournaments(gender, league, season)

	if err != nil {
		return errors.Wrap(err, "loading the client tournament list failed")
	}

	persistedTournaments := []*volleynet.FullTournament{}
	toDownload := []*volleynet.Tournament{}

	for _, t := range current {
		persisted, err := s.TournamentRepo.Get(t.ID)

		if err != nil {
			return errors.Wrap(err, "loading the persisted tournament failed")
		}

		syncInfo := Tournaments(persisted, t)

		if syncInfo.Type == SyncTournamentNoUpdate {
			continue
		} else if syncInfo.Type != SyncTournamentNew {
			persisted.Teams, err = s.TeamRepo.ByTournament(t.ID)

			if err != nil {
				return errors.Wrap(err, "loading the persisted tournament teams failed")
			}

			persistedTournaments = append(persistedTournaments, persisted)
		}

		toDownload = append(toDownload, t)
	}

	if len(toDownload) == 0 {
		return nil
	}

	currentTournaments, err := s.Client.ComplementMultipleTournaments(toDownload)

	s.syncTournaments(report, persistedTournaments, currentTournaments)

	err = s.persistChanges(report)

	s.publishEndScrapeEvent(report, time.Now())

	return errors.Wrap(err, "sync failed")
}

func (s *Service) persistChanges(report *Changes) error {
	err := s.addMissingPlayers(report.Team.New)

	if err != nil {
		return err
	}

	err = s.persistTournaments(report.Tournament)

	if err != nil {
		return err
	}

	return s.persistTeams(report.Team)
}
