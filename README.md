# feed_proxy

## Running
```
make
./bin/feed_proxy
```

open http://localhost:8080/

## Environment variables
* `PORT`(optional): http listening port. default is `8080`
* `SENTRY_DSN`, `SENTRY_RELEASE`, `SENTRY_ENVIRONMENT`(optional): [Sentry](https://sentry.io/) configuration. See followings
  * https://github.com/getsentry/sentry-go
  * https://docs.sentry.io/platforms/go/configuration/
