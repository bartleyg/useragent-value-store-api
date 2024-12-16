FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd

FROM alpine:latest
COPY --from=builder /app/main .
EXPOSE 80
CMD ["./main"]
