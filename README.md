# feed_squeezer
_feed_squeezer_ returns a new feed squeezed by any keyword in feed

[![Latest Version](https://img.shields.io/github/v/release/sue445/feed_squeezer)](https://github.com/sue445/feed_squeezer/releases)
[![test](https://github.com/sue445/feed_squeezer/actions/workflows/test.yml/badge.svg)](https://github.com/sue445/feed_squeezer/actions/workflows/test.yml)
[![docker-ghcr](https://github.com/sue445/feed_squeezer/actions/workflows/docker-ghcr.yml/badge.svg)](https://github.com/sue445/feed_squeezer/actions/workflows/docker-ghcr.yml)
[![docker-gcp](https://github.com/sue445/feed_squeezer/actions/workflows/docker-gcp.yml/badge.svg)](https://github.com/sue445/feed_squeezer/actions/workflows/docker-gcp.yml)
[![Coverage Status](https://coveralls.io/repos/github/sue445/feed_squeezer/badge.svg?branch=main)](https://coveralls.io/github/sue445/feed_squeezer?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/sue445/feed_squeezer)](https://goreportcard.com/report/github.com/sue445/feed_squeezer)
[![Go Reference](https://pkg.go.dev/badge/github.com/sue445/feed_squeezer.svg)](https://pkg.go.dev/github.com/sue445/feed_squeezer)

## Getting started
Run _feed_squeezer_ where it can be accessed from the Internet.

e.g.

### Docker
This application is provided as a Docker image, so you can run it wherever you like.

#### Images
* [GitHub Container Registry](https://github.com/sue445/feed_squeezer/pkgs/container/feed_squeezer) **(Recommended)**
  * `ghcr.io/sue445/feed_squeezer:latest`: Use latest version
  * `ghcr.io/sue445/feed_squeezer:vX.Y.Z`: Use specified version
* [Google Artifact Registry](https://console.cloud.google.com/artifacts/docker/feed-squeezer/us/feed-squeezer/app): If you want to run this app on [Cloud Run](https://cloud.google.com/run), use this image
  * `us-docker.pkg.dev/feed-squeezer/feed-squeezer/app:latest`: Use latest version
  * `us-docker.pkg.dev/feed-squeezer/feed-squeezer/app:vX.Y.Z`: Use specified version
  * `us-docker.pkg.dev/feed-squeezer/feed-squeezer/app:edge`: The contents of the main branch are pushed to this tag

```bash
docker run --rm -p 8080:8080 ghcr.io/sue445/feed_squeezer:latest
```

open http://localhost:8080/

### standalone binary
Download latest binary from https://github.com/sue445/feed_squeezer/releases

```bash
./feed_squeezer
```

open http://localhost:8080/

### Build yourself
```bash
git clone https://github.com/sue445/feed_squeezer.git
cd feed_squeezer
make
./bin/feed_squeezer
```

open http://localhost:8080/

## Environment variables
* `PORT`(optional): http listening port. default is `8080`
* `SENTRY_DSN`, `SENTRY_RELEASE`, `SENTRY_ENVIRONMENT`(optional): [Sentry](https://sentry.io/) configuration. See followings
  * https://github.com/getsentry/sentry-go
  * https://docs.sentry.io/platforms/go/configuration/

## Cli usage
```bash
$ ./feed_squeezer --help
Usage of ./bin/feed_squeezer:
  -version
        Whether showing version
```

## Endpoint
### GET /
Display a simple form to generate _feed_squeezer_ URL

![top](doc/top.png)

### GET /api/feed
returns a new feed squeezed by any keyword in feed

#### Parameters
All parameters must be URL encoded

* `url` : source feed url
* `query` : query to squeeze feed
  * The following formats are supported
  * `AAA BBB` : Includes all (AND search)
  * `AAA | BBB` : Includes any (OR search)
  * `(AAA BBB) | CCC`, `(AAA | BBB) CCC` : Evaluate conditions in brackets first

### GET /api/version
Returns app version (same to `feed_squeezer -version`)

## LICENSE
All programs are licensed under the [MIT License](LICENSE) Copyright (c) 2024 sue445.

But only [favicon](public/favicon.svg)'s LICENCE belongs to [TopeconHeros](https://icooon-mono.com/)

Original icon is [here](https://icooon-mono.com/11500-rss%E3%81%AE%E3%82%A2%E3%82%A4%E3%82%B3%E3%83%B3/)
