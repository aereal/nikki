/* tslint:disable:max-classes-per-file */

import * as React from "react";
import * as ReactDOM from "react-dom";

import { AuthenticationComponent } from "./components/authentication";
import { API_ORIGIN } from "./endpoints";
import { Article, isPostedArticle , PostedArticle} from "./models/article";
import { AuthedUser } from "./models/user";

function getInitialProps<T>(): T | null {
  const rawInitialProps = document.body.dataset.initialProps;
  if (rawInitialProps === undefined || rawInitialProps === null) {
    return null;
  }
  const initialProps = JSON.parse(rawInitialProps);
  return initialProps as T;
}

const postArticle = (author: AuthedUser, article: Article): Promise<PostedArticle> => {
  const req = window.fetch(`${API_ORIGIN}/articles`, {
    body: JSON.stringify({
      body: article.body,
      title: article.title,
    }),
    credentials: "same-origin",
    headers: {
      "visitor-key": author.authKey,
    },
    method: "POST",
  });
  return req
    .then((res) => res.json())
    .then((json) => {
      if (isPostedArticle(json)) {
        return json;
      } else {
        throw new Error("Invalid response");
      }
    });
};

const updateArticle = (author: AuthedUser, article: PostedArticle): Promise<PostedArticle> => {
  const req = window.fetch(`https://api.nikki.dev/articles/${article.id}`, {
    body: JSON.stringify({
      body: article.body,
      title: article.title,
    }),
    credentials: "same-origin",
    headers: {
      "visitor-key": author.authKey,
    },
    method: "PUT",
  });
  return req
    .then((res) => res.json())
    .then((json) => {
      if (isPostedArticle(json)) {
        return json;
      } else {
        throw new Error("Invalid response");
      }
    });
};

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
  onSubmit: (editingArticle: Article) => void;
  article: Article;
}
interface EditorComponentState {
  title: string;
  body: string;
}
class EditorComponent extends React.PureComponent<EditorComponentProps, EditorComponentState> {
  constructor(props: EditorComponentProps) {
    super(props);
    this.state = {
      body: props.article.body,
      title: props.article.title,
    };
  }

  public render() {
    const onSubmit: React.FormEventHandler<HTMLFormElement> = (e) => {
      e.preventDefault();
      const editingArticle: Article = { title: this.state.title, body: this.state.body};
      this.props.onSubmit(editingArticle);
    };
    return (
      <>
        <form style={{height: "100%", display: "flex", flexDirection: "column"}} onSubmit={onSubmit}>
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

interface EditArticlePageComponentProps {
  authedUser: AuthedUser | null;
  article: PostedArticle;
}
const EditArticlePageComponent: React.SFC<EditArticlePageComponentProps> = ({ authedUser, article }) => {
  const onSubmit = authedUser === undefined || authedUser === null ?
    () => {} :
    (editingArticle: Article) => {
      updateArticle(authedUser, { ...editingArticle, id: article.id }).then((postedArticle) => {
        console.log(postedArticle);
      });
      alert("publish");
    };
  return (
    <AuthenticationComponent
      authenticated={() => authedUser !== null }
      authenticatedView={<EditorComponent headerHeight="10vh" onSubmit={onSubmit} article={article} />}
      authenticationView={<SignInComponent />} />
  );
};

interface RootProps {
  authedUser: AuthedUser | null;
}
class RootComponent extends React.PureComponent<RootProps, {}> {
  public render(): React.ReactNode {
    const authedUser = this.props.authedUser
    const onSubmit = authedUser === undefined || authedUser === null ?
      () => {} :
      (article: Article) => {
        postArticle(authedUser, article).then((postedArticle) => {
          console.log(postedArticle);
        });
        alert("publish");
      };
    const newArticle = { body: "", title: "" };
    return (
      <AuthenticationComponent
        authenticated={() => this.props.authedUser !== null }
        authenticatedView={<EditorComponent headerHeight="10vh" onSubmit={onSubmit} article={newArticle} />}
        authenticationView={<SignInComponent />} />
    );
  }
}

const Router: React.SFC<{ location: Location }> = ({ location }) => {
  switch (location.pathname) {
    case "/":
      const rootProps = getInitialProps<RootProps>();
      if (rootProps === null) {
        throw new Error("Invalid initial props");
      }
      return (<RootComponent {...rootProps} />);
    default:
      if (location.pathname.match(/^\/articles\/\d+/)) {
        const props = getInitialProps<EditArticlePageComponentProps>();
        if (props === null) {
          throw new Error("Invalid initial props");
        }
        return (<EditArticlePageComponent {...props} />);
      } else {
        return null;
      }
  }
};

const entrypoint = document.getElementById("entrypoint");

ReactDOM.render(
  <Router location={window.location} />,
  entrypoint,
);
