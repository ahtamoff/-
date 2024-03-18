FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go .
RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/

RUN apk add --no-cache postgresql-dev
RUN apk add --no-cache ca-certificates

ENV PGUSER=posgres\
    PGPASSWORD=postgres\
    PGDATABASE=film_library\
    PGHOST=localhost\
    PGPORT=5432

CMD ["/app/main"]
