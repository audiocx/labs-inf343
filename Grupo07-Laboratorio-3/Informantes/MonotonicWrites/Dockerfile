FROM golang:latest

WORKDIR /app

COPY . .
RUN go mod download

CMD ["go", "run", "Informantes/ReadYourWrites/main.go", "dist028:4000", "dist025:5000", "dist026:5001", "dist027:5002"]