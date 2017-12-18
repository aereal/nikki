import * as React from "react";

export const SignInComponent: React.SFC<{}> = () => {
  return (
    <div className="row valign-wrapper" style={{minHeight: "100vh"}}>
      <div className="col s12">
        <a className="waves-effect waves-light btn-large" href="/auth/google_oauth2">
          <i className="material-icons left">input</i>
          Sign in with Google
        </a>
      </div>
    </div>
  );
};
