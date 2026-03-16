FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o apollo ./cmd/apollo

# ---

FROM alpine:3.18

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/apollo .

EXPOSE 4000
ENTRYPOINT ["./apollo"]
