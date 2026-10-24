FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /usr/local/bin/privai-server ./cmd/server

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /usr/local/bin/privai-server /usr/local/bin/privai-server
WORKDIR /data
ENV DB_PATH=/data/privai.db
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/privai-server"]
