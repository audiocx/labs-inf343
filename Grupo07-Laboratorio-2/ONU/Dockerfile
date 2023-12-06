FROM golang:latest

WORKDIR /app

COPY . .
RUN go mod download

ENTRYPOINT ["go", "run", "ONU/main.go"]