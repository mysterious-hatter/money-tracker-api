FROM golang:1.22 AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./back-frontend

FROM scratch

WORKDIR /app
COPY --from=builder /build/back-frontend ./back-frontend
COPY --from=builder /build/.env ./env

ENTRYPOINT ["/app/back-frontend"]