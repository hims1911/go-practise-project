# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app
ADD . /app/

RUN go build -o ./out/go-practise-project .

EXPOSE 8000

ENTRYPOINT ["./out/go-practise-project"]