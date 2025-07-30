FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY *.mod *.sum ./
RUN go mod download
COPY . .
RUN go build -o subscription_service cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/subscription_service ./
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./subscription_service"]