package sync

import (
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/raphi011/scores/volleynet"
)

// TournamentChanges lists the tournaments that are `New`, `Delete`'d and `Update`'d
// during a sync job
type TournamentChanges struct {
	New    []volleynet.FullTournament
	Delete []volleynet.FullTournament
	Update []volleynet.FullTournament
}

// TournamentSyncInformation contains sync information for two `Tournament`s
type TournamentSyncInformation struct {
	IsNew         bool
	Type          string
	OldTournament *volleynet.FullTournament
	NewTournament *volleynet.Tournament
}

// represents the various tournament sync states.
const (
	SyncTournamentNoUpdate           = "SyncTournamentNoUpdate"
	SyncTournamentUpcomingToCanceled = "SyncTournamentUpcomingToCanceled"
	SyncTournamentUpcoming           = "SyncTournamentUpcoming"
	SyncTournamentUpcomingToDone     = "SyncTournamentUpcomingToDone"
	SyncTournamentNew                = "SyncTournamentNew"
)

func (s *Service) syncTournaments(changes *Changes, oldTournaments, newTournaments []volleynet.FullTournament) {
	oldTournamentMap := createTournamentMap(oldTournaments)
	newTournamentMap := createTournamentMap(newTournaments)

	for key, newTournament := range newTournamentMap {
		oldTournament, ok := oldTournamentMap[key]
		if !ok {
			changes.Tournament.New = append(changes.Tournament.New, newTournament)
		} else {
			mergedTournament := MergeTournament(oldTournament, newTournament)

			if hasTournamentChanged(oldTournament, mergedTournament) {
				changes.Tournament.Update = append(changes.Tournament.Update, mergedTournament)
			}
		}

		s.syncTournamentTeams(changes.Team, oldTournament.Teams, newTournament.Teams)
	}

	for key, oldTournament := range oldTournamentMap {
		if _, ok := oldTournamentMap[key]; !ok {
			changes.Tournament.Delete = append(changes.Tournament.Delete, oldTournament)
		}
	}
}

func createTournamentMap(tournaments []volleynet.FullTournament) map[int]volleynet.FullTournament {
	tournamentMap := make(map[int]volleynet.FullTournament)

	for i := range tournaments {
		t := tournaments[i]
		tournamentMap[t.ID] = t
	}

	return tournamentMap
}

func tournamentSyncType(persisted *volleynet.FullTournament, current *volleynet.Tournament) string {
	if persisted == nil {
		return SyncTournamentNew
	}
	if persisted.Status != volleynet.StatusUpcoming {
		return SyncTournamentNoUpdate
	}
	if current.Status == volleynet.StatusCanceled {
		return SyncTournamentUpcomingToCanceled
	}
	if current.Status == volleynet.StatusUpcoming {
		return SyncTournamentUpcoming
	}
	if current.Status == volleynet.StatusDone {
		return SyncTournamentUpcomingToDone
	}

	return ""
}

func hasTournamentChanged(old, new volleynet.FullTournament) bool {
	new.UpdatedAt = time.Time{}
	old.UpdatedAt = time.Time{}

	return !cmp.Equal(new, old)
}

// Tournaments figures out if and how a tournament needs to be synchronized
func Tournaments(persisted *volleynet.FullTournament, current *volleynet.Tournament) TournamentSyncInformation {
	syncType := tournamentSyncType(persisted, current)

	return TournamentSyncInformation{
		OldTournament: persisted,
		NewTournament: current,
		IsNew:         persisted == nil,
		Type:          syncType,
	}
}

func (s *Service) persistTournaments(changes *TournamentChanges) error {
	for _, new := range changes.New {
		err := s.TournamentRepository.New(&new)

		if err != nil {
			return errors.Wrap(err, "persisting new tournament failed")
		}
	}

	for _, update := range changes.Update {
		err := s.TournamentRepository.Update(&update)

		if err != nil {
			return errors.Wrap(err, "persisting updated tournament failed")
		}
	}

	return nil
}