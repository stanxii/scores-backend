import React from 'react';

import { MuiThemeProviderProps } from '@material-ui/core/styles/MuiThemeProvider';
import Document, {
  Head,
  Main,
  NextScript,
  NextDocumentContext,
} from 'next/document';
import flush from 'styled-jsx/server';

interface Props {
  pageContext: MuiThemeProviderProps;
  url: string;
}

class MyDocument extends Document<Props> {
  static getInitialProps = (ctx: NextDocumentContext): any => {
    let pageContext: MuiThemeProviderProps;

    const page = ctx.renderPage((Component: any) => {
      const WrappedComponent = (props: Props) => {
        pageContext = props.pageContext;
        return <Component {...props} />;
      };

      return WrappedComponent;
    });

    return {
      ...page,
      // @ts-ignore
      pageContext,
      // Styles fragment is rendered after the app and page rendering finish.
      styles: (
        <>
          <style
            id="jss-server-side"
            // eslint-disable-next-line react/no-danger
            dangerouslySetInnerHTML={{
              // @ts-ignore
              __html: pageContext.sheetsRegistry.toString(),
            }}
          />
          {flush() || null}
        </>
      ),
    };
  };

  render() {
    const { pageContext } = this.props;

    const theme =
      typeof pageContext.theme === 'function'
        ? pageContext.theme(null)
        : pageContext.theme;

    const themeColor = theme.palette.primary.main;

    const directives = [
      "default-src 'self'",
      "script-src 'self'",
      "img-src 'self' https://*.googleusercontent.com",
      "child-src 'none'",
      "object-src 'none'",
      "font-src 'self' https://fonts.gstatic.com",
      "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
      // 'report-uri /api/csp-violation-report', TODO .. this does not work via <meta> elements
    ];

    if (process.env.NODE_ENV === 'development') {
      // webpack needs 'unsafe-eval'
      directives[1] = "script-src 'unsafe-eval' 'unsafe-inline' 'self'";
      directives.push("connect-src 'self' wss://localhost:*");
    }

    const csp = directives.join(';');

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
          <meta httpEquiv="Content-Security-Policy" content={csp} />
          <meta name="theme-color" content={themeColor} />
          <link
            rel="stylesheet"
            href="https://fonts.googleapis.com/css?family=Roboto:300,400,500"
          />
        </Head>
        <body style={{ overflowY: 'scroll' }}>
          <Main />
          <NextScript />
        </body>
      </html>
    );
  }
}

export default MyDocument;
