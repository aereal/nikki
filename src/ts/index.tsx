/* tslint:disable:max-classes-per-file */

import * as React from "react";
import * as ReactDOM from "react-dom";

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
  public render() {
    const { authenticated, authenticationView, authenticatedView } = this.props;
    if (authenticated()) {
      return authenticatedView;
    } else {
      return authenticationView;
    }
  }
}

class SignInComponent extends React.PureComponent<{}, {}> {
  public render() {
    return (
      <div className="row valign-wrapper" style={{minHeight: "100vh"}}>
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
interface EditorComponentState {
  title: string;
  body: string;
}
class EditorComponent extends React.PureComponent<EditorComponentProps, EditorComponentState> {
  constructor(props: EditorComponentProps) {
    super(props);
    this.state = {
      body: "",
      title: "",
    };
  }

  public render() {
    return (
      <>
        <form style={{height: "100%", display: "flex", flexDirection: "column"}} onSubmit={this.props.onSubmit}>
          { this.renderHeader() }
          { this.renderTextarea() }
        </form>
      </>
    );
  }

  private renderHeader(): React.ReactNode {
    const childStyle: React.CSSProperties = {flexGrow: 0, flexShrink: 0, flexBasis: this.props.headerHeight};
    const parentStyle: React.CSSProperties  = {display: "flex", flexDirection: "row"};
    const headerStyle: React.CSSProperties = { ...childStyle, ...parentStyle };
    return (
      <div className="" style={headerStyle}>
        <div className="input-field" style={{flex: "1 1 auto"}}>
          <input
            className="validate"
            type="text"
            placeholder="Title"
            value={this.state.title}
            onChange={(e) => this.setState({ title: e.target.value })} />
        </div>
        <div style={{minWidth: "10%", marginTop: "14px"}}>
          <button className="btn waves-effect waves-light"><i className="material-icons">publish</i></button>
        </div>
      </div>
    );
  }

  private renderTextarea(): React.ReactNode {
    const textareaStyle: React.CSSProperties = {height: `calc(100% - ${this.props.headerHeight})`};
    return (
      <div className="input-field" style={{flexGrow: 1, flexShrink: 0, flexBasis: "auto", height: 0}}>
        <textarea
          className="materialize-textarea"
          style={textareaStyle}
          placeholder="Body"
          onChange={(e) => this.setState({ body: e.target.value })}
          value={this.state.body}>
        </textarea>
      </div>
    );
  }
}

interface InitialProps {
  authedUser: AuthedUser | null;
}
class RootComponent extends React.PureComponent<{}, {}> {
  public render() {
    const rawInitialProps = document.body.dataset.initialProps;
    if (rawInitialProps === undefined || rawInitialProps === null) {
      throw new Error("Invalid initial props");
    }
    const initialProps: InitialProps = JSON.parse(rawInitialProps);
    const onSubmit: React.FormEventHandler<HTMLFormElement> = (event) => {
      event.preventDefault();
      alert("publish");
    };
    return (
      <AuthenticationComponent
        authenticated={() => initialProps.authedUser !== null }
        authenticatedView={<EditorComponent headerHeight="10vh" onSubmit={onSubmit} />}
        authenticationView={<SignInComponent />} />
    );
  }
}

const entrypoint = document.getElementById("entrypoint");

ReactDOM.render(
  <>
    <RootComponent />
  </>,
  entrypoint,
);
