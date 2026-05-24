FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/gateway-service ./cmd/app

FROM alpine:3.20
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 app
WORKDIR /app
COPY --from=builder /out/gateway-service /app/gateway-service
USER app
EXPOSE 8081
ENTRYPOINT ["/app/gateway-service"]
