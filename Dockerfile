# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем зависимости и скачиваем их
COPY app/go.mod ./
# RUN go mod download

COPY app .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/http-server ./cmd/http-server/main.go

RUN apk add --no-cache openssl && \
    openssl req -x509 -newkey rsa:4096 -nodes -keyout localhost.key -out localhost.crt -days 365 -subj "/CN=localhost"

FROM alpine:3.21 AS runner

RUN adduser -D -u 1000 appuser

USER appuser

# Копируем бинарник
COPY --from=builder --chown=appuser /bin/http-server /app/

# Копируем статику и LTS
COPY --chown=appuser app/cmd/http-server/static /app/static
COPY --from=builder --chown=appuser /app/localhost.crt /app/localhost.key /app/

WORKDIR /app

EXPOSE 4443

CMD ["./http-server"]



