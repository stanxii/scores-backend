// @flow

import React from 'react';
import { withStyles } from 'material-ui/styles';
import List, { ListItem, ListItemText } from 'material-ui/List';
import Chip from 'material-ui/Chip';
import Typography from 'material-ui/Typography';
import { formatDate } from '../utils/dateFormat';
import type { Match, Team } from '../types';

const styles = () => ({
  root: {
    width: '100%',
  },
  listContainer: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    width: '100%',
  },
  team: { flex: '1 1 0' },
  points: { fontWeight: 'lighter', flex: '2 2 0' },
});

type Props = {
  matches: Array<Match>,
  onMatchClick: Match => void,
  classes: Object,
};

function getTeamName(team: Team) {
  if (team.Name) return team.Name;

  return (
    <span>
      {team.player1.name}
      <br />
      {team.player2.name}
    </span>
  );
}

function isSameDay(date1: Date, date2: Date): boolean {
  return (
    date1.getFullYear() === date2.getFullYear() &&
    date1.getMonth() === date2.getMonth() &&
    date1.getDate() === date2.getDate()
  );
}

type DayHeaderProps = {
  date: Date,
};

const DayHeader = ({ date }: DayHeaderProps) => (
  <ListItem dense style={{ justifyContent: 'center' }}>
    <Chip label={formatDate(date)} />
  </ListItem>
);

class MatchList extends React.PureComponent<Props> {
  render() {
    const { matches = [], onMatchClick, classes } = this.props;

    return (
      <List className={classes.root}>
        {matches.map((m, i) => {
          const currentDate = new Date(m.createdAt);
          const lastDate = i ? new Date(matches[i - 1].createdAt) : null;
          const showHeader = !lastDate || !isSameDay(currentDate, lastDate);
          const matchKey = m.id;

          return (
            <React.Fragment key={matchKey}>
              {showHeader ? <DayHeader date={currentDate} /> : null}
              <ListItem divider button onClick={() => onMatchClick(m)}>
                <ListItemText
                  primary={
                    <div className={classes.listContainer}>
                      <Typography className={classes.team} variant="body1">
                        {getTeamName(m.team1)}
                      </Typography>
                      <Typography
                        className={classes.points}
                        variant="display2"
                        align="center"
                      >
                        {m.scoreTeam1} - {m.scoreTeam2}
                      </Typography>
                      <Typography
                        className={classes.team}
                        variant="body1"
                        align="right"
                      >
                        {getTeamName(m.team2)}
                      </Typography>
                    </div>
                  }
                  // secondary={formatDateTime(currentDate)}
                />
              </ListItem>
            </React.Fragment>
          );
        })}
      </List>
    );
  }
}

const StyledMatchList = withStyles(styles)(MatchList);

export default StyledMatchList;
