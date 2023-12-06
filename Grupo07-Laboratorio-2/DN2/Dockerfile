FROM golang:latest

WORKDIR /app

COPY . .
RUN go mod download

CMD ["go", "run", "DN2/main.go"]