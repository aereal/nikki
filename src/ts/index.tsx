import * as React from 'react';
import * as ReactDOM from 'react-dom';

const entrypoint = document.querySelector('main');

interface AuthedUser {
  name: string;
  slug: string;
}

interface LoginComponentProps {
  authedUser: AuthedUser | null;
}
class LoginComponent extends React.PureComponent<LoginComponentProps, {}> {
  render() {
    const { authedUser } = this.props;
    if (authedUser !== null) {
      return (
        <nav>
          <ul>
            <li>User: {authedUser.name}</li>
            <li><a href="/auth/-/logout">Logout</a></li>
          </ul>
        </nav>
      );
    } else {
      return (
        <nav>
          <ul>
            <li><a href="/auth/google_oauth2">Sign in with Google</a></li>
          </ul>
        </nav>
      );
    }
  }
}

type InitialProps = LoginComponentProps;
class RootComponent extends React.PureComponent<{}, {}> {
  render() {
    const rawInitialProps = document.body.dataset['initialProps'];
    if (rawInitialProps === undefined || rawInitialProps === null) {
      throw new Error("Invalid initial props");
    }
    const initialProps: InitialProps = JSON.parse(rawInitialProps);
    return (
      <>
        <LoginComponent authedUser={initialProps.authedUser} />
      </>
    );
  }
}

ReactDOM.render(
  <>
    <RootComponent />
  </>,
  entrypoint
);
