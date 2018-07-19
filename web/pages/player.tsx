import React from 'react';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import CircularProgress from '@material-ui/core/CircularProgress';
import Router from 'next/router';

import MatchList from '../containers/MatchListContainer';
import withAuth from '../containers/AuthContainer';
import StatisticList from '../components/StatisticList';

import Layout from '../containers/LayoutContainer';
import {
  loadPlayerTeamStatisticAction,
  loadPlayerStatisticAction,
  loadPlayerMatchesAction,
} from '../redux/actions/entities';
import { multiApiAction } from '../redux/actions/api';
import { userSelector } from '../redux/reducers/auth';
import PlayerView from '../components/PlayerView';
import {
  playerSelector,
  statisticByPlayerSelector,
  statisticByPlayerTeamSelector,
  matchesByPlayerSelector,
} from '../redux/reducers/entities';
import { Player, Statistic, Match, User, TeamStatistic } from '../types';

interface Props {
  player: Player;
  user: User;
  statistic: Statistic;
  teamStatistic: TeamStatistic[];
  matches: Match[];
  playerId: number;
  loadMatches: (playerId: number, after?: string) => Promise<any>;
}

interface State {
  tabOpen: number;
  loading: boolean;
  hasMore: boolean;
}

class PlayerInfo extends React.Component<Props, State> {
  static getParameters(query) {
    const { id } = query;

    const playerId = Number.parseInt(id, 10);

    if (Number.isInteger(playerId)) {
      return { playerId };
    }

    return {};
  }

  static shouldComponentUpdate(lastProps, nextProps) {
    return lastProps.playerId !== nextProps.playerId;
  }

  static buildActions({ playerId, user }) {
    const loadPlayerId = playerId || user.playerId;

    return [
      multiApiAction([
        loadPlayerMatchesAction(loadPlayerId),
        loadPlayerStatisticAction(loadPlayerId),
        loadPlayerTeamStatisticAction(loadPlayerId),
      ]),
    ];
  }

  static mapStateToProps(state, ownProps) {
    const { user } = userSelector(state);

    let { playerId } = ownProps;

    playerId = playerId || user.playerId;

    const player = playerSelector(state, playerId);
    const statistic = statisticByPlayerSelector(state, playerId);
    const matches = matchesByPlayerSelector(state, playerId);
    const teamStatistic = statisticByPlayerTeamSelector(state, playerId);

    return {
      player,
      statistic,
      matches,
      teamStatistic,
      user,
    };
  }

  static mapDispatchToProps = {
    loadMatches: loadPlayerMatchesAction,
  };

  state = {
    tabOpen: 0,
    loading: false,
    hasMore: true,
  };

  onLoadMore = async () => {
    const { loadMatches, matches } = this.props;
    const playerId = this.playerId();

    this.setState({ loading: true });

    const lastElement = matches[matches.length - 1];

    const after = lastElement ? lastElement.createdAt : '';

    const newState = {
      loading: false,
      hasMore: true,
    };

    try {
      const result = await loadMatches(playerId, after);
      newState.hasMore = !result.empty;
    } catch (e) {
      newState.hasMore = false;
    } finally {
      this.setState(newState);
    }
  };

  onRowClick = playerId => {
    Router.push(`/player?id=${playerId}`);
  };

  onTabClick = (event, index) => {
    this.setState({ tabOpen: index });
  };

  playerId = () => {
    const { playerId, user } = this.props;

    return playerId || user.playerId;
  };

  render() {
    const { player, matches, statistic, teamStatistic } = this.props;
    const { loading, hasMore } = this.state;

    const playerId = this.playerId();

    const loadingPlayer = !(player && statistic);

    return (
      <Layout title="Players">
        {loadingPlayer ? (
          <CircularProgress />
        ) : (
          <PlayerView player={player} statistic={statistic} />
        )}
        <Tabs
          onChange={this.onTabClick}
          value={this.state.tabOpen}
          textColor="primary"
          fullWidth
        >
          <Tab label={`Matches (${matches.length})`} />
          <Tab label="Teams" />
        </Tabs>
        {this.state.tabOpen === 0 ? (
          <MatchList
            matches={matches}
            highlightPlayerId={playerId}
            onLoadMore={this.onLoadMore}
            loading={loading}
            hasMore={hasMore}
          />
        ) : (
          <StatisticList
            statistics={teamStatistic}
            onPlayerClick={this.onRowClick}
          />
        )}
      </Layout>
    );
  }
}

export default withAuth(PlayerInfo);
