# feed_squeezer
[![test](https://github.com/sue445/feed_squeezer/actions/workflows/test.yml/badge.svg)](https://github.com/sue445/feed_squeezer/actions/workflows/test.yml)
[![docker-ghcr](https://github.com/sue445/feed_squeezer/actions/workflows/docker-ghcr.yml/badge.svg)](https://github.com/sue445/feed_squeezer/actions/workflows/docker-ghcr.yml)

## Running
```
make
./bin/feed_squeezer
```

open http://localhost:8080/

## Environment variables
* `PORT`(optional): http listening port. default is `8080`
* `SENTRY_DSN`, `SENTRY_RELEASE`, `SENTRY_ENVIRONMENT`(optional): [Sentry](https://sentry.io/) configuration. See followings
  * https://github.com/getsentry/sentry-go
  * https://docs.sentry.io/platforms/go/configuration/
