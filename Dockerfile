FROM golang:alpine

WORKDIR /app
COPY . .

RUN go build ./cmd/server

CMD ["./server"]