FROM golang:1.24-alpine

WORKDIR /app/server

COPY ./server /app/server

RUN go mod download
RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"] 