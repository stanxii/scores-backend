import React from "react";
import { withStyles } from "material-ui/styles";
import fetch from "isomorphic-unfetch";
import PersonIcon from "material-ui-icons/Person";
import Button from "material-ui/Button";
import BackIcon from "material-ui-icons/KeyboardArrowLeft";
import NextIcon from "material-ui-icons/KeyboardArrowRight";
import DoneIcon from "material-ui-icons/Done";
import MobileStepper from "material-ui/MobileStepper";
import Router from "next/router";

import Layout from "../components/Layout";
import CreateMatch from "../components/CreateMatch";
import SetScores from "../components/SetScores";
import withRoot from "../components/withRoot";

const styles = theme => ({
  root: {
    width: "90%"
  },
  button: {
    marginRight: theme.spacing.unit
  },
  stepContainer: {
    padding: "0 20px"
  },
  actionsContainer: {
    marginTop: theme.spacing.unit
  },
  resetContainer: {
    marginTop: 0,
    padding: theme.spacing.unit * 3
  },
  transition: {
    paddingBottom: 4
  },
  button: {
    margin: theme.spacing.unit
  }
});

class NewMatch extends React.Component {
  state = {
    activeStep: 0,
    teamsComplete: false,
    player1ID: 0,
    player2ID: 0,
    player3ID: 0,
    player4ID: 0,
    scoreTeam1: "",
    scoreTeam2: ""
  };

  static async getInitialProps() {
    try {
      const playersResponse = await fetch(
        `${process.env.BACKEND_URL}/api/players`
      );
      const players = await playersResponse.json();

      const playerMap = {};
      const playerIDs = [];
      players.forEach(p => {
        playerIDs.push(p.ID);
        playerMap[p.ID] = p;
      });

      return { playerMap, playerIDs };
    } catch (error) {
      return { error };
    }
  }

  onUnsetPlayer = selected => {
    this.setState({ [`player${selected}ID`]: 0, teamsComplete: false });
  };

  onSetPlayer = (unassigned, ID, teamsComplete) => {
    const activeStep = teamsComplete ? 1 : 0;
    this.setState({ [`player${unassigned}ID`]: ID, teamsComplete, activeStep });
  };

  onSelectTeam = () => {
    this.setState({ activeStep: 0 });
  };

  onSetScores = () => {
    this.setState({ activeStep: 1 });
  };

  getPlayers = () => {
    const { playerIDs = [], playerMap } = this.props;

    return playerIDs.map(p => playerMap[p]);
  };

  onChangeScore = (teamNr, score) => {
    this.setState({ ["scoreTeam" + teamNr]: score });
  };

  onCreateMatch = async () => {
    const { formState, ...match } = this.state;

    try {
      await fetch(`${process.env.BACKEND_URL}/api/matches`, {
        method: "POST",
        body: JSON.stringify(match)
      });

      Router.push("/");
    } catch (e) {
      this.setState({ status: "Error creating match" });
    }
  };

  render() {
    const { playerMap, classes, error } = this.props;
    const {
      teamsComplete,
      activeStep,
      scoreTeam1,
      scoreTeam2,
      ...selectedIDs
    } = this.state;

    const players = this.getPlayers();

    let body;

    const player1 = playerMap[selectedIDs.player1ID];
    const player2 = playerMap[selectedIDs.player2ID];
    const player3 = playerMap[selectedIDs.player3ID];
    const player4 = playerMap[selectedIDs.player4ID];

    return (
      <Layout status={this.state.status} title="New Match">
        <div>
          <MobileStepper
            position="static"
            className={classes.mobileStepper}
            steps={2}
            activeStep={activeStep}
            orientation="vertical"
            backButton={
              <Button
                dense
                disabled={activeStep == 0}
                onClick={this.onSelectTeam}
              >
                <BackIcon className={classes.button} />
                Back
              </Button>
            }
            nextButton={
              <Button
                onClick={this.onSetScores}
                dense
                disabled={activeStep == 1 || !teamsComplete}
              >
                Next
                <NextIcon className={classes.button} />
              </Button>
            }
          />
          <div className={classes.stepContainer}>
            {activeStep == 0 ? (
              <CreateMatch
                {...selectedIDs}
                players={players}
                onSetPlayer={this.onSetPlayer}
                onUnsetPlayer={this.onUnsetPlayer}
              />
            ) : (
              <div>
                <SetScores
                  player1={player1}
                  player2={player2}
                  player3={player3}
                  player4={player4}
                  scoreTeam1={scoreTeam1}
                  scoreTeam2={scoreTeam2}
                  onChangeScore={this.onChangeScore}
                />
                <Button
                  className={classes.button}
                  onClick={this.onCreateMatch}
                  raised
                  color="primary"
                >
                  <DoneIcon className={classes.leftIcon} />
                  Submit
                </Button>
              </div>
            )}
          </div>
        </div>
      </Layout>
    );
  }
}

export default withRoot(withStyles(styles)(NewMatch));
