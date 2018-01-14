import * as React from "react";

import { postArticle } from "../actions/articles";
import { isSignedIn } from "../authentication";
import { AuthenticationComponent } from "../components/authentication";
import { EditorComponent } from "../components/editor";
import { Article } from "../models/article";
import { SignInComponent } from "../presentations/signIn";

export interface Props {
  token?: string;
}

export class NewArticlePageComponent extends React.PureComponent<Props, {}> {
  public render(): React.ReactNode {
    const { token } = this.props;
    const onSubmit = isSignedIn(token) ?
      (article: Article) => {
        postArticle(token, article);
        alert("publish");
      } :
      () => {}; // tslint:disable-line:no-empty
    const newArticle = { body: "", title: "" };
    return (
      <AuthenticationComponent
        authenticated={() => isSignedIn(token)}
        authenticatedView={<EditorComponent headerHeight="10vh" onSubmit={onSubmit} article={newArticle} />}
        authenticationView={<SignInComponent />} />
    );
  }
}
