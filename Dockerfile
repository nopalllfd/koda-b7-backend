FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

# bikin folder di runtime (INI YANG BENAR)
RUN mkdir -p uploads

EXPOSE 8080

CMD ["./main"]