import * as React from 'react';
import * as ReactDOM from 'react-dom';

interface AuthedUser {
  name: string;
  slug: string;
}

interface LoginComponentProps {
  authedUser: AuthedUser | null;
}
class LoginComponent extends React.PureComponent<LoginComponentProps, {}> {
  private renderMenu(): React.ReactNode {
    return (
      <li><a href="#"><i className="material-icons">menu</i></a></li>
    );
  }

  private renderSignInLink(): React.ReactNode {
    return (
      <li><a href="/auth/google_oauth2"><i className="material-icons">input</i></a></li>
    );
  }

  render() {
    const { authedUser } = this.props;
    return (
      <nav className="blue-grey">
        <div className="nav-wrapper">
          <ul className="right">
            { authedUser === null ? this.renderSignInLink() : this.renderMenu() }
            <li><a href="#"><i className="material-icons">menu</i></a></li>
          </ul>
        </div>
      </nav>
    );
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

const entrypoint = document.getElementById('entrypoint');

ReactDOM.render(
  <>
    <RootComponent />
  </>,
  entrypoint
);
