#Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
#COPY go.mod go.sum ./
#RUN go mod download

COPY . .
RUN go build -v -o main main.go

#Run Stage
FROM alpine:3.13
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/main .
COPY app.env .
EXPOSE 8080

CMD ["/usr/src/app/main"]
