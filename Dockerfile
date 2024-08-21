FROM golang:1.22

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /go-import-service

EXPOSE 8080

EXPOSE 8081

CMD ["/go-import-service"]
