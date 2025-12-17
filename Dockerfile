FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api cmd/api/main.go

FROM scratch

COPY --from=builder /app/api /api

EXPOSE 8080

CMD ["/api"]
