FROM golang:1.23 AS builder

WORKDIR pinger

COPY . .

RUN go mod download && go mod verify
RUN go build -v -o /usr/bin/pinger ./cmd/pinger/main.go

FROM ubuntu:25.04

WORKDIR app

COPY config/config.yaml config.yaml
COPY --from=builder /usr/bin/pinger /usr/bin/pinger

RUN apt update && apt install -y iputils-ping

CMD ["pinger", "-c", "./config.yaml"]