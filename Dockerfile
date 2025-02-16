FROM golang:1.22 AS builder

WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY .env ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o fin-backend

FROM scratch

WORKDIR /app
COPY --from=builder /build/fin-backend ./fin-backend
COPY --from=builder /build/.env ./env

ENTRYPOINT ["/app/fin-backend"]