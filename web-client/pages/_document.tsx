import React from 'react';

import { MuiThemeProviderProps } from '@material-ui/core/styles/MuiThemeProvider';
import Document, { Head, Main, NextScript } from 'next/document';
import flush from 'styled-jsx/server';

class MyDocument extends Document<{
  pageContext: MuiThemeProviderProps;
}> {
  render() {
    const { pageContext } = this.props;

    const theme =
      typeof pageContext.theme === 'function'
        ? pageContext.theme(null)
        : pageContext.theme;

    const themeColor = theme.palette.primary.main;

    let csp: string;

    if (process.env.NODE_ENV === 'development') {
      // webpack needs 'unsafe-evel'
      csp =
        "default-src 'self'; img-src 'self' https://*.googleusercontent.com; child-src 'none'; script-src 'unsafe-eval' 'unsafe-inline' 'self'; object-src 'none'; font-src 'self' https://fonts.gstatic.com; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com";
    } else {
      // TODO: remove script-src 'unsafe-inline' as soon this does not break next.js anymore
      csp =
        "default-src 'self'; img-src 'self' https://*.googleusercontent.com; child-src 'none'; script-src 'unsafe-inline' 'self'; object-src 'none'; font-src 'self' https://fonts.gstatic.com; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com";
    }

    return (
      <html lang="en" dir="ltr">
        <Head>
          <meta charSet="utf-8" />
          <meta
            name="viewport"
            content={
              'user-scalable=0, initial-scale=1, ' +
              'minimum-scale=1, width=device-width, height=device-height'
            }
          />
          <meta http-equiv="Content-Security-Policy" content={csp} />
          <meta name="theme-color" content={themeColor} />
          <link
            rel="stylesheet"
            href="https://fonts.googleapis.com/css?family=Roboto:300,400,500"
          />
        </Head>
        <body>
          <Main />
          <NextScript />
        </body>
      </html>
    );
  }
}

MyDocument.getInitialProps = ctx => {
  let pageContext;

  const page = ctx.renderPage(Component => {
    const WrappedComponent = props => {
      pageContext = props.pageContext;
      return <Component {...props} />;
    };

    return WrappedComponent;
  });

  return {
    ...page,
    pageContext,
    // Styles fragment is rendered after the app and page rendering finish.
    styles: (
      <>
        <style
          id="jss-server-side"
          // eslint-disable-next-line react/no-danger
          dangerouslySetInnerHTML={{
            __html: pageContext.sheetsRegistry.toString(),
          }}
        />
        {flush() || null}
      </>
    ),
  };
};

export default MyDocument;
