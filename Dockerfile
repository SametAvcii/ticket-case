FROM golang:alpine as builder

LABEL maintainer="Samet Avcı"

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download 

COPY . .

# Build the Go app
RUN go build -o main /app/cmd/main.go
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .     

EXPOSE 8080

CMD ["./main"]