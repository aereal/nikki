import * as React from "react";

import { updateArticle } from "../actions/articles";
import { AuthenticationComponent } from "../components/authentication";
import { EditorComponent } from "../components/editor";
import { Article, PostedArticle } from "../models/article";
import { AuthedUser } from "../models/user";
import { SignInComponent } from "../presentations/signIn";

export interface Props {
  authedUser: AuthedUser | null;
  article: PostedArticle;
}

export const EditArticlePageComponent: React.SFC<Props> = ({ authedUser, article }) => {
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
