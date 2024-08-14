FROM golang:1.23-alpine3.20 AS build-env

RUN apk add --no-cache make

ADD . /work
WORKDIR /work

ARG REVISION=dev
RUN make REVISION=${REVISION}

FROM alpine:3.20
COPY --from=build-env /work/bin/feed_squeezer /app/feed_squeezer

ENTRYPOINT ["/app/feed_squeezer"]
