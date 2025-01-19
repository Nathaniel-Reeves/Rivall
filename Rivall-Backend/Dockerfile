FROM golang:1.23.4-alpine
WORKDIR /Rivall-Backend

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/api ./cmd/api

CMD ["/Rivall-Backend/bin/api"]
EXPOSE 8080