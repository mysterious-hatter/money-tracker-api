FROM golang:1.22 AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM scratch

WORKDIR /app
COPY --from=builder /build/main ./main
COPY --from=builder /build/.env ./env

ENTRYPOINT ["/app/main"]