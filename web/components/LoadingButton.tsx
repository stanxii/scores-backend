import Button from '@material-ui/core/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import { createStyles, Theme, withStyles } from '@material-ui/core/styles';
import React from 'react';

const styles = (theme: Theme) =>
  createStyles({
    buttonProgress: {
      left: '50%',
      marginLeft: -12,
      marginTop: -12,
      position: 'absolute',
      top: '50%',
    },
    wrapper: {
      margin: theme.spacing.unit,
      position: 'relative',
    },
  });

const LoadingButton = ({ children, loading, classes }) => (
  <div className={classes.wrapper}>
    <Button
      color="primary"
      fullWidth
      variant="raised"
      disabled={loading}
      type="submit"
    >
      {children}
    </Button>
    {loading && (
      <CircularProgress size={24} className={classes.buttonProgress} />
    )}
  </div>
);

export default withStyles(styles)(LoadingButton);
