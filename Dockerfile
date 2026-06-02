FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# IMPORTANT: static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd

FROM scratch

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/uploads ./uploads

EXPOSE 8080

CMD ["./main"]