FROM golang:1.17-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .

RUN go build -o /docker-restapi

## Deploy
FROM alpine:latest

WORKDIR /

COPY --from=builder /docker-restapi /docker-restapi
COPY --from=builder /app /

EXPOSE 8083

ENTRYPOINT ["/docker-restapi"]