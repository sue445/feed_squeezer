FROM golang:1.25-alpine AS build-env

RUN apk add --no-cache make

ADD . /work
WORKDIR /work

ARG REVISION=dev
RUN make REVISION=${REVISION}

FROM alpine
COPY --from=build-env /work/bin/feed_squeezer /app/feed_squeezer

ENTRYPOINT ["/app/feed_squeezer"]
