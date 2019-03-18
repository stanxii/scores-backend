import React from 'react';

import { createStyles, WithStyles, withStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import SearchIcon from '@material-ui/icons/Search';
import { connect } from 'react-redux';

import { searchPlayersAction } from '../../redux/entities/actions';
import { searchVolleynetplayerSelector } from '../../redux/entities/selectors';
import SearchPlayerList from './SearchPlayerList';

import { SearchPlayer } from '../../types';
import LoadingButton from '../LoadingButton';
import { Store } from '../../redux/store';

const styles = createStyles({
  container: {
    padding: '0 10px',
  },
});

interface Props extends WithStyles<typeof styles> {
  gender: string;
  foundPlayers: SearchPlayer[];

  onSelectPlayer: (player: SearchPlayer | null) => void;
  searchVolleynetPlayers: (params: {
    fname: string;
    lname: string;
    bday: string;
  }) => void;
}

interface State {
  firstName: string;
  lastName: string;
  birthday: string;
  searching: boolean;
}

class PlayerSearch extends React.Component<Props, State> {
  state = {
    birthday: '',
    firstName: '',
    lastName: '',
    searching: false,
  };

  onChangeFirstname = (event: React.ChangeEvent<HTMLInputElement>) => {
    const firstName = event.target.value;

    this.setState({ firstName });
  };

  onChangeLastname = (event: React.ChangeEvent<HTMLInputElement>) => {
    const lastName = event.target.value;

    this.setState({ lastName });
  };
  onChangeBirthday = (event: React.ChangeEvent<HTMLInputElement>) => {
    const birthday = event.target.value;

    this.setState({ birthday });
  };

  onSearch = async (e: React.FormEvent<HTMLFormElement>) => {
    const { firstName: fname, lastName: lname, birthday: bday } = this.state;
    const { searchVolleynetPlayers } = this.props;

    this.setState({ searching: true });

    e.preventDefault();

    try {
      await searchVolleynetPlayers({ fname, lname, bday });
    } finally {
      this.setState({ searching: false });
    }
  };

  render() {
    const { onSelectPlayer, foundPlayers, classes } = this.props;
    const { firstName, lastName, birthday, searching } = this.state;

    return (
      <form onSubmit={this.onSearch} className={classes.container}>
        <TextField
          label="Firstname"
          type="search"
          margin="normal"
          fullWidth
          onChange={this.onChangeFirstname}
          value={firstName}
        />
        <TextField
          label="Lastname"
          type="search"
          margin="normal"
          fullWidth
          onChange={this.onChangeLastname}
          value={lastName}
        />
        <TextField
          label="Birthday"
          type="search"
          margin="normal"
          fullWidth
          onChange={this.onChangeBirthday}
          value={birthday}
        />
        <SearchPlayerList
          players={foundPlayers}
          onPlayerClick={onSelectPlayer}
        />
        <LoadingButton loading={searching}>
          <SearchIcon />
          <span>Search</span>
        </LoadingButton>
      </form>
    );
  }
}

const mapDispatchToProps = {
  searchVolleynetPlayers: searchPlayersAction,
};

function mapStateToProps(state: Store) {
  const foundPlayers = searchVolleynetplayerSelector(state);

  return { foundPlayers };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(withStyles(styles)(PlayerSearch));
