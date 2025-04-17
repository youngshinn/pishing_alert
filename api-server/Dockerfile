FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

#  정적 빌드 설정 (CGO 비활성화, glibc 의존 제거)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

#  실행용 이미지: 가볍고 안전한 alpine
FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/server .

# 실행 권한 부여
RUN chmod +x /app/server

EXPOSE 8080
CMD ["./server"]
