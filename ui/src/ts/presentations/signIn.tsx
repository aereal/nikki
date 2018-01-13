import * as React from "react";
import { GoogleLogin, GoogleLoginResponse } from "react-google-login";

const isOnlineResponse = (res: any): res is GoogleLoginResponse => {
  return !("code" in res);
};

export const SignInComponent: React.SFC<{}> = () => {
  const clientId = process.env.GOOGLE_OAUTH_CLIENT_ID;
  if (clientId === undefined) {
    throw new Error("GOOGLE_OAUTH_CLIENT_ID required");
  }
  return (
    <div className="row valign-wrapper" style={{minHeight: "100vh"}}>
      <div className="col s12">
        <GoogleLogin
          clientId={clientId}
          buttonText="Sign in with Google"
          onSuccess={(res) => {
            if (isOnlineResponse(res)) {
              const auth = res.getAuthResponse();
              window.localStorage.setItem("id_token", auth.id_token);
              window.location.reload();
            }
          }}
          onFailure={(res) => {
            console.log(res);
            debugger;
          }}
        />
      </div>
    </div>
  );
};
