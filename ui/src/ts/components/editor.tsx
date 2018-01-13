import * as React from "react";

import { Article } from "../models/article";

interface Props {
  headerHeight: string | number;
  onSubmit: (editingArticle: Article) => void;
  article: Article;
}

interface State {
  title: string;
  body: string;
}

export class EditorComponent extends React.PureComponent<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      body: props.article.body,
      title: props.article.title,
    };
  }

  public render() {
    const onSubmit: React.FormEventHandler<HTMLFormElement> = (e) => {
      e.preventDefault();
      const editingArticle: Article = { title: this.state.title, body: this.state.body};
      this.props.onSubmit(editingArticle);
    };
    return (
      <>
        <form style={{height: "100%", display: "flex", flexDirection: "column"}} onSubmit={onSubmit}>
          { this.renderHeader() }
          { this.renderTextarea() }
        </form>
      </>
    );
  }

  private renderHeader(): React.ReactNode {
    const childStyle: React.CSSProperties = {flexGrow: 0, flexShrink: 0, flexBasis: this.props.headerHeight};
    const parentStyle: React.CSSProperties  = {display: "flex", flexDirection: "row"};
    const headerStyle: React.CSSProperties = { ...childStyle, ...parentStyle };
    return (
      <div className="" style={headerStyle}>
        <div className="input-field" style={{flex: "1 1 auto"}}>
          <input
            className="validate"
            type="text"
            placeholder="Title"
            value={this.state.title}
            onChange={(e) => this.setState({ title: e.target.value })} />
        </div>
        <div style={{minWidth: "10%", marginTop: "14px"}}>
          <button className="btn waves-effect waves-light"><i className="material-icons">publish</i></button>
        </div>
      </div>
    );
  }

  private renderTextarea(): React.ReactNode {
    const textareaStyle: React.CSSProperties = {height: `calc(100% - ${this.props.headerHeight})`};
    return (
      <div className="input-field" style={{flexGrow: 1, flexShrink: 0, flexBasis: "auto", height: 0}}>
        <textarea
          className="materialize-textarea"
          style={textareaStyle}
          placeholder="Body"
          onChange={(e) => this.setState({ body: e.target.value })}
          value={this.state.body}>
        </textarea>
      </div>
    );
  }
}
