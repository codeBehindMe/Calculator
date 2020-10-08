FROM golang:1.14

RUN mkdir /app

WORKDIR /app

COPY . .

ENV GO111MODULE=on

RUN go build main.go

ENTRYPOINT ["./main"]