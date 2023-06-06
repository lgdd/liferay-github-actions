FROM golang:1.20.4-alpine3.18

RUN apk update && \
  apk upgrade && \
  apk add git

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY main.go .

RUN go build -o /bin/app main.go

ENTRYPOINT ["app"]