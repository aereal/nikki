/* tslint:disable:max-classes-per-file */

import * as React from "react";
import * as ReactDOM from "react-dom";

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
