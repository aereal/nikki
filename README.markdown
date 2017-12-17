# Nikki

```
env GOOGLE_OAUTH_CLIENT_ID='xxx' GOOGLE_OAUTH_CLIENT_SECRET='xxx' bundle exec rackup
```

# Setup

```
env GOOGLE_OAUTH_CLIENT_ID='xxx' GOOGLE_OAUTH_CLIENT_SECRET='xxx' docker-compose up # Run DB
cat db/nikki_dump.sql | docker-compose run --rm db psql -h db -U postgres
sudo cat <<EOS >> /etc/hosts
127.0.0.1\tapi.nikki.dev
127.0.0.1\tadmin.nikki.dev
EOS
open etc/certs/ca.crt # Always trust
```

# See also

- [Creating a Google API Console project and client ID  |  Google Sign-In for Websites  |  Google Developers](https://developers.google.com/identity/sign-in/web/devconsole-project)
  - and make sure Google+ API available because omniauth-google-oauth2 requires it
