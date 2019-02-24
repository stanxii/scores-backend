import React from 'react';

import Link from 'next/link';

import { Theme } from '@material-ui/core';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import { createStyles, WithStyles, withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';

import Login from '../../components/volleynet/Login';
import SearchPlayer from '../../components/volleynet/SearchPlayer';
import withAuth from '../../containers/AuthContainer';
import Layout from '../../containers/LayoutContainer';
import { userSelector } from '../../redux/auth/selectors';
import {
  loadTournamentAction,
  tournamentSignupAction,
} from '../../redux/entities/actions';
import { tournamentSelector } from '../../redux/entities/selectors';
import { link } from '../../styles/shared';

import {
  Tournament,
  User,
  SearchPlayer as SearchPlayerType,
} from '../../types';
import { Query } from '../../utils/query';
import { Store } from '../../redux/store';

const styles = (theme: Theme) =>
  createStyles({
    container: {
      padding: theme.spacing.unit * 2,
    },
    link,
  });

interface Props extends WithStyles<typeof styles> {
  tournamentId: string;
  tournament?: Tournament;
  user: User;

  signup: (info: {
    username: string;
    password: string;
    partnerId: number;
    tournamentId: string;
    // partnerName: string;
    rememberMe: boolean;
  }) => void;
}

interface State {
  partner: SearchPlayerType | null;
}

class Signup extends React.Component<Props, State> {
  static mapDispatchToProps = {
    signup: tournamentSignupAction,
  };
  static getParameters(query: Query) {
    const { id } = query;

    return { tournamentId: id };
  }

  static buildActions({ tournamentId }: Props) {
    return [loadTournamentAction(tournamentId)];
  }

  static mapStateToProps(state: Store, { tournamentId }: Props) {
    const tournament = tournamentSelector(state, tournamentId);
    const { user } = userSelector(state);

    return { tournament, user };
  }

  state: State = {
    partner: null,
  };

  onSelectPlayer = (partner: SearchPlayerType | null) => {
    this.setState({ partner });
  };

  onSignup = async (
    username: string,
    password: string,
    rememberMe: boolean,
  ) => {
    const { tournamentId, signup } = this.props;
    const { partner } = this.state;

    const partnerId = partner && partner.id;
    // const partnerName = partner && partner.login;

    if (!partnerId) {
      return;
    }

    const body = {
      partnerId,
      // partnerName,
      password,
      rememberMe,
      tournamentId,
      username,
    };

    signup(body);
  };

  render() {
    const { partner } = this.state;
    const { tournament, user, classes } = this.props;

    if (!tournament) {
      return null;
    }

    let content: React.ReactNode;

    if (partner) {
      content = (
        <>
          <Typography variant="h6">{`Partner: ${partner.firstName} ${
            partner.lastName
          }`}</Typography>
          <Login onLogin={this.onSignup} username={user.volleynetLogin} />
        </>
      );
    } else {
      content = (
        <SearchPlayer
          gender={tournament.gender}
          onSelectPlayer={this.onSelectPlayer}
        />
      );
    }

    return (
      <Layout
        title={{
          href: `/tournament?id=${tournament.id}`,
          text: 'Signup',
        }}
      >
        <div className={classes.container}>
          <div className={classes.link}>
            <Link href={`/tournament?id=${tournament.id}`}>
              <a className={classes.link}>
                <Typography variant="h4">{tournament.name}</Typography>
              </a>
            </Link>
          </div>
          <Card>
            <CardContent>{content}</CardContent>
          </Card>
        </div>
      </Layout>
    );
  }
}

export default withAuth(withStyles(styles)(Signup));
