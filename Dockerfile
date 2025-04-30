# ---------- Build stage ----------

# 使用 Go 官方映像
FROM golang:1.21-slim

# 安裝必要工具（包含 git 和 ca-certificates）
RUN apk add --no-cache git ca-certificates

# 設定工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum，先做依賴快取
COPY go.mod go.sum ./
RUN go mod download

# 複製專案所有檔案
COPY . .

# 編譯 Go 程式為靜態檔案（包含 main.go）
RUN go build -o server .

# 加入 .env（供開發測試用，Render 會忽略）
# 注意：此行不是在容器內跑 .env，而是可以選擇在 docker run 時搭配
# docker run --env-file .env ...
# 不需要 COPY .env 到映像中以避免 secrets 外洩

# 設定容器啟動指令
CMD ["./server"]
