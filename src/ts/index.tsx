/* tslint:disable:max-classes-per-file */

import * as React from "react";
import * as ReactDOM from "react-dom";

import { postArticle } from "./actions/articles";
import { AuthenticationComponent } from "./components/authentication";
import { EditorComponent } from "./components/editor";
import { Article } from "./models/article";
import { AuthedUser } from "./models/user";
import { EditArticlePageComponent, Props as EditArticlePageComponentProps } from "./pages/editArticle";
import { SignInComponent } from "./presentations/signIn";

function getInitialProps<T>(): T | null {
  const rawInitialProps = document.body.dataset.initialProps;
  if (rawInitialProps === undefined || rawInitialProps === null) {
    return null;
  }
  const initialProps = JSON.parse(rawInitialProps);
  return initialProps as T;
}

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
