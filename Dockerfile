FROM golang:alpine

COPY wait-for-it.sh /wait-for-it.sh

COPY go.mod go.sum /go/src/app/
COPY ./src /go/src/app

WORKDIR /go/src/app

RUN go mod download
RUN go build -o rinha-de-backend-2024-q1 .

RUN chmod +x /wait-for-it.sh
RUN apk update && apk add bash

EXPOSE 8080

CMD ["/wait-for-it.sh", "db:5432", "--", "./rinha-de-backend-2024-q1"]