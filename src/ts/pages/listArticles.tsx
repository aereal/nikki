import * as React from "react";

import { fetchArticles } from "../actions/articles";
import { AuthenticationComponent } from "../components/authentication";
import { PostedArticle } from "../models/article";
import { AuthedUser } from "../models/user";
import { SignInComponent } from "../presentations/signIn";

export interface Props {
  authedUser: AuthedUser | null;
}

interface State {
  articles: PostedArticle[];
}

const ListArticlesComponent: React.SFC<{ articles: PostedArticle[] }> = ({ articles }) => {
  return (
    <div className="collection">
      {articles.map((article, idx) => {
        return (<a
          href={`/articles/${article.id}`}
          key={idx}
          className="collection-item">
          #{article.id} {article.title}
        </a>);
      })}
    </div>
  );
};

export class ListArticlesPageComponent extends React.PureComponent<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      articles: [],
    };
  }

  public componentDidMount() {
    const { authedUser } = this.props;
    if (authedUser !== null) {
      fetchArticles(authedUser).then((articles) => {
        this.setState({ articles });
      });
    }
  }

  public render() {
    const { authedUser } = this.props;
    const { articles } = this.state;
    return (
      <AuthenticationComponent
        authenticated={() => authedUser !== null}
        authenticatedView={<ListArticlesComponent articles={articles} />}
        authenticationView={<SignInComponent />} />
    );
  }
}
