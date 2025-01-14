# Base Image
FROM golang:1.23.4

WORKDIR /go/src/app
COPY ./src .

RUN go get -d -v
RUN go build -v

CMD ["./Rivall-Backend"]