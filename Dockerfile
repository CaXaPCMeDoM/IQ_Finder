FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN apk add --no-cache git
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -tags migrate -o migrate-service -ldflags="-X 'main.migrationsPath=file:///app/migrations'" ./cmd/migrate/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o name_iq_finder ./cmd/app/main.go

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /app/migrate-service .
COPY --from=builder /app/name_iq_finder .
COPY --from=builder /app/migrations ./migrations

RUN chmod +x migrate-service name_iq_finder

CMD ["sh", "-c", "sleep 20 && ./migrate-service && ./auth-service"]