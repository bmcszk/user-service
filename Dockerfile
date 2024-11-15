FROM golang:1.23-alpine3.20 AS base
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

FROM base AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /app/service .

FROM gcr.io/distroless/static-debian12 AS runner
WORKDIR /app
COPY --from=builder /app/service .
COPY db/migrations /app/db/migrations
EXPOSE 8080
CMD ["/app/service"]
