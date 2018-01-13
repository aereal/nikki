import GraphiQL = require("graphiql");
import * as React from "react";
import * as ReactDOM from "react-dom";

import { isSignedIn } from "./authentication";
import { API_ORIGIN } from "./endpoints";
import { EditArticlePageComponent, Props as EditArticlePageComponentProps } from "./pages/editArticle";
import { NewArticlePageComponent, Props as NewArticlePageComponentProps } from "./pages/newArticle";

function getInitialProps<T>(): T | null {
  const rawInitialProps = document.body.dataset.initialProps;
  if (rawInitialProps === undefined || rawInitialProps === null) {
    return null;
  }
  const initialProps = JSON.parse(rawInitialProps);
  initialProps.token = window.localStorage.getItem("id_token");
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
      const graphqlProps = getInitialProps<{ token?: string }>();
      const token = graphqlProps !== null ? graphqlProps.token : undefined;
      const fetcher = (params: any): Promise<any> => {
        const headers = new Headers({
          "content-type": "application/json",
        });
        if (isSignedIn(token)) {
          headers.append("authorization", `bearer ${token}`);
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
