FROM golang:1.22-alpine3.20 AS build-env

RUN apk add --no-cache make

ADD . /work
WORKDIR /work

RUN make

FROM alpine:3.20
COPY --from=build-env /work/bin/feed_squeezer /app/feed_squeezer

ENTRYPOINT ["/app/feed_squeezer"]
