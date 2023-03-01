# syntax=docker/dockerfile:1

## Build
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /docker-gs-ping

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /docker-gs-ping /docker-gs-ping

EXPOSE 3000

USER nonroot:nonroot

ENV DB_USER=doadmin
ENV DB_PASSWORD=AVNS_6w64qBFeRo1FGhSsE24
ENV DB_HOST=db-mysql-fmf-do-user-7517862-0.b.db.ondigitalocean.com
ENV DB_PORT=25060
ENV DB_DATABASE=fmf

ENTRYPOINT ["/docker-gs-ping"]
