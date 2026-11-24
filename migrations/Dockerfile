FROM golang:1.27.3-alpine AS builder

RUN apk add --no-cache bash curl

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.20.1/migrate.linux-amd64.tar.gz | tar xvz -C /usr/local/bin && \
    chmod +x /usr/local/bin/migrate

FROM alpine:3.23

RUN apk add --no-cache bash

COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

WORKDIR /app/migrate

ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_USER=user
ENV DB_PASSWORD=password
ENV DB_NAME=name
ENV MIGRATIONS_DIR=dir

COPY migrate.sh .
COPY migrations/ ./migrations/

RUN chmod +x migrate.sh
RUN apk --update --no-cache add bash postgresql-client curl jq

CMD ["./migrate.sh"]

