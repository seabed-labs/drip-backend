FROM 434618599721.dkr.ecr.us-east-1.amazonaws.com/golang:latest

RUN apk add --no-cache git
RUN apk add --no-cache bash

ARG DISCORD_ACCESS_TOKEN
ARG DISCORD_WEBHOOK_ID
ARG DRIP_BACKEND_WALLET
ARG DRIP_PROGRAM_ID
ARG PORT
ARG ENV
ARG GOOGLE_CLIENT_ID
ARG NETWORK
ARG PSQL_USER
ARG PSQL_PASS
ARG PSQL_DBNAME
ARG PSQL_PORT
ARG PSQL_HOST

ENV DISCORD_ACCESS_TOKEN=${DISCORD_ACCESS_TOKEN}
ENV DISCORD_WEBHOOK_ID=${DISCORD_WEBHOOK_ID}
ENV DRIP_BACKEND_WALLET=${DRIP_BACKEND_WALLET}
ENV DRIP_PROGRAM_ID=${DRIP_PROGRAM_ID}
ENV PORT=${PORT}
ENV ENV=${ENV}
ENV GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID}
ENV NETWORK=${NETWORK}
ENV PSQL_USER=${PSQL_USER}
ENV PSQL_PASS=${PSQL_PASS}
ENV PSQL_DBNAME=${PSQL_DBNAME}
ENV PSQL_PORT=${PSQL_PORT}
ENV PSQL_HOST=${PSQL_HOST}

WORKDIR /drip-backend

COPY go.mod .
COPY go.sum .
RUN go mod download
# https://github.com/montanaflynn/golang-docker-cache
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

COPY main.go main.go
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg


RUN go build -o app ./main.go

CMD ["./app"]