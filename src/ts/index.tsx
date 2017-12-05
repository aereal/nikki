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

class EditorComponent extends React.PureComponent<{}, {}> {
  render() {
    return (
      <>
        <div id="editor" className="row" style={{height: '75%', marginBottom: 0}}>
          <form className="col s12" style={{height: '100%'}}>
            <div className="col s12" style={{height: '100%', display: 'flex', flexDirection: 'column'}}>
              <div className="input-field col s12" style={{minHeight: '80px', flexGrow: 0, flexShrink: 0, flexBasis: '80px'}}>
                <input className="validate" type="text" placeholder="Title" />
              </div>
              <div className="input-field col s12" style={{flexGrow: 1, flexShrink: 0, flexBasis: '20%'}}>
                <textarea className="materialize-textarea" style={{height: '100%'}} placeholder="Body"></textarea>
              </div>
            </div>
          </form>
        </div>
      </>
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
        { initialProps.authedUser !== null ? <EditorComponent /> : null  }
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
