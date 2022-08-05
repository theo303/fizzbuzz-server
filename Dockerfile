# syntax=docker/dockerfile:1

FROM golang:1.19

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN git clone
RUN go build -o /fizzbuzz-server

ENV PORT=8080

CMD [ "/fizzbuzz-server" ]