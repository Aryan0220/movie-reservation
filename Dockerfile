FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM golang:1.26.2-alpine

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]