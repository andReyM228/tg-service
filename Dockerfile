FROM golang:latest
WORKDIR /app
COPY . .
RUN chmod +x ./cmd
RUN go build cmd/main.go

CMD ["./main"]