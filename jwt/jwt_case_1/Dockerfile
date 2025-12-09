# --- Stage 1: Builder (Derleme Aşaması) ---
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o jwt-case-1-app .

# --- Stage 2: Runner (Çalıştırma Aşaması) ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/jwt-case-1-app .
EXPOSE 8080
CMD ["./jwt-case-1-app"]