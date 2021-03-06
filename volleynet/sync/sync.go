package sync

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/raphi011/scores"
	"github.com/raphi011/scores/events"
	"github.com/raphi011/scores/repo"
	"github.com/raphi011/scores/volleynet"
	"github.com/raphi011/scores/volleynet/client"
)

// Changes contains metrics of a scrape job
type Changes struct {
	TournamentInfo TournamentChanges
	Team           TeamChanges
	ScrapeDuration time.Duration
	Success        bool
}

// Service allows loading and synchronizing of the volleynetpage.
type Service struct {
	Log logrus.FieldLogger

	TeamRepo       repo.TeamRepository
	TournamentRepo repo.TournamentRepository
	PlayerRepo     repo.PlayerRepository

	Client        client.Client
	Subscriptions events.Publisher
}

// Tournaments loads tournaments of a certain `gender`, `league` and `season` and
// synchronizes + updates them (if necessary) in the repository.
func (s *Service) Tournaments(gender, league string, season int) error {
	report := &Changes{TournamentInfo: TournamentChanges{}, Team: TeamChanges{}}
	s.publishStartScrapeEvent("tournaments", time.Now())

	current, err := s.Client.Tournaments(gender, league, season)

	if err != nil {
		return errors.Wrap(err, "loading the client tournament list failed")
	}

	persistedTournaments := []*volleynet.Tournament{}
	toDownload := []*volleynet.TournamentInfo{}

	for _, t := range current {
		persisted, err := s.TournamentRepo.Get(t.ID)

		if errors.Cause(err) == scores.ErrNotFound {
			persisted = nil
		} else if err != nil {
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

	currentTournaments := make([]*volleynet.Tournament, len(toDownload))

	for i, t := range toDownload {
		currentTournaments[i], err = s.Client.ComplementTournament(t)

		if err != nil {
			s.Log.Warnf("error loading touappend(slice[:s], slice[s+1:]...)rnament: %v", err)

			// remove it from the tournaments for now
			currentTournaments = append(currentTournaments[:i], currentTournaments[i+1:]...)
		}
	}

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

	err = s.persistTournaments(&report.TournamentInfo)

	if err != nil {
		return err
	}

	return s.persistTeams(&report.Team)
}
