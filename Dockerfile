FROM golang:1.22.3-alpine as base
WORKDIR /usr/src/app

FROM base as build
RUN mkdir -p /temp/dev
COPY go.mod go.sum /temp/dev/
RUN cd /temp/dev && go mod download && go mod verify

COPY . .
RUN go build

FROM alpine:3.20.0 as prod
COPY --from=build /usr/src/app/Gateway ./gateway

ENTRYPOINT ["./gateway"]