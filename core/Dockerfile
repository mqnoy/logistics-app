# Stage builder
FROM golang:1.21.11-alpine3.20 AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy && \
    go mod download && \
    go build --o core && \
    chmod +x core

# Stage runner
FROM golang:1.21.11-alpine3.20

# Install dependencies
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY --from=builder /app/core .

CMD ["/app/core"]