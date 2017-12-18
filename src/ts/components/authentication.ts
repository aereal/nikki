import { PureComponent, ReactNode } from "react";

interface Props {
  authenticatedView: ReactNode;
  authenticationView: ReactNode;
  authenticated: () => boolean;
}

export class AuthenticationComponent extends PureComponent<Props, {}> {
  public render() {
    const { authenticated, authenticationView, authenticatedView } = this.props;
    if (authenticated()) {
      return authenticatedView;
    } else {
      return authenticationView;
    }
  }
}
