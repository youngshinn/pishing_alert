FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

#  정적 빌드: glibc 문제 해결
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o collector main.go

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/collector .

RUN chmod +x /app/collector

CMD ["./collector"]
