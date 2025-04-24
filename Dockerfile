FROM golang:1.23.4 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o lyriclink .

FROM debian:stable-slim

COPY --from=builder /app/lyriclink /bin/lyriclink

COPY frontend /frontend
COPY internal /internal
COPY sql /sql
COPY supabase /supabase

CMD ["/bin/lyriclink"]
