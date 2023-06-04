FROM golang:1.20.4-alpine3.18

WORKDIR /app

COPY main.go .

RUN go build -o /bin/app main.go

ENTRYPOINT ["app"]