FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod init client || true
RUN go mod tidy

RUN go build -o client .

CMD ["./client"]
