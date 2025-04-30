# ---------- Build stage ----------
FROM golang:1.23 AS builder

WORKDIR /app

# 加入 go module 檔案並下載依賴
COPY go.mod go.sum ./
RUN go mod download

# 複製程式與金鑰目錄
COPY . .
COPY keys ./keys

# 編譯 Go 程式
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# ---------- Deploy stage ----------
FROM gcr.io/distroless/static

# 設定執行目錄
WORKDIR /

# 複製編譯後執行檔
COPY --from=builder /app/main .

# 指定啟動指令
CMD ["/main"]

# 如果你的服務監聽 8080，就打開它（Render 使用）
EXPOSE 8080
