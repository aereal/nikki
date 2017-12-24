import GraphiQL = require("graphiql");
import * as React from "react";
import * as ReactDOM from "react-dom";

import { API_ORIGIN } from "./endpoints";
import { AuthedUser } from "./models/user";
import { EditArticlePageComponent, Props as EditArticlePageComponentProps } from "./pages/editArticle";
import { NewArticlePageComponent, Props as NewArticlePageComponentProps } from "./pages/newArticle";

function getInitialProps<T>(): T | null {
  const rawInitialProps = document.body.dataset.initialProps;
  if (rawInitialProps === undefined || rawInitialProps === null) {
    return null;
  }
  const initialProps = JSON.parse(rawInitialProps);
  return initialProps as T;
}

const Router: React.SFC<{ location: Location }> = ({ location }) => {
  switch (location.pathname) {
    case "/":
      const rootProps = getInitialProps<NewArticlePageComponentProps>();
      if (rootProps === null) {
        throw new Error("Invalid initial props");
      }
      return (<NewArticlePageComponent {...rootProps} />);
    case "/graphql":
      const graphqlProps = getInitialProps<{ authedUser: AuthedUser | null }>();
      if (graphqlProps === null) {
        throw new Error("Invalid initial props");
      }

      const fetcher = (params: any): Promise<any> => {
        const headers = new Headers({
          "content-type": "application/json",
        });
        if (graphqlProps.authedUser !== null) {
          headers.append("visitor-key", graphqlProps.authedUser.authKey);
        }

        return window.fetch(`${API_ORIGIN}/graphql`, {
          body: JSON.stringify(params),
          headers,
          method: "post",
        }).then((res) => res.json());
      };
      return (
        <GraphiQL fetcher={fetcher} />
      );
    default:
      if (location.pathname.match(/^\/articles\/\d+/) !== null) {
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
