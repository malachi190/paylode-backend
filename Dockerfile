# ---- build stage ----
FROM golang:1.23-alpine AS builder
WORKDIR /app
# install ca-certs for HTTPS calls (useful later)
RUN apk add --no-cache ca-certificates
# copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download
# copy source & build static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd

# ---- final minimal image ----
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# copy the single static binary
COPY --from=builder /app/server .

EXPOSE 8000
CMD ["./server"]