import * as React from "react";

import { updateArticle } from "../actions/articles";
import { isSignedIn } from "../authentication";
import { AuthenticationComponent } from "../components/authentication";
import { EditorComponent } from "../components/editor";
import { Article, PostedArticle } from "../models/article";
import { SignInComponent } from "../presentations/signIn";

export interface Props {
  token?: string;
  article: PostedArticle;
}

export const EditArticlePageComponent: React.SFC<Props> = ({ token, article }) => {
  const onSubmit = isSignedIn(token) ?
    (editingArticle: Article) => {
      updateArticle(token, { ...editingArticle, id: article.id });
      alert("publish");
    } :
    () => {}; // tslint:disable-line:no-empty
  return (
    <AuthenticationComponent
      authenticated={() => isSignedIn(token)}
      authenticatedView={<EditorComponent headerHeight="10vh" onSubmit={onSubmit} article={article} />}
      authenticationView={<SignInComponent />} />
  );
};
