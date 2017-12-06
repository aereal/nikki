import * as React from 'react';
import * as ReactDOM from 'react-dom';

interface AuthedUser {
  name: string;
  slug: string;
}

interface AuthenticationComponentProps {
  authenticatedView: React.ReactNode;
  authenticationView: React.ReactNode;
  authenticated: () => boolean;
}
class AuthenticationComponent extends React.PureComponent<AuthenticationComponentProps, {}> {
  render() {
    const { authenticated, authenticationView, authenticatedView } = this.props;
    if (authenticated()) {
      return authenticatedView;
    } else {
      return authenticationView;
    }
  }
}

class SignInComponent extends React.PureComponent<{}, {}> {
  render() {
    return (
      <div className="row valign-wrapper" style={{minHeight: '100vh'}}>
        <div className="col s12">
          <a className="waves-effect waves-light btn-large" href="/auth/google_oauth2">
            <i className="material-icons left">input</i>
            Sign in with Google
          </a>
        </div>
      </div>
    );
  }
}

interface EditorComponentProps {
  headerHeight: string | number;
  onSubmit: React.FormEventHandler<HTMLFormElement>;
}
class EditorComponent extends React.PureComponent<EditorComponentProps, {}> {
  private renderHeader(): React.ReactNode {
    const childStyle: React.CSSProperties = {flexGrow: 0, flexShrink: 0, flexBasis: this.props.headerHeight};
    const parentStyle: React.CSSProperties  = {display: 'flex', flexDirection: 'row'};
    const headerStyle: React.CSSProperties = { ...childStyle, ...parentStyle };
    return (
      <div className="" style={headerStyle}>
        <div className="input-field" style={{flex: '1 1 auto'}}>
          <input className="validate" type="text" placeholder="Title" />
        </div>
        <div style={{minWidth: '10%', marginTop:'14px'}}>
          <button className="btn waves-effect waves-light"><i className="material-icons">publish</i></button>
        </div>
      </div>
    );
  }

  private renderTextarea(): React.ReactNode {
    return (
      <div className="input-field" style={{flexGrow: 1, flexShrink: 0, flexBasis: 'auto', height: 0}}>
        <textarea className="materialize-textarea" style={{height: `calc(100% - ${this.props.headerHeight})`}} placeholder="Body"></textarea>
      </div>
    );
  }

  render() {
    return (
      <>
        <form style={{height: '100%', display: 'flex', flexDirection: 'column'}} onSubmit={this.props.onSubmit}>
          { this.renderHeader() }
          { this.renderTextarea() }
        </form>
      </>
    );
  }
}

interface InitialProps {
  authedUser: AuthedUser | null;
}
class RootComponent extends React.PureComponent<{}, {}> {
  render() {
    const rawInitialProps = document.body.dataset['initialProps'];
    if (rawInitialProps === undefined || rawInitialProps === null) {
      throw new Error("Invalid initial props");
    }
    const initialProps: InitialProps = JSON.parse(rawInitialProps);
    return (
      <AuthenticationComponent
        authenticated={() => initialProps.authedUser !== null }
        authenticatedView={<EditorComponent headerHeight='10vh' onSubmit={(event) => { event.preventDefault(); alert('publish') }} />}
        authenticationView={<SignInComponent />} />
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
