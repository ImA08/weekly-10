FROM golang:1.24.4-alpine3.22 AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ENV GO_ENV=production

RUN go build -o server ./cmd/main.go

# STAGE 2
FROM alpine:3.21.3

WORKDIR /root/

COPY --from=builder /app/server .

RUN chmod +x server

EXPOSE 8080

CMD [ "./server" ]