# ---------- Build Stage ----------
FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /app/bin

# Build các binary
RUN go build -o /app/bin/go-blockchain ./cmd/node
RUN go build -o /app/bin/create_wallet ./cmd/cli/create_wallet.go
RUN go build -o /app/bin/send_tx ./cmd/cli/send_tx.go
RUN go build -o /app/bin/status ./cmd/cli/status.go

# Copy wait-for-it.sh nếu bạn có file đó trong source
COPY ./wait-for-it.sh /app/bin/wait-for-it.sh
RUN chmod +x /app/bin/wait-for-it.sh

# ---------- Run Stage ----------
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache libc6-compat

# Copy binary đã build
COPY --from=builder /app/bin/* ./

# Tạo thư mục ví và dữ liệu
RUN mkdir -p /app/wallets /app/data

CMD ["./go-blockchain"]
