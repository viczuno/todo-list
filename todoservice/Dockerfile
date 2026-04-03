FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY . .
COPY .env ./
RUN go mod download

RUN CGO_ENABLED=0 go build -o bin/main ./cmd/main.go


FROM alpine:3.21
WORKDIR /
COPY --from=builder /app/bin/main /main
COPY --from=builder /app/.env ./

ENTRYPOINT ["./main"]
EXPOSE 5000

