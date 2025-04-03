FROM golang:latest AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./bin ./cmd/api/main.go

FROM ubuntu:latest
WORKDIR /app
COPY --from=builder /app/bin .
CMD ["./bin"]