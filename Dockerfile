FROM golang:1.22

RUN apt-get update && apt-get install -y git ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .

CMD ["./server"]
