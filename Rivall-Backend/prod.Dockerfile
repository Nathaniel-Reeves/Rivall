# Build environment
# -----------------
    FROM golang:1.23.4-alpine as build-env
    WORKDIR /Rivall-Backend
    
    RUN apk add --no-cache gcc musl-dev
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    
    RUN go build -ldflags '-w -s' -a -o ./bin/api ./cmd/api
    
    
    # Deployment environment
    # ----------------------
    FROM alpine
    
    COPY --from=build-env /Rivall-Backend/bin/api /Rivall-Backend/
    
    EXPOSE 8080
    CMD ["/Rivall-Backend/api"]