FROM golang:1.22.3

WORKDIR /app

COPY . .

RUN go get github.com/tursodatabase/libsql-client-go/libsql

RUN go build -o joker

EXPOSE 8080

CMD ["./joker"]
