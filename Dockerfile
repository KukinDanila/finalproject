FROM golang:1.26.8 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/scheduler .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/scheduler /app/scheduler
COPY --from=builder /app/web /app/web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db

EXPOSE 7540

CMD ["/app/scheduler"]