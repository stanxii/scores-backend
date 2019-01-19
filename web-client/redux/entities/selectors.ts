import { createSelector } from 'reselect';

import { denorm } from '../entitySchemas';

import { EntityName } from '../../types';
import { Store } from '../store';

const groupMap = state => state.entities.group.values;
const playerMap = state => state.entities.player.values;
const teamMap = state => state.entities.team.values;
const matchMap = state => state.entities.match.values;
const statisticMap = state => state.entities.statistic.values;
const tournamentMap = state => state.entities.tournament.values;
const userMap = state => state.entities.user.values;
const volleynetplayerMap = state => state.entities.volleynetplayer.values;

export const entityMapSelector = createSelector(
  groupMap,
  playerMap,
  teamMap,
  matchMap,
  statisticMap,
  tournamentMap,
  userMap,
  volleynetplayerMap,
  (
    group,
    player,
    team,
    match,
    statistic,
    tournament,
    user,
    volleynetplayer,
  ) => ({
    group,
    match,
    player,
    statistic,
    team,
    tournament,
    user,
    volleynetplayer,
  }),
);

export const allUsersSelector = (state: Store) =>
  denorm(
    EntityName.User,
    entityMapSelector(state),
    state.entities.user.list.all,
  );

export const allPlayersSelector = (state: Store) =>
  denorm(
    EntityName.Player,
    entityMapSelector(state),
    state.entities.player.list.all,
  );

export const groupSelector = (state: Store, groupId: number) =>
  denorm(EntityName.User, entityMapSelector(state), groupId);

export const groupPlayersSelector = (state: Store, groupId: number) =>
  denorm(
    EntityName.Player,
    entityMapSelector(state),
    state.entities.player.by.group[groupId],
  );

export const matchSelector = (state: Store, id: number) =>
  denorm(EntityName.Match, entityMapSelector(state), id);

const allMatchIdsSelector = state =>
  state.entities.match.all.length ? state.entities.match.all : [];

export const allMatchesSelector = createSelector(
  allMatchIdsSelector,
  entityMapSelector,
  (ids, entities) => denorm(EntityName.Match, entities, ids),
);

export const playerSelector = (state: Store, playerId: number) =>
  denorm(EntityName.Player, entityMapSelector(state), playerId);

export const matchesByGroupSelector = (state: Store, groupId: number) =>
  denorm(
    EntityName.Match,
    entityMapSelector(state),
    state.entities.match.by.group[groupId],
  );

export const matchesByPlayerSelector = (state: Store, playerId: number) =>
  denorm(
    EntityName.Match,
    entityMapSelector(state),
    state.entities.match.by.player[playerId],
  );

export const allStatisticSelector = (state: Store) =>
  denorm(
    EntityName.Statistic,
    entityMapSelector(state),
    state.entities.statistic.list.all,
  );

export const statisticByPlayerTeamSelector = (state: Store, playerId: number) =>
  denorm(
    EntityName.Statistic,
    entityMapSelector(state),
    state.entities.statistic.by.playerTeam[playerId],
  );

export const statisticByGroupSelector = (state: Store, groupId: number) =>
  denorm(
    EntityName.Statistic,
    entityMapSelector(state),
    state.entities.statistic.by.group[groupId],
  );

export const statisticByPlayerSelector = (state: Store, playerId: number) =>
  denorm(
    EntityName.Statistic,
    entityMapSelector(state),
    state.entities.statistic.by.player[playerId][0],
  );

export const tournamentsByLeagueSelector = (
  state: Store,
  leagues: string[],
) => {
  let tournaments = [];
  const entityMap = entityMapSelector(state);

  leagues.forEach(league => {
    tournaments = [
      ...tournaments,
      ...denorm(
        EntityName.TournamentInfo,
        entityMap,
        state.entities.tournament.by.league[league] || [],
      ),
    ];
  });

  return tournaments;
};

export const tournamentSelector = (state: Store, tournamentId: number) =>
  denorm(EntityName.TournamentInfo, entityMapSelector(state), tournamentId);

export const ladderVolleynetplayerSelector = (state: Store, gender: string) =>
  denorm(
    EntityName.VolleynetPlayer,
    entityMapSelector(state),
    state.entities.volleynetplayer.by.ladder[gender],
  );

export const searchVolleynetplayerSelector = (state: Store) =>
  denorm(
    EntityName.VolleynetPlayer,
    entityMapSelector(state),
    state.entities.volleynetplayer.list.search,
  );
