export const isSignedIn = (token?: string): token is string => {
  return token !== undefined && token !== null;
};
