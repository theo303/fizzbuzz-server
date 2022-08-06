# syntax=docker/dockerfile:1

FROM golang:1.19 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o /fizzbuzz-server

FROM debian:stable-slim

WORKDIR /

COPY --from=builder /fizzbuzz-server /fizzbuzz-server

EXPOSE 8080
ENV PORT=8080

CMD [ "/fizzbuzz-server" ]