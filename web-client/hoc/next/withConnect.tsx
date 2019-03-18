import React from 'react';

import { NextComponentClass, NextContext } from 'next';
import Error from 'next/error';
import { connect } from 'react-redux';
import { Dispatch, Store as ReduxStore } from 'redux';

import { Store } from '../../redux/store';
import { dispatchActions } from '../../redux/actions';

interface Props {
  reduxStore: ReduxStore;
  fromServer: boolean;
  dispatch: Dispatch;
  initialReduxState: Store;
  error?: { responseCode: number };
}

export interface Context extends NextContext {
  store: ReduxStore;
}

export interface ClientContext extends Context {
  fromServer: boolean;
}

export default (Component: any): NextComponentClass<Props> => {
  class WithConnect extends React.Component<Props> {
    static async getInitialProps(ctx: Context) {
      try {
        const { store, res, req } = ctx;

        const isServer = !!req;

        let props = {
          fromServer: isServer,
          initialReduxState: store.getState(),
          reduxStore: store,
        };

        if (typeof Component.getInitialProps === 'function') {
          const initialProps = await Component.getInitialProps(ctx);

          props = {
            ...props,
            ...initialProps,
          };
        }

        // Execute these only on the server side to avoid waiting for
        // api calls before rendering the page
        if (isServer && Component.buildActions) {
          const actions = Component.buildActions(props);

          await dispatchActions(store.dispatch, actions, req, res);
        }

        return props;
      } catch (e) {
        return { e };
      }
    }

    async componentDidMount() {
      const { fromServer, reduxStore } = this.props;

      if (!Component.buildActions || fromServer) {
        return;
      }

      const actions = Component.buildActions(this.props);
      await dispatchActions(reduxStore.dispatch, actions);
    }

    async componentDidUpdate(nextProps: any, nextState: any) {
      if (
        !Component.buildActions ||
        !Component.shouldComponentFetch ||
        !Component.shouldComponentFetch(nextProps, nextState)
      ) {
        return;
      }

      const { reduxStore } = nextProps;

      const actions = Component.buildActions(nextProps);

      await dispatchActions(reduxStore.dispatch, actions);
    }

    render() {
      const { error, ...props } = this.props;

      if (error) {
        return <Error statusCode={error.responseCode} />;
      }
      return <Component {...props} />;
    }
  }

  return connect(
    Component.mapStateToProps,
    Component.mapDispatchToProps,
  )(WithConnect);
};