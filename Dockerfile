# ---------- Build stage ----------

# 使用 Debian-based Go 官方映像
FROM golang:1.21

# 安裝必要工具（git 和 ca-certificates）
RUN apt-get update && apt-get install -y \
  git \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*

# 設定工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum，先做依賴快取
COPY go.mod go.sum ./
RUN go mod download

# 複製專案所有檔案
COPY . .

# 編譯 Go 程式為 server 可執行檔
RUN go build -o server .

# 不要 COPY .env，Render 上會用環境變數設定；本機使用 `--env-file .env`

# 設定容器啟動指令
CMD ["./server"]
