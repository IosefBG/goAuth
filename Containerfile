FROM golang:latest

WORKDIR /goAuth

COPY ../go.mod .
COPY ./cmd/main.go .

RUN go mod download

COPY . .

RUN go build -o app

EXPOSE 8080

CMD ["./app"]
