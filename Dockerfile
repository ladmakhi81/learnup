FROM golang:latest AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./bin ./cmd/app/api/main.go

FROM golang:latest AS tusdbuilder
RUN git clone https://github.com/tus/tusd.git /app/tusd
WORKDIR /app/tusd
ENV GOPROXY=https://goproxy.cn
RUN go mod download
RUN go build -o /app/bin/tusd ./cmd/tusd/main.go

FROM ubuntu:latest
RUN apt-get update && apt-get install -y ffmpeg
WORKDIR /app
COPY --from=builder /app/bin .
COPY --from=builder /app/translations /app/translations
CMD ["./bin"]