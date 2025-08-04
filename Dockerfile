FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env* ./

EXPOSE 8080

CMD ["./main"]

# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: secretdb
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - secret-network

  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://user:password@postgres:5432/secretdb?sslmode=disable
      JWT_SECRET: your-super-secret-jwt-key
      ENCRYPTION_KEY: your-32-byte-encryption-key-here12
      API_KEYS: api-key-1,api-key-2,api-key-3
      PORT: 8080
    depends_on:
      - postgres
    networks:
      - secret-network

volumes:
  postgres_data:

networks:
  secret-network:
    driver: bridge

# .env.example
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/secretdb?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key
ENCRYPTION_KEY=your-32-byte-encryption-key-here12
API_KEYS=api-key-1,api-key-2,api-key-3
