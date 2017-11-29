import React from "react";
import { withStyles } from "material-ui/styles";
import AppBar from "material-ui/AppBar";
import Toolbar from "material-ui/Toolbar";
import Typography from "material-ui/Typography";
import Button from "material-ui/Button";
import IconButton from "material-ui/IconButton";
import MenuIcon from "material-ui-icons/Menu";
import Tooltip from "material-ui/Tooltip";

const styles = theme => ({
  root: {
    width: "100%"
  },
  flex: {
    flex: 1
  },
  menuButton: {
    marginLeft: -12,
    marginRight: 20
  }
});

function ButtonAppBar({
  onOpenMenu,
  loginRoute,
  title,
  isLoggedIn,
  user,
  onLogout,
  classes
}) {
  const button = isLoggedIn ? (
    <Tooltip title={user} placement="bottom">
      <Button color="contrast" onClick={onLogout}>
        Logout
      </Button>
    </Tooltip>
  ) : (
    <Button color="contrast" href={loginRoute}>
      Login
    </Button>
  );

  return (
    <div className={classes.root}>
      <AppBar position="fixed">
        <Toolbar>
          <IconButton
            onClick={onOpenMenu}
            className={classes.menuButton}
            color="contrast"
            aria-label="Menu"
          >
            <MenuIcon />
          </IconButton>
          <Typography type="title" color="inherit" className={classes.flex}>
            {title}
          </Typography>
          {button}
        </Toolbar>
      </AppBar>
    </div>
  );
}

export default withStyles(styles)(ButtonAppBar);