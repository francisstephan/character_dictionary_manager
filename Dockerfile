# syntax=docker/dockerfile:1

##
## Build the application from source
##
# https://docs.docker.com/language/golang/build-images/

FROM golang:1.19 AS build-stage

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags "-s -w" -buildvcs=false -o /zdmanager

##
## Deploy the application binary into a lean image
##

FROM busybox AS build-release-stage

ENV GIN_MODE=release

WORKDIR /

COPY --from=build-stage /zdmanager /zdmanager

COPY ./vol ./vol

EXPOSE 8080

ENTRYPOINT ["/zdmanager"]
