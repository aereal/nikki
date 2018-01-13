import * as React from "react";

import { postArticle } from "../actions/articles";
import { AuthenticationComponent } from "../components/authentication";
import { EditorComponent } from "../components/editor";
import { Article } from "../models/article";
import { AuthedUser } from "../models/user";
import { SignInComponent } from "../presentations/signIn";

export interface Props {
  authedUser: AuthedUser | null;
}

export class NewArticlePageComponent extends React.PureComponent<Props, {}> {
  public render(): React.ReactNode {
    const authedUser = this.props.authedUser;
    const onSubmit = authedUser === undefined || authedUser === null ?
      () => {} : // tslint:disable-line:no-empty
      (article: Article) => {
        postArticle(authedUser, article);
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
