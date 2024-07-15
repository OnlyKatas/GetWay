FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go build  ./cmd/main.go

FROM alpine:3.15
# Copy binary from builder

WORKDIR /app
COPY --from=builder /app/.env .
COPY --from=builder /app/config/local.yaml ./config/local.yaml
COPY --from=builder /app/main .
EXPOSE 50051

# Set entrypoint
CMD ["/app/main"]