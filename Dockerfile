#Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
#COPY go.mod go.sum ./
#RUN go mod download

COPY . .
RUN go build -v -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.0/migrate.linux-amd64.tar.gz | tar xvz

#Run Stage
FROM alpine:3.13
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/main .
COPY --from=builder /usr/src/app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY db/migration ./migration
EXPOSE 8080

CMD ["/usr/src/app/main"]
ENTRYPOINT ["/usr/src/app/start.sh"]
