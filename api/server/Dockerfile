FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod init server || true
RUN go mod tidy

RUN go build -o server .

EXPOSE 8080

CMD ["./server"]
