# Build stage
FROM golang:1.25.3-alpine3.22 as builder

WORKDIR /app

COPY . .

RUN go build -o main main.go

# Run Stage
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080

CMD ["/app/main"]