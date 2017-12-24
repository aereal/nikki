import * as React from "react";

export = GraphiQL;
declare class GraphiQL extends React.Component<Props, {}> {}

interface GraphiQLVariables {
  [name: string]: any;
}

interface GraphiQLParams {
  query: string;
  variables?: GraphiQLVariables;
  operationName?: string;
}

interface Props {
  fetcher: (params: GraphiQLParams) => Promise<any>;
  // schema: any;
  // query: string | null;
  // variables: string | null;
  // operationName: string | null;
  // response: string | null;
  // storage: any;
  // defaultQuery: string | null;
  // onEditQuery: any;
  // onEditVariables: any;
  // onEditOperationName: any;
  // onToggleDocs: any;
  // getDefaultFieldNames: any;
  // editorTheme: any;
}
